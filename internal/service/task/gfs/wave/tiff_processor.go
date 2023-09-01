package wave

import (
	"fmt"
	"strings"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/lukeroth/gdal"
	"golang.org/x/exp/slices"
	"yh-process/internal/consts"
	"yh-process/internal/model"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
)

/*
   @File: tiff_processor.go
   @Author: khaosles
   @Time: 2023/8/22 23:13
   @Desc:
*/

type GfsWaveTiffDataProcessor struct {
	node2.BaseNode
}

func NewGfsWaveTiffDataProcessor(id, proxy string, cache node2.Cache) node2.Node {
	var n GfsWaveTiffDataProcessor
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *GfsWaveTiffDataProcessor) Execute() *node2.Result {

	var err error
	var op16, op25, op16nc, op25nc, tifPath string
	var result node2.Result
	var nodata float32
	var jsonObject = gjson.NewJsonObject()
	var elements []*model.SysElementInfo
	var nclist []string

	glog.Infof("execute nc to tif ===> %s", n.Id)

	// 解析id
	idsplit := strings.Split(n.Id, ":")

	dependentId := n.Dependents.Values()[0]
	parent, ok := n.NodeCahce.Get(dependentId)
	if !ok {
		err = fmt.Errorf("node cache not exist ===> %s", dependentId)
		goto returntag
	}
	// 保存路径
	tifPath, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisTifDir, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 无效值
	nodata, err = redisutil.GetOfF32(fmt.Sprintf(consts.RedisNodata, idsplit[0]))
	if err != nil {
		goto returntag
	}

	op16, err = parent.GetResult().Msg.GetString("0p16")
	if err != nil {
		goto returntag
	}
	op25, err = parent.GetResult().Msg.GetString("0p25")
	if err != nil {
		goto returntag
	}
	op16nc = fmt.Sprintf("%s.nc", strings.TrimSuffix(op16, ".grib2"))
	op25nc = fmt.Sprintf("%s.nc", strings.TrimSuffix(op25, ".grib2"))
	// nc文件路径
	nclist = []string{op16nc, op25nc}
	// 执行 wgrib2
	for i, gribfile := range []string{op16, op25} {
		if !fileutil.IsFile(gribfile) {
			err = fmt.Errorf("file not exist ===> %s", op16)
			goto returntag
		}
		if !fileutil.IsFile(nclist[i]) {
			// 执行cmd grib2转cn
			err = util.ExecGrib2Netcdf(gribfile, nclist[i])
			if err != nil {
				goto returntag
			}
		}
		if !fileutil.IsFile(nclist[i]) {
			err = fmt.Errorf("file not exist ===> %s", nclist[i])
			goto returntag
		}
	}
	// 获取需要提取的要素
	elements, err = redisutil.GetOfElements(fmt.Sprintf(consts.RedisFieldBase, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 提取基础要素
	for _, element := range elements {
		outfile := fileutil.Join(tifPath, element.Elemname, fmt.Sprintf("%s_%s.tif", element.Elemname, idsplit[2]))
		err = Nc2tif(op16nc, op25nc, outfile, element, nodata)
		if err != nil {
			goto returntag
		}
		if !fileutil.IsFile(outfile) {
			err = fmt.Errorf("file not exist ===> %s", outfile)
			goto returntag
		}
		jsonObject.Put(element.Elemname, outfile)
	}

returntag:
	if err != nil {
		// 失败
		result.Success = false
		result.Err = err.Error()
		n.Status = node2.Failed
	} else {
		// 成功
		result.Success = true
		result.Msg = jsonObject
		n.Status = node2.Succeed
	}
	result.Id = n.Id
	n.Result = &result
	return &result
}

func readData(infile string, element *model.SysElementInfo, nodata float32) ([][]float32, float64, float64, float64, float64, error) {
	// 校验文件
	if !fileutil.IsFile(infile) {
		return nil, 0, 0, 0, 0, fmt.Errorf("file not exist: %s", infile)
	}

	// 打开文件
	ds, err := netcdf.Open(infile)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	defer ds.Close()

	// //////////////////////////////////////////
	// 经度
	lonVar, err := ds.GetVariable("longitude")
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	lon, _ := lonVar.Values.([]float64)
	// //////////////////////////////////////////
	// 纬度
	latVar, err := ds.GetVariable("latitude")
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	lat, _ := latVar.Values.([]float64)
	// 处理经度
	for i := 0; i < len(lon); i++ {
		lon[i] -= 180
	}
	// 反转纬度
	var reverse bool
	if lat[0] < lat[1] {
		reverse = true
		lat = util.Reverse1d(lat)
	}
	// 数据字段
	dataVar, err := ds.GetVariable(element.Field)
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}
	// 获取无效值
	fillValue, _ := dataVar.Attributes.Get("_FillValue")
	f, _ := fillValue.(float32)
	// 获取数据
	d3, _ := dataVar.Values.([][][]float32)
	data := d3[0]
	mid := len(data[0]) / 2
	// 镜像反转
	for i, datum := range data {
		vv := make([]float32, len(datum))
		copy(vv[:mid], datum[mid:])
		copy(vv[mid:], datum[:mid])
		data[i] = vv
	}

	// 取出无效值
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			if data[i][j] == f {
				data[i][j] = nodata
				continue
			}
			data[i][j] = data[i][j]*element.Gain + element.Offset
		}
	}
	// 反转数据
	if reverse {
		data = util.Reverse2d(data)
	}
	return data, slices.Min(lon), slices.Max(lon), slices.Min(lat), slices.Max(lat), nil
}

func Nc2tif(op16nc, op25nc, outfile string, element *model.SysElementInfo, nodata float32) error {
	// 读取数据 .16
	data16, _, _, _, ymax16, err := readData(op16nc, element, nodata)
	if err != nil {
		return nil
	}
	// 读取数据 .25
	data25, _, _, _, ymax25, err := readData(op25nc, element, nodata)
	if err != nil {
		return nil
	}
	// 保存为dataset
	ds25, err := util.ImgWriteMem(&util.Img[float32]{Data: data25,
		Tran: [6]float64{-180, 0.25, 0, 90, 0, -0.25}, Proj: consts.Projection})
	if err != nil {
		return err
	}
	defer ds25.Close()

	// 创建选项对象并设置选项
	options := []string{
		"-of", "MEM",
		"-te", "-180", "-90", "180", "90", // outputBounds
		"-ts", "2160", "1080", // width, height
		"-r", "near",
	}
	// 重采样
	warpDs, err := gdal.Warp("", nil, []gdal.Dataset{*ds25}, options)
	defer warpDs.Close()
	if err != nil {
		return err
	}
	// 读取采样后数据
	dataWarp := make([]float32, 2160*1080)
	err = warpDs.RasterBand(1).IO(
		gdal.Read, 0, 0, 2160, 1080,
		dataWarp, 2160, 1080, 0, 0,
	)
	if err != nil {
		return err
	}
	// 融合后的数据
	var data = make([][]float32, 1080)
	// 计算.16开始索引
	start := int((ymax25 - ymax16) / 0.166666666666)
	// 计算.16结束索引
	end := len(data16) + start
	// 写入数据
	for i := 0; i < 1080; i++ {
		if i >= start && i < end {
			// 使用.16数据
			var tmp []float32
			for j := 0; j < 2160; j++ {
				if data16[i-start][j] != nodata {
					tmp = append(tmp, data16[i-start][j])
				} else {
					tmp = append(tmp, dataWarp[i*2160+j])
				}
				data[i] = tmp
			}
		} else {
			// 使用.25数据
			data[i] = dataWarp[i*2160 : (i+1)*2160]
		}
	}

	// 膨胀
	for i := 0; i < 5; i++ {
		data = util.Dilation(data, 3, nodata)
	}
	// 创建文件夹
	_ = fileutil.MkdirP(outfile)
	// 保存tiff
	return util.ImgWrite(
		outfile,
		&util.Img[float32]{
			Data:   data,
			Tran:   warpDs.GeoTransform(),
			Proj:   warpDs.Projection(),
			Nodata: -9999,
		})
}

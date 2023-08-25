package atmos

import (
	"fmt"
	"strings"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/model"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
	"yh-process/internal/vo"
)

/*
   @File: tiff_rocessor.go
   @Author: khaosles
   @Time: 2023/8/19 23:57
   @Desc:
*/

// GfsAtmosTiffDataProcessor nc转为tiff模块
type GfsAtmosTiffDataProcessor struct {
	node2.BaseNode
}

func NewGfsAtmosTiffDataProcessor(id, proxy string, cache node2.Cache) node2.Node {
	var n GfsAtmosTiffDataProcessor
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *GfsAtmosTiffDataProcessor) Execute() *node2.Result {

	var err error
	var gribFile, ncFile, tifPath string
	var result node2.Result
	var nodata float32
	var jsonObject = gjson.NewJsonObject()
	var elements []*model.SysElementInfo
	var uvs []*vo.UV

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

	// grib2文件路径
	gribFile, err = parent.GetResult().Msg.GetString("gribfile")
	if err != nil {
		goto returntag
	}
	if !fileutil.IsFile(gribFile) {
		err = fmt.Errorf("file not exist ===> %s", gribFile)
		goto returntag
	}

	// nc文件路径
	ncFile = fmt.Sprintf("%s.nc", gribFile)
	if !fileutil.IsFile(ncFile) {
		// 执行cmd grib2转cn
		err = util.ExecGrib2Netcdf(gribFile, ncFile)
		if err != nil {
			goto returntag
		}
	}
	if !fileutil.IsFile(ncFile) {
		err = fmt.Errorf("file not exist ===> %s", ncFile)
		goto returntag
	}

	// 获取需要提取的要素
	elements, err = redisutil.GetOfElements(fmt.Sprintf(consts.RedisFieldBase, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 提取基础要素
	for _, element := range elements {
		outfile := fileutil.Join(tifPath, element.Elemname, fmt.Sprintf("%s_%s.tif", element.Elemname, idsplit[2]))
		err = Nc2tif(ncFile, outfile, element, nodata)
		if err != nil {
			goto returntag
		}
		if !fileutil.IsFile(outfile) {
			err = fmt.Errorf("file not exist ===> %s", outfile)
			goto returntag
		}
		jsonObject.Put(element.Elemname, outfile)
	}

	// 处理uv转sd
	uvs, err = redisutil.GetOfUV(fmt.Sprintf(consts.RedisFieldUV, idsplit[0]))
	if err != nil {
		goto returntag
	}
	for _, uv := range uvs {
		ufile := fileutil.Join(tifPath, uv.U.Elemname, fmt.Sprintf("%s_%s.tif", uv.U.Elemname, idsplit[2]))
		vfile := fileutil.Join(tifPath, uv.V.Elemname, fmt.Sprintf("%s_%s.tif", uv.V.Elemname, idsplit[2]))
		sfile := fileutil.Join(tifPath, uv.Speed.Elemname, fmt.Sprintf("%s_%s.tif", uv.Speed.Elemname, idsplit[2]))
		dfile := fileutil.Join(tifPath, uv.Direction.Elemname, fmt.Sprintf("%s_%s.tif", uv.Direction.Elemname, idsplit[2]))
		err := util.Uv2sd(ufile, vfile, sfile, dfile, nodata)
		if err != nil {
			goto returntag
		}
		if !fileutil.IsFile(ufile) {
			err = fmt.Errorf("file not exist ===> %s", ufile)
			goto returntag
		}
		if !fileutil.IsFile(vfile) {
			err = fmt.Errorf("file not exist ===> %s", vfile)
			goto returntag
		}
		jsonObject.Put(uv.U.Elemname, ufile)
		jsonObject.Put(uv.V.Elemname, vfile)
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

func Nc2tif(infile, outfile string, element *model.SysElementInfo, nodata float32) error {
	// 校验文件
	if !fileutil.IsFile(infile) {
		return fmt.Errorf("file not exist: %s", infile)
	}

	// 打开文件
	ds, err := netcdf.Open(infile)
	if err != nil {
		return err
	}
	defer ds.Close()

	// //////////////////////////////////////////
	// 经度
	lonVar, err := ds.GetVariable("longitude")
	if err != nil {
		return err
	}
	lon, _ := lonVar.Values.([]float64)
	// //////////////////////////////////////////
	// 纬度
	latVar, err := ds.GetVariable("latitude")
	if err != nil {
		return err
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
		return err
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
	// 创建文件夹
	_ = fileutil.MkdirP(outfile)
	// tiff
	return util.ImgWrite(
		outfile,
		&util.Img[float32]{
			Data: data,
			Tran: [6]float64{lon[0], lon[1] - lon[0], 0, lat[0], 0, lat[1] - lat[0]},
			Proj: consts.Projection,
		})
}

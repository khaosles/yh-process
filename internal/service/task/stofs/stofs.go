package stofs

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/khaosles/giz/fileutil"
	"yh-process/internal/consts"
	"yh-process/internal/service/task/util"
)

/*
   @File: stofs.go
   @Author: khaosles
   @Time: 2023/8/24 16:18
   @Desc:
*/

var timeComplie = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)

const (
	xMin float64 = -180.
	xMax         = 180.
	yMin         = -90.
	yMax         = 90.
	res          = 0.1
)

// ToGrid 按照时刻转为格网文件
func ToGrid(infile, outpath, dataType string, gain, offset, nodata float64) error {

	if !fileutil.IsFile(infile) {
		return fmt.Errorf("file not exist: %s", infile)
	}
	// 打开文件
	group, err := netcdf.Open(infile)
	if err != nil {
		return err
	}

	// 经度
	variable, err := group.GetVariable("x")
	if err != nil {
		return err
	}
	lon, _ := variable.Values.([]float64)

	// 纬度
	variable, err = group.GetVariable("y")
	if err != nil {
		return err
	}
	lat, _ := variable.Values.([]float64)

	// 时间
	variable, err = group.GetVariable("time")
	if err != nil {
		return err
	}
	ts, _ := variable.Values.([]float64)
	val, _ := variable.Attributes.Get("units")
	s, _ := val.(string)
	timestr := timeComplie.FindString(s)
	baseTime, err := time.Parse(time.DateTime, timestr)
	if err != nil {
		return err
	}

	// 数据
	getter, err := group.GetVarGetter("zeta")
	if err != nil {
		return err
	}
	obj, _ := getter.Attributes().Get("_FillValue")
	fillvalue, _ := obj.(float64)

	// 生成格网
	var i int64
	wg := sync.WaitGroup{}
	limit := make(chan struct{}, 3)
	for i = 0; i < int64(len(ts)); i++ {
		wg.Add(1)
		iCp := i
		go func(i int64) {
			limit <- struct{}{}
			defer wg.Done()
			defer func() { <-limit }()

			grib := genGrib()                   // 生成格网
			initGrib(grib, float32(nodata))     // 格网初始化
			slice, _ := getter.GetSlice(i, i+1) // 读取数据
			data2d, _ := slice.([][]float64)
			data1d := data2d[0]
			for j := 0; j < len(lon); j++ {
				if data1d[j] == fillvalue {
					continue
				}
				xoff, yoff := calcIdByLonLat(lon[j], lat[j])
				v := float32(data1d[j]*gain + offset)
				if grib[yoff][xoff] == float32(nodata) {
					// 该点无值则直接填充
					grib[yoff][xoff] = v
				} else {
					// 有值则计算均值
					grib[yoff][xoff] = (grib[yoff][xoff] + v) / 2
				}
			}
			// 膨胀
			for k := 0; k < 10; k++ {
				grib = util.Dilation(grib, 3, float32(nodata))
			}
			ymdh := baseTime.Add(time.Duration(ts[i]) * time.Second).Format("2006-01-02_15")
			outfile := fileutil.Join(outpath, fmt.Sprintf("%s-%s.tif", ymdh, dataType))
			// 保存
			_ = util.ImgWrite(outfile, &util.Img[float32]{
				Data: grib, Tran: [6]float64{xMin, res, 0, yMax, 0, -res},
				Nodata: nodata, Proj: consts.Projection,
			})
		}(iCp)
	}
	wg.Wait()
	return nil
}

// 生成格网
func genGrib() [][]float32 {
	width := int((xMax - xMin) / res)
	height := int((yMax - yMin) / res)
	grib := make([][]float32, height)
	for i := 0; i < height; i++ {
		grib[i] = make([]float32, width)
	}
	return grib
}

// 初始化格网
func initGrib(grib [][]float32, fillValue float32) {
	for i := 0; i < len(grib); i++ {
		for j := 0; j < len(grib[0]); j++ {
			grib[i][j] = fillValue
		}
	}
}

// 根据经纬度计算索引位置
func calcIdByLonLat(lon, lat float64) (xOff, yOff int) {
	xOff = int((lon - xMin) / res)
	yOff = int((lat - yMax) / -res)
	return
}

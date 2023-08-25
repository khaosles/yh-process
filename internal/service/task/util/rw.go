package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/g"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/lukeroth/gdal"
)

/*
   @File: rw.go
   @Author: khaosles
   @Time: 2023/7/24 18:21
   @Desc:
*/

type Img[T g.Numeric] struct {
	Data   [][]T
	Tran   [6]float64
	Proj   string
	Nodata float64
}

// ImgRead 读取单波段数据
func ImgRead[T g.Numeric](infile string, isFlagTranOrProj, isFlagNodata bool) (*Img[T], error) {
	if !fileutil.IsFile(infile) {
		return nil, fmt.Errorf("file not exist: %s", infile)
	}

	ds, err := gdal.Open(infile, gdal.ReadOnly)
	if err != nil {
		return nil, fmt.Errorf("error opening GDAL file: %v", err)
	}
	defer ds.Close()

	// 获取栅格波段
	band := ds.RasterBand(1)
	// 获取栅格宽度和高度
	width := band.XSize()
	height := band.YSize()

	// 读取栅格数据
	data := make([]T, width*height)
	err = band.IO(gdal.Read, 0, 0, width, height, data, width, height, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("Error reading raster data: %v", err)
	}

	var img Img[T]
	// 将一维数据转换为二维数组
	img.Data = make([][]T, height)
	for i := 0; i < height; i++ {
		img.Data[i] = data[i*width : (i+1)*width]
	}
	// 读取投影
	if isFlagTranOrProj {
		img.Tran = ds.GeoTransform()
		img.Proj = ds.Projection()
	}
	// 读取无效值
	if isFlagNodata {
		val, valid := band.NoDataValue()
		if valid {
			img.Nodata = val
		}
	}
	return &img, nil
}

// ImgWrite 写入单波段数据
func ImgWrite[T g.Numeric](outfile string, img *Img[T]) error {
	_ = fileutil.MkdirP(outfile)
	format := "GTiff"
	if outfile == "" {
		format = "MEM"
	}
	// 创建文件
	driver, err := gdal.GetDriverByName(format)
	if err != nil {
		return fmt.Errorf("GDAL TIFF driver not available")
	}

	width := len(img.Data[0])
	height := len(img.Data)

	var options []string
	if format == "GTiff" {
		options = append(options, "TILED=YES")
		options = append(options, "COMPRESS=LZW")
	}

	// 获取文件类型
	var t gdal.DataType
	switch reflect.TypeOf(img.Data).Elem().Elem().String() {
	case "float64":
		t = gdal.Float64
	case "float32":
		t = gdal.Float32
	case "int32":
		t = gdal.Int32
	case "uint32":
		t = gdal.UInt32
	case "int16":
		t = gdal.Int16
	case "uint16":
		t = gdal.UInt16
	case "uint8":
		t = gdal.Byte
	default:
		t = gdal.Float32
	}

	// 创建数据集
	ds := driver.Create(outfile, width, height, 1, t, options)
	defer ds.Close()

	if img.Proj != "" {
		proj := img.Proj
		if strings.HasPrefix(img.Proj, "EPSG:") {
			epsg, err := strconv.ParseInt(strings.TrimPrefix(img.Proj, "EPSG:"), 0, 0)
			glog.Info(epsg)
			if err != nil {
				return err
			}
			srs := gdal.CreateSpatialReference("")
			err = srs.FromEPSG(int(epsg))
			if err != nil {
				return err
			}
			proj, err = srs.ToWKT()
			if err != nil {
				return err
			}
		}
		err = ds.SetProjection(proj)
	}
	err = ds.SetGeoTransform(img.Tran)
	// 获取栅格波段
	band := ds.RasterBand(1)
	err = band.SetNoDataValue(img.Nodata)
	if err != nil {
		glog.Error(err.Error())
		return err
	}

	// 将二维数组转换为一维数据
	grid := make([]T, width*height)
	for i := 0; i < height; i++ {
		copy(grid[i*width:(i+1)*width], img.Data[i])
	}
	// 写入栅格数据
	err = band.IO(gdal.Write, 0, 0, width, height, grid, width, height, 0, 0)
	if err != nil {
		return fmt.Errorf("error writing raster data: %v", err)
	}
	return nil
}

// ImgWriteMem 写入单波段数据
func ImgWriteMem[T g.Numeric](img *Img[T]) (*gdal.Dataset, error) {
	format := "MEM"
	// 创建文件
	driver, err := gdal.GetDriverByName(format)
	if err != nil {
		return nil, fmt.Errorf("GDAL TIFF driver not available")
	}

	width := len(img.Data[0])
	height := len(img.Data)

	// 获取文件类型
	var t gdal.DataType
	switch reflect.TypeOf(img.Data).Elem().Elem().String() {
	case "float64":
		t = gdal.Float64
	case "float32":
		t = gdal.Float32
	case "int32":
		t = gdal.Int32
	case "uint32":
		t = gdal.UInt32
	case "int16":
		t = gdal.Int16
	case "uint16":
		t = gdal.UInt16
	case "uint8":
		t = gdal.Byte
	default:
		t = gdal.Float32
	}

	// 创建数据集
	ds := driver.Create("", width, height, 1, t, nil)

	if img.Proj != "" {
		proj := img.Proj
		if strings.HasPrefix(img.Proj, "EPSG:") {
			epsg, err := strconv.ParseInt(strings.TrimPrefix(img.Proj, "EPSG:"), 0, 0)
			glog.Info(epsg)
			if err != nil {
				return nil, err
			}
			srs := gdal.CreateSpatialReference("")
			err = srs.FromEPSG(int(epsg))
			if err != nil {
				return nil, err
			}
			proj, err = srs.ToWKT()
			if err != nil {
				return nil, err
			}
		}
		err = ds.SetProjection(proj)
	}

	err = ds.SetGeoTransform(img.Tran)
	// 获取栅格波段
	band := ds.RasterBand(1)
	err = band.SetNoDataValue(img.Nodata)
	if err != nil {
		return nil, err
	}
	// 将二维数组转换为一维数据
	grid := make([]T, width*height)
	for i := 0; i < height; i++ {
		copy(grid[i*width:(i+1)*width], img.Data[i])
	}
	// 写入栅格数据
	err = band.IO(gdal.Write, 0, 0, width, height, grid, width, height, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("error writing raster data: %v", err)
	}
	return &ds, nil
}

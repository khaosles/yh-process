package task

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/g"
	"github.com/lukeroth/gdal"
)

/*
   @File: stofs_test.go
   @Author: khaosles
   @Time: 2023/8/31 10:47
   @Desc:
*/

func TestName(t *testing.T) {
	path := "/Users/yherguot/doraemon/data/gfs/product/stofs/tif/2023-07-23_12-cwl.tif"
	data, err := ImgRead[float32](path, 88)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", data[637])
	fileList, _ := filepath.Glob("/Users/yherguot/doraemon/data/gfs/product/stofs/tif/*.tif")
	dat := make([][]float32, 1)
	sort.Strings(fileList)
	for _, file := range fileList {
		d, _ := ImgRead[float32](file, 88)
		dat = append(dat, d)
	}
}

func ImgRead[T g.Numeric](infile string, line int) ([]T, error) {
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

	// 读取栅格数据
	data := make([]T, width*1)
	err = band.IO(gdal.Read, 0, line, width, 1, data, width, 1, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("error reading raster data: %v", err)
	}

	return data, nil
}

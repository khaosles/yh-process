package util

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/mathutil"
	"github.com/khaosles/giz/slice"
	glog "github.com/khaosles/go-contrib/core/log"
	"gonum.org/v1/gonum/interp"
)

/*
   @File: interp.go
   @Author: khaosles
   @Time: 2023/7/25 12:22
   @Desc:
*/

func Inter(inpath string, element string, hours []string, nodata float64) error {

	// 获取原数据
	filelist := make([]string, len(hours))
	srcTime := make([]float64, len(hours))
	for i, hour := range hours {
		file := fileutil.Join(inpath, fmt.Sprintf("%s_%s.tif", element, hour))
		if !fileutil.IsFile(file) {
			return fmt.Errorf("file not found ===> %s", file)
		}
		f, _ := strconv.ParseFloat(hour, 0)
		srcTime[i] = f
		filelist[i] = file
	}
	// 数据排序
	sort.Strings(filelist)
	// 时间序列数组
	arr := make([][][]float64, len(filelist))

	var tran [6]float64 // 仿射参数
	var proj string     // 投影

	// 读取数据
	for i, file := range filelist {
		img, _ := ImgRead[float64](file, true, true)
		arr[i] = img.Data
		if proj == "" {
			proj = img.Proj
		}
		if tran[0] == 0 && tran[1] == 0 {
			tran = img.Tran
		}
	}
	// 高 宽
	height, width := len(arr[0]), len(arr[0][0])

	// 计算续约插值的时刻
	var interpTime []float64
	maxhour := mathutil.Max[float64](srcTime...)
	for i := 0; i < int(maxhour); i++ {
		if !slice.Contain[float64](srcTime, float64(i)) {
			interpTime = append(interpTime, float64(i))
		}
	}

	// 插值结果
	var result = make([][][]float64, len(interpTime))
	// 初始化切片
	for i := 0; i < len(interpTime); i++ {
		result[i] = make([][]float64, height)
		for j := 0; j < height; j++ {
			result[i][j] = make([]float64, width)
		}
	}

	// 插值
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// 是否全部为无效值
			allNodata := true
			// 单点数据
			srcData := make([]float64, len(srcTime))
			// 读取数据
			for k := 0; k < len(srcTime); k++ {
				val := arr[k][i][j]
				if !math.IsNaN(val) && val != nodata {
					allNodata = false
				}
				srcData[k] = val
			}
			if !allNodata {
				// 插值
				spline := interp.PiecewiseLinear{}
				err := spline.Fit(srcTime, srcData)
				if err != nil {
					return err
				}
				for idx, t := range interpTime {
					result[idx][i][j] = spline.Predict(t)
				}
			} else {
				// 全部为无效值则用无效值填补
				for idx := range interpTime {
					result[idx][i][j] = nodata
				}
			}
		}
	}

	// 保存数据
	wg := sync.WaitGroup{}
	limit := make(chan struct{}, 5)
	for i, t := range interpTime {
		wg.Add(1)
		go func(idx, t int) {
			limit <- struct{}{}
			defer wg.Done()
			defer func() { <-limit }()
			filename := fileutil.Join(inpath, fmt.Sprintf("%s_%03d.tif", fileutil.Basename(inpath), t))
			// 保存
			err := ImgWrite(filename, &Img[float64]{Data: result[idx], Tran: tran, Proj: proj, Nodata: nodata}) //Proj:   proj,
			if err != nil {
				glog.Error(err)
			}
		}(i, int(t))
	}
	wg.Wait()

	return nil
}

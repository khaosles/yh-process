package util

import (
	"fmt"
	"math"

	"github.com/khaosles/giz/fileutil"
)

/*
   @File: uvsd.go
   @Author: khaosles
   @Time: 2023/8/18 12:27
   @Desc:
*/

// Uv2sd uv 转为 风速风向
func Uv2sd(ufile, vfile, speedfile, dirfile string, nodata float32) error {
	// 校验文件是否存在
	if !fileutil.IsFile(ufile) || !fileutil.IsFile(vfile) {
		return fmt.Errorf("err")
	}
	//glog.Infof("u v -> s d ===> %s | %s", ufile, vfile)
	u, _ := ImgRead[float32](ufile, true, true)
	v, _ := ImgRead[float32](vfile, true, true)

	// 计算风速风向
	s, d, err := uv2sd(u.Data, v.Data, nodata)
	if err != nil {
		return err
	}
	// tiff
	if err := ImgWrite(speedfile, &Img[float32]{Data: s, Tran: u.Tran, Proj: ""}); err != nil {
		return err
	}
	if err := ImgWrite(dirfile, &Img[float32]{Data: d, Tran: u.Tran, Proj: ""}); err != nil {
		return err
	}
	return nil
}

func uv2sd(u, v [][]float32, nodata float32) ([][]float32, [][]float32, error) {
	if len(u) != len(v) || len(u[0]) != len(v[0]) {
		return nil, nil, fmt.Errorf("input arrays must have the same dimensions")
	}

	speed := make([][]float32, len(u))
	direction := make([][]float32, len(u))

	for i := 0; i < len(u); i++ {
		speed[i] = make([]float32, len(u[0]))
		direction[i] = make([]float32, len(u[0]))

		for j := 0; j < len(u[0]); j++ {
			if u[i][j] == nodata || v[i][j] == nodata {
				speed[i][j] = nodata
				direction[i][j] = nodata
			} else {
				speed[i][j] = float32(math.Sqrt(float64(u[i][j]*u[i][j] + v[i][j]*v[i][j])))
				direction[i][j] = float32(math.Atan2(float64(v[i][j]), float64(u[i][j])) * 180.0 / math.Pi)
			}
		}
	}
	return speed, direction, nil
}

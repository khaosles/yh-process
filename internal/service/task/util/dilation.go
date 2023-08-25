package util

import (
	"math"

	"github.com/khaosles/giz/g"
)

/*
   @File: dilation.go
   @Author: khaosles
   @Time: 2023/8/18 14:13
   @Desc:
*/

// Dilation 膨胀算法
func Dilation[T g.Numeric](input [][]T, windowSize int, nodata T) [][]T {
	height := len(input)
	width := len(input[0])
	result := make([][]T, height)
	for i := range result {
		result[i] = make([]T, width)
	}
	// 进行膨胀操作
	for i := 1; i < height-1; i++ {
		for j := 1; j < width-1; j++ {
			// 检查中心像素是否为无效值
			if input[i][j] == nodata || math.IsNaN(float64(input[i][j])) {
				// 3x3 区域内所有非无效值的均值
				sum := T(0)
				count := 0
				for m := -windowSize / 2; m <= windowSize/2; m++ {
					for n := -windowSize / 2; n <= windowSize/2; n++ {
						if input[i+m][j+n] != nodata && !math.IsNaN(float64(input[i+m][j+n])) {
							sum += input[i+m][j+n]
							count++
						}
					}
				}
				if count > 0 {
					// 使用均值填充中心像素
					result[i][j] = sum / T(count)
				} else {
					result[i][j] = nodata
				}
			} else {
				// 否则保留原数据
				result[i][j] = input[i][j]
			}
		}
	}
	return result
}

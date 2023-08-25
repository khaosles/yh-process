package util

/*
   @File: reserve.go
   @Author: khaosles
   @Time: 2023/8/22 23:34
   @Desc:
*/

func Reverse2d(slice [][]float32) [][]float32 {
	reversed := make([][]float32, len(slice))
	for i := 0; i < len(slice); i++ {
		reversed[i] = slice[len(slice)-1-i]
	}
	return reversed
}

func Reverse1d(slice []float64) []float64 {
	reversed := make([]float64, len(slice))
	for i := 0; i < len(slice); i++ {
		reversed[i] = slice[len(slice)-1-i]
	}
	return reversed
}

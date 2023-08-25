package util

import "time"

/*
   @File: times.go
   @Author: khaosles
   @Time: 2023/8/22 16:37
   @Desc:
*/

func GetYMDH(t time.Time) (string, string, string, string) {
	return t.Format("2006"), t.Format("01"), t.Format("02"), t.Format("15")
}

package gfs

import (
	"fmt"
	"strings"
	"time"

	"yh-process/internal/vo"
)

/*
   @File: util.go
   @Author: khaosles
   @Time: 2023/8/21 10:25
   @Desc:
*/

func GetDataTime() time.Time {
	// 获得当前utc
	curTime := time.Now().UTC()
	year := curTime.Year()
	month := curTime.Month()
	day := curTime.Day()
	hour := curTime.Hour()
	var date time.Time
	switch {
	case hour < 3:
		date = time.Date(year, month, day, 18, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
	case hour >= 3 && hour < 9:
		date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	case hour >= 9 && hour < 15:
		date = time.Date(year, month, day, 6, 0, 0, 0, time.UTC)
	case hour >= 15 && hour < 21:
		date = time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	default:
		date = time.Date(year, month, day, 18, 0, 0, 0, time.UTC)
	}
	return date
}

func TimeFormatURLInfo(URLInfos []*vo.URLInfo, t time.Time) {
	yyyymmdd := t.Format("20060102")
	hh := t.Format("15")
	for _, info := range URLInfos {
		info.Url = strings.Replace(info.Url, "{YYYYMMDD}", yyyymmdd, -1)
		info.Url = strings.Replace(info.Url, "{HH}", hh, -1)
		info.Filename = strings.Replace(info.Filename, "{HH}", hh, -1)
	}
}

func GetHours() []string {
	var hours []string
	for i := 0; i < 120; i++ {
		hours = append(hours, fmt.Sprintf("%03d", i))
	}
	for i := 120; i < 385; i += 3 {
		hours = append(hours, fmt.Sprintf("%03d", i))
	}
	return hours
}

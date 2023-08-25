package util

import (
	"strconv"
	"strings"

	"yh-process/internal/dao"
)

/*
   @File: threads.go
   @Author: khaosles
   @Time: 2023/8/23 11:22
   @Desc:
*/

func GetThreads(key string) []int {
	value := dao.SysConfigDao.GetValue(key)
	var threads []int
	for _, s := range strings.Split(value, ",") {
		i, _ := strconv.ParseInt(s, 0, 0)
		threads = append(threads, int(i))
	}
	return threads
}

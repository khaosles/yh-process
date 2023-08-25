package vo

import "yh-process/internal/model"

/*
   @File: redis.go
   @Author: khaosles
   @Time: 2023/8/21 10:27
   @Desc:
*/

// URLInfo 基础信息
type URLInfo struct {
	Id       string `json:"id"`
	Url      string `json:"url"`
	Filename string `json:"filename"`
}

type UV struct {
	Speed     *model.SysElementInfo
	Direction *model.SysElementInfo
	U         *model.SysElementInfo
	V         *model.SysElementInfo
}

package model

/*
   @File: sys_white_list.go
   @Author: khaosles
   @Time: 2023/8/22 17:38
   @Desc:
*/

type SysWhiteList struct {
	BaseModel
	Ip string `json:"ip" gorm:"type:varchar(100)"`
}

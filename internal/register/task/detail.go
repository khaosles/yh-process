package task

import "time"

/*
   @File: info.go
   @Author: khaosles
   @Time: 2023/8/4 00:17
   @Desc:
*/

type Detail struct {
	Name         string    `json:"name"`
	Spec         string    `json:"spec"`
	RegisterTime time.Time `json:"registerTime"`
	Description  string    `json:"description"`
}

package model

import (
	"time"

	"yh-process/internal/service/task/dag/node"
)

/*
   @File: df_task.go
   @Author: khaosles
   @Time: 2023/8/22 16:57
   @Desc:
*/

type Task struct {
	Id         string
	CreateTime time.Time `json:"createTime,omitempty" gorm:"autoCreateTime;column:create_time;type:timestamptz;comment:创建时间"` // 创建时间
	UpdateTime time.Time `json:"updateTime,omitempty" gorm:"autoUpdateTime;column:update_time;type:timestamptz;comment:更新时间"` // 更新时间
	Status     node.Status
	NodeCount  int
	Remarks    string `json:"remarks,omitempty" gorm:"column:remarks;default:null;comment:备注"` // 备注
}

func (Task) TableName() string {
	return "df_task"
}

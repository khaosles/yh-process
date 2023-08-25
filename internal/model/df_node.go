package model

import (
	"time"

	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	"yh-process/internal/service/task/dag/node"
)

/*
   @File: df_node.go
   @Author: khaosles
   @Time: 2023/8/22 16:58
   @Desc:
*/

type Node struct {
	Id         string                   `json:"id" gorm:"primaryKey;type:varchar(100);comment:id"`
	CreateTime time.Time                `json:"createTime,omitempty" gorm:"autoCreateTime;column:create_time;type:timestamptz;comment:创建时间"` // 创建时间
	UpdateTime time.Time                `json:"updateTime,omitempty" gorm:"autoUpdateTime;column:update_time;type:timestamptz;comment:更新时间"` // 更新时间
	Proxy      string                   `json:"proxy" gorm:"type:varchar(100);comment:代理名称"`
	Status     node.Status              `json:"status" gorm:"type:int;comment:节点状态"`
	Dependents *safe.OrderedSet[string] `json:"dependents" gorm:"type:varchar[];comment:依赖id"`
	Msg        *gjson.JsonObject        `json:"msg" gorm:"type:text;comment:结果消息"`
	Err        string                   `json:"err" gorm:"type:text;comment:错误信息"`
	TaskId     string                   `json:"taskId" gorm:"type:varchar(100);comment:任务id"`
	Remarks    string                   `json:"remarks,omitempty" gorm:"column:remarks;default:null;comment:备注"` // 备注
}

func (Node) TableName() string {
	return "df_node"
}

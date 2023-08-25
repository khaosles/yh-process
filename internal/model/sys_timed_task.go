package model

/*
   @File: sys_timed_task.go
   @Author: khaosles
   @Time: 2023-06-23 20:34:47
   @Desc: Automatic code generation
*/

type TimedTask struct {
	BaseModel
	TaskName string `json:"taskName,omitempty" gorm:"type:varchar(255);comment:任务名称"`
	Spec     string `json:"spec,omitempty"  gorm:"type:varchar(255);comment:定时任务表达式"`
}

func (TimedTask) TableName() string {
	return "sys_timed_task"
}

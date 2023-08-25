package dao

import (
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: task_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:12:13
   @Desc: auto-generated code
*/

var TaskDao *taskDao

func init() {
	var dao taskDao
	dao.DB = configure.GetDB()
	TaskDao = &dao
}

type taskDao struct {
	dfBaseDao[model.Task]
}

package dao

import (
	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: timed_task_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:08:12
   @Desc: auto-generated code
*/

var TimedTaskDao *timedTaskdao

func init() {
	var dao timedTaskdao
	dao.DB = configure.GetDB()
	TimedTaskDao = &dao
}

type timedTaskdao struct {
	g.BaseDao[model.TimedTask]
}

func (dao timedTaskdao) GetByTaskName(taskName string) string {
	entity, err := dao.SelectOne(&model.TimedTask{TaskName: taskName})
	if err != nil {
		glog.Error(err)
		return ""
	}
	return entity.Spec
}

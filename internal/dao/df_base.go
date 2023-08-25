package dao

import (
	"strconv"

	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
)

/*
   @File: df_base.go
   @Author: khaosles
   @Time: 2023/7/28 00:48
   @Desc:
*/

type dfDao[T any] interface {
	g.Dao[T]
	DeleteHistory(history int)
	CreateOrGet(*T)
}

type dfBaseDao[T any] struct {
	g.BaseDao[T]
}

func (dao dfBaseDao[T]) DeleteHistory(history int) {
	err := dao.DeleteHardByCondition(
		g.NewConditions().
			AndLessThan("create_time", "NOW( ) - INTERVAL '"+strconv.Itoa(history)+" DAYS'"),
	)
	if err != nil {
		glog.Error(err)
		return
	}
}

func (dao dfBaseDao[T]) CreateOrGet(obj *T) {
	err := dao.InsertOrSelect(obj)
	if err != nil {
		glog.Error(err)
	}
}

package dao

import (
	"strconv"

	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: sys_config_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:08:12
   @Desc: auto-generated code
*/

var SysConfigDao *sysConfigDao

func init() {
	var dao sysConfigDao
	dao.DB = configure.GetDB()
	SysConfigDao = &dao
}

type sysConfigDao struct {
	g.BaseDao[model.SysConfig]
}

func (dao sysConfigDao) GetValue(code string) string {
	obj, err := dao.SelectOneByConditions(
		&model.SysConfig{Code: code},
		g.NewConditions().Select("value"),
	)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return obj.Value
}

func (dao sysConfigDao) GetValueOfInteger(code string) int {
	value, err := strconv.Atoi(dao.GetValue(code))
	if err != nil {
		glog.Error(err)
	}
	return value
}

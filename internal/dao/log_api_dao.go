package dao

import (
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: log_api_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:18:12
   @Desc: auto-generated code
*/

var LogApiDao *logApiDao

func init() {
	var dao logApiDao
	dao.DB = configure.GetDB()
	LogApiDao = &dao
}

type logApiDao struct {
	dfBaseDao[model.LogApi]
}

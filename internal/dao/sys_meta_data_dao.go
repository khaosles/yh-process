package dao

import (
	"time"

	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: sys_meta_data_dao.go
   @Author: khaosles
   @Time: 2023-07-01 16:46:41
   @Desc: auto-generated code
*/

var SysMetaDataDao *sysMetaDataDao

func init() {
	var dao sysMetaDataDao
	dao.DB = configure.GetDB()
	SysMetaDataDao = &dao
}

type sysMetaDataDao struct {
	g.BaseDao[model.SysMetaData]
}

func (dao sysMetaDataDao) FindByQ(dataType, elemName string) *model.SysMetaData {
	metaData, err := dao.SelectOne(&model.SysMetaData{
		DataType: dataType,
		ElemName: elemName,
	})
	if err != nil {
		glog.Error(err)
		return nil
	}
	return metaData
}

func (dao sysMetaDataDao) UpdateVersion(data *model.SysMetaData, dataTime, reportTime time.Time, version int64) {
	err := dao.UpdateSelective(data, map[string]any{
		"data_time":   dataTime,
		"report_time": reportTime,
		"version":     version,
	})
	if err != nil {
		glog.Error(err)
	}
}

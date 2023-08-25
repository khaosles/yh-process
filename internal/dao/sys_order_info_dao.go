package dao

import (
	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
	"yh-process/internal/configure"
	"yh-process/internal/model"
)

/*
   @File: sys_order_info_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:08:12
   @Desc: auto-generated code
*/

var SysOrderInfoDao *sysOrderInfoDao

func init() {
	var dao sysOrderInfoDao
	dao.DB = configure.GetDB()
	SysOrderInfoDao = &dao
}

type sysOrderInfoDao struct {
	g.BaseDao[model.SysOrderInfo]
}

func (dao sysOrderInfoDao) GetOrderInfo(orderName string) *model.SysOrderInfo {
	orderModel, err := dao.SelectOne(&model.SysOrderInfo{OrderName: orderName})
	if err != nil {
		glog.Error(err)
		return nil
	}
	return orderModel
}

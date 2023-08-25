package dao

import (
	glog "github.com/khaosles/go-contrib/core/log"
	g "github.com/khaosles/go-contrib/gorm/orm"
	"yh-process/internal/configure"
	"yh-process/internal/model"
	"yh-process/internal/vo"
)

/*
   @File: sys_element_info_info_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:08:12
   @Desc: auto-generated code
*/

var SysElementInfoDao *sysElementInfoDao

func init() {
	var dao sysElementInfoDao
	dao.DB = configure.GetDB()
	SysElementInfoDao = &dao
}

type sysElementInfoDao struct {
	g.BaseDao[model.SysElementInfo]
}

func (dao sysElementInfoDao) GetElementInfoByOrderId(orderId string) []*model.SysElementInfo {
	// 查询订单对应产品
	ElementInfos, err := dao.SelectByCondition(
		g.NewConditions().
			AndEqualTo("sys_order_id", orderId).
			AndIsNull("relevance"),
	)
	if err != nil {
		glog.Error(err)
		return nil
	}
	return ElementInfos
}

func (dao sysElementInfoDao) GetWindyByOrderId(orderId string) []*vo.UV {
	// 查询所有需要转换的风速
	ElementInfos, err := dao.SelectByCondition(
		g.NewConditions().
			AndEqualTo("sys_order_id", orderId).
			AndIsNotNull("relevance").
			AndIsNotNull("match"),
	)
	if err != nil {
		glog.Error(err)
	}
	var uvs []*vo.UV
	// 查询对应的风向
	for _, speed := range ElementInfos {
		direction, err := dao.SelectById(speed.Match)
		if err != nil {
			glog.Error(err)
		}
		u, err := dao.SelectById(speed.Relevance)
		if err != nil {
			glog.Error(err)
		}
		v, err := dao.SelectById(direction.Relevance)
		if err != nil {
			glog.Error(err)
		}
		uv := &vo.UV{
			Speed:     speed,
			Direction: direction,
			U:         u,
			V:         v,
		}
		uvs = append(uvs, uv)
	}
	return uvs
}

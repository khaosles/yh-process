package dao

import (
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	"yh-process/internal/configure"
	"yh-process/internal/model"
	"yh-process/internal/service/task/dag/node"
)

/*
   @File: product_dao.go
   @Author: khaosles
   @Time: 2023-07-06 14:12:13
   @Desc: auto-generated code
*/

var NodeDao *nodeDao

func init() {
	var dao nodeDao
	dao.DB = configure.GetDB()
	NodeDao = &dao
}

type nodeDao struct {
	dfBaseDao[model.Node]
}

func (d *nodeDao) SaveNode(id, proxy string, status node.Status,
	dependents *safe.OrderedSet[string], msg *gjson.JsonObject, err string) {
	var n model.Node
	n.Id = id
	n.Proxy = proxy
	n.Status = status
	n.Dependents = dependents
	n.Msg = msg
	n.Err = err
	d.Save(&n)
	//d.DB.
}

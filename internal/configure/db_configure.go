package configure

import (
	"github.com/khaosles/go-contrib/gorm/pgsql"
	"gorm.io/gorm"
)

/*
   @File: db_configure.go
   @Author: khaosles
   @Time: 2023/6/24 10:49
   @Desc:
*/

func GetDB() *gorm.DB {
	//_ = pgsql.DB.AutoMigrate(&model.Order{})
	//_ = pgsql.DB.AutoMigrate(&model.Task{})
	//_ = pgsql.DB.AutoMigrate(&model.Product{})
	return pgsql.DB
}

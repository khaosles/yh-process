package model

import (
	"time"

	"github.com/khaosles/giz/gen"
	"gorm.io/gorm"
)

/*
   @File: base_model.go
   @Author: khaosles
   @Time: 2023/8/21 10:06
   @Desc:
*/

type BaseModel struct {
	Id         string         `json:"id" gorm:"primaryKey;column:id;type:varchar(32);comment:主键"`                                  // 主键ID
	CreateTime time.Time      `json:"createTime,omitempty" gorm:"autoCreateTime;column:create_time;type:timestamptz;comment:创建时间"` // 创建时间
	UpdateTime time.Time      `json:"updateTime,omitempty" gorm:"autoUpdateTime;column:update_time;type:timestamptz;comment:更新时间"` // 更新时间
	DeleteTime gorm.DeletedAt `json:"-" gorm:"index;column:delete_time;type:timestamptz;comment:删除时间"`
	Remarks    string         `json:"remarks,omitempty" gorm:"column:remarks;default:null;comment:备注"` // 备注
}

func (m *BaseModel) BeforeCreate(*gorm.DB) error {
	m.Id = gen.UuidNoSeparator()
	return nil
}

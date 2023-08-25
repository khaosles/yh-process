package model

/*
   @File: sys_config.go
   @Author: khaosles
   @Time: 2023-06-23 20:34:47
   @Desc: Automatic code generation
*/

type SysConfig struct {
	BaseModel
	Code  string `json:"code" gorm:"type:varchar(40);unique;index;comment:配置编码"`
	Value string `json:"value" gorm:"type:text;comment:配置值"`
	Note  string `json:"note" gorm:"type:varchar(100);comment:配置名称"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

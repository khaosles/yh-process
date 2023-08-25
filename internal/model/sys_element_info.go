package model

import "yh-process/internal/consts"

/*
   @File: sys_element_info.go
   @Author: khaosles
   @Time: 2023-06-23 20:34:47
   @Desc: Automatic code generation
*/

type SysElementInfo struct {
	BaseModel
	SysOrderId string      `json:"sysOrderId" gorm:"type:varchar(32);index;comment:订单文件配置id"`
	Elemname   string      `json:"elemName" gorm:"type:varchar(40);comment:要素名称"`
	Field      string      `json:"field" gorm:"type:varchar(40);comment:对应文件中字段名称"`
	Lev        string      `json:"lev" gorm:"default:NULL;type:varchar(40);comment:下载层级"`
	Var        string      `json:"var" gorm:"default:NULL;type:varchar(40);comment:下载变量名称"`
	Gain       float32     `json:"gain" gorm:"default:1;type:float;comment:系数"`
	Offset     float32     `json:"offset" gorm:"default:0;type:float;comment:偏移"`
	Interval   int         `json:"interval"  gorm:"default:0;type:int2;comment:数据时间间隔 0 前120逐小时 后面逐三小时 1 逐三天"`
	IsPng      consts.Bool `json:"isNumerical" gorm:"default:0;type:int2;comment:是否处理成数值产品数据 0不处理 1处理"`
	IsBin      consts.Bool `json:"isLayer" gorm:"default:0;type:int2;comment:处理成全球png图片 0不处理 1处理"`
	IsFill     int         `json:"isFill" gorm:"default:1;type:int2;comment:提取格点时是否填充缺失数据 0不填充 1填充"`
	Relevance  string      `json:"relevance" gorm:"default:NULL;type:varchar(32);comment:自关联的字段, 风速此字段关联U 风向此字段关联V"`
	Match      string      `json:"match" gorm:"default:NULL;type:varchar(32);comment:自关联的字段，风速关联风向 证明次这是一对风的属性，根据对应的UV值进行转换"`
}

func (SysElementInfo) TableName() string {
	return "sys_element_info"
}

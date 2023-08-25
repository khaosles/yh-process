package model

import (
	"time"
)

/*
   @File: sys_meta_data.go
   @Author: khaosles
   @Time: 2023/6/16 17:06
   @Desc: auto-generated code
*/

type SysMetaData struct {
	BaseModel
	DataType   string    `json:"dataType" gorm:"type:varchar(100);index;comment:数据类型"`
	ElemName   string    `json:"elemName" gorm:"type:varchar(100);index;comment:要素名称"`
	Version    int64     `json:"version" gorm:"type:int;comment:当前版本号"`
	DataTime   time.Time `json:"dataTime" gorm:"type:timestamptz;comment:数据生产时间"`
	ReportTime time.Time `json:"reportTime" gorm:"type:timestamptz;comment:数据起报时间"`
	FileName   string    `json:"fileName" gorm:"type:varchar(255);comment:文件名称"`
	LeftLon    float32   `json:"leftLon" gorm:"type:float4;comment:数据最小经度"`
	RightLon   float32   `json:"rightLon" gorm:"type:float4;comment:数据最大经度"`
	BottomLat  float32   `json:"bottomLat" gorm:"type:float4;comment:数据最小纬度"`
	TopLat     float32   `json:"topLat" gorm:"type:float4;comment:数据最大纬度"`
	GridSize   float32   `json:"gridSize" gorm:"type:float4;comment:格网大小"`
	Gain       float32   `json:"gain" gorm:"type:float;default:1;comment:数据系数"`
	Offset     float32   `json:"offset" gorm:"type:float;default:0;comment:数据偏移"`
	Length     int       `json:"length" gorm:"type:int;comment:数据长度"`
	FillValue  int       `json:"fillValue" gorm:"type:int;default:-9999;comment:数据无效值"`
	Bytes      int       `json:"bytes" gorm:"type:int;default:2;comment:单时刻数据位数"`
	Unit       string    `json:"unit" gorm:"type:varchar(40);comment:数据单位"`
}

func (SysMetaData) TableName() string {
	return "sys_meta_data"
}

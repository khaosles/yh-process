package model

/*
   @File: sys_order_info.go
   @Author: khaosles
   @Time: 2023-06-23 20:34:47
   @Desc: Automatic code generation
*/

type SysOrderInfo struct {
	BaseModel
	DataType  string `json:"dataType" gorm:"type:varchar(100);comment:类型"`
	OrderName string `json:"orderName" gorm:"type:varchar(20);unique;index;comment:订单名称"`
	FileName  string `json:"fileName" gorm:"type:varchar(100);comment:文件名称"`
	BaseUrl   string `json:"fmtUrl" gorm:"type:text;comment:带格式化字符串的下载url"`
}

func (SysOrderInfo) TableName() string {
	return "sys_order_info"
}

package model

import (
	"time"
)

/*
   @File: log_api.go
   @Author: khaosles
   @Time: 2023-06-23 20:34:56
   @Desc: Automatic code generation
*/

type LogApi struct {
	BaseModel
	IP           string        `json:"ip" gorm:"type:varchar(100);comment:请求ip"`                          // 请求ip
	Location     string        `json:"location" gorm:"type:varchar(100);comment:请求地址"`                    // 请求地址
	Method       string        `json:"method" gorm:"type:varchar(10);comment:请求方法"`                       // 请求方法
	Path         string        `json:"path" gorm:"type:varchar(100);comment:请求路径"`                        // 请求路径
	UserAgent    string        `json:"userAgent" gorm:"type:varchar(255);comment:代理"`                     // 代理
	Os           string        `json:"os" gorm:"type:varchar(10);comment:操作系统"`                           // 操作系统
	Browser      string        `json:"browser" gorm:"type:varchar(40);browser;comment:浏览器"`               // 浏览器
	ReqParam     string        `json:"reqParam" gorm:"type:text;comment:请求参数"`                            // 请求Body
	RespResult   string        `json:"respResult" gorm:"type:text;comment:响应Body"`                        // 响应Body
	Status       uint          `json:"status" gorm:"default:1;type:uint;comment:请求状态 1成功 2失败"`            // 请求状态
	Latency      time.Duration `json:"latency" gorm:"type:numeric(10,2);comment:延迟" swaggertype:"string"` // 延迟
	ErrorMessage string        `json:"errorMessage" gorm:"comment:错误信息"`                                  // 错误信息
}

func (LogApi) TableName() string {
	return "log_api"
}

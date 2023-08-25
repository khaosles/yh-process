package vo

import "time"

/*
   @File: update.go
   @Author: khaosles
   @Time: 2023/8/23 13:16
   @Desc:
*/

type UpdateVo struct {
	Filepath   string    `json:"filepath,omitempty"`
	Version    int64     `json:"version,omitempty"`
	ReportTime time.Time `json:"reportTime"`
	DataTime   time.Time `json:"dataTime"`
	DataType   string    `json:"dataType,omitempty"`
	Element    string    `json:"element,omitempty"`
}

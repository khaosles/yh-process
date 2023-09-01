package consts

import "time"

/*
   @File: consts.go
   @Author: khaosles
   @Time: 2023/7/11 10:06
   @Desc:
*/

const (
	SUFFIX        = ".yh.downloading"
	RETRY         = 3
	MAX_RUNTIME   = 5*time.Hour + 30*time.Minute
	TICK_INTERVAL = 5 * time.Minute
	TASK_TIMEOUT  = 5 * time.Minute
)

const (
	DAILY = 0
	THREE = 1
)

const Projection = ""

//const Projection = "EPSG:4326"

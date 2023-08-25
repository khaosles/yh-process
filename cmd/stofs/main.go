package main

import (
	timer2 "github.com/khaosles/giz/timer"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/model"
	"yh-process/internal/service/task/stofs"
)

/*
   @File: main.go
   @Author: khaosles
   @Time: 2023/8/25 10:35
   @Desc:
*/

func main() {
	timer := timer2.NewTimerTask()

	defer func() {
		timer.StopTask(consts.DataGfsStofs)
		timer.Close()
	}()

	timedTask, _ := dao.TimedTaskDao.SelectOne(&model.TimedTask{TaskName: consts.DataGfsStofs})
	_, _ = timer.AddTaskByFunc(consts.DataGfsStofs, timedTask.Spec, func() {
		stofs.ExecGFSStofsDAG()
	})
	glog.Infof("start cron => %s", consts.DataGfsStofs)

	timer.StartTask(consts.DataGfsStofs)

	select {}
}

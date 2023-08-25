package main

import (
	timer2 "github.com/khaosles/giz/timer"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/model"
	"yh-process/internal/service/clear"
	"yh-process/internal/service/task/gfs"
)

/*
   @File: main.go
   @Author: khaosles
   @Time: 2023/8/22 17:09
   @Desc:
*/

func main() {
	timer := timer2.NewTimerTask()

	defer func() {
		timer.StopTask(consts.DataGfsAtmos)
		timer.StopTask(consts.DataGfsWave)
		timer.StopTask(consts.ClearGfs)
		timer.Close()
	}()

	timedTask, _ := dao.TimedTaskDao.SelectOne(&model.TimedTask{TaskName: consts.DataGfsWave})
	_, _ = timer.AddTaskByFunc(consts.DataGfsAtmos, timedTask.Spec, func() {
		gfs.ExecGFSWaveDAG()
	})
	timer.StartTask(consts.DataGfsWave)
	glog.Infof("start cron => %s", consts.DataGfsWave)
	
	timedTask, _ = dao.TimedTaskDao.SelectOne(&model.TimedTask{TaskName: consts.DataGfsAtmos})
	_, _ = timer.AddTaskByFunc(consts.DataGfsAtmos, timedTask.Spec, func() {
		gfs.ExecGFSAtmosDAG()
	})
	timer.StartTask(consts.DataGfsAtmos)
	glog.Infof("start cron => %s", consts.DataGfsAtmos)

	timedTask, _ = dao.TimedTaskDao.SelectOne(&model.TimedTask{TaskName: consts.ClearGfs})
	_, _ = timer.AddTaskByFunc(consts.ClearGfs, timedTask.Spec, func() {
		clear.Clear("gfs")
	})

	timer.StartTask(consts.ClearGfs)
	glog.Infof("start cron => %s", consts.ClearGfs)

	select {}
}

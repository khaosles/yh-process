package gfs

import (
	"yh-process/internal/consts"
	wave2 "yh-process/internal/service/task/gfs/wave"
)

/*
   @File: gfs_wave_dag.go
   @Author: khaosles
   @Time: 2023/8/22 23:11
   @Desc:
*/

func ExecGFSWaveDAG() {
	ExecGFSDAG(consts.DataGfsWave, wave2.GenURL, wave2.NewGfsWaveTiffDataProcessor)
}

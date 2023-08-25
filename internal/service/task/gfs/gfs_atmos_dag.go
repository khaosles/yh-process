package gfs

import (
	"yh-process/internal/consts"
	atmos2 "yh-process/internal/service/task/gfs/atmos"
)

/*
   @File: gfs_atmos_dag.go
   @Author: khaosles
   @Time: 2023/8/22 00:42
   @Desc:
*/

func ExecGFSAtmosDAG() {
	ExecGFSDAG(consts.DataGfsAtmos, atmos2.GenURL, atmos2.NewGfsAtmosTiffDataProcessor)
}

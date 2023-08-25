package util

import (
	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/go-contrib/cmd"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
)

/*
   @File: wgrb.go
   @Author: khaosles
   @Time: 2023/8/18 11:00
   @Desc:
*/

func ExecGrib2Netcdf(src, dst string) error {
	_ = fileutil.MkdirP(dst)
	gribPath := dao.SysConfigDao.GetValue(consts.GribPath)
	var params []string
	params = append(params, src)
	params = append(params, "-netcdf")
	params = append(params, dst)
	err := cmd.Sync(gribPath, params...)
	if err != nil {
		return err
	}
	return nil
}

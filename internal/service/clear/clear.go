package clear

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/khaosles/giz/datetime"
	"github.com/khaosles/giz/fileutil"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/service/task/util"
)

/*
   @File: clear.go
   @Author: khaosles
   @Time: 2023/8/24 20:33
   @Desc:
*/

func Clear(dataType string) {
	rootPath := dao.SysConfigDao.GetValue(consts.RootPath)
	sourcePath := fileutil.Join(rootPath, consts.Source, dataType)
	del(sourcePath)
	productPath := fileutil.Join(rootPath, consts.Product, dataType)
	del(productPath)
}

func del(path string) {
	day := datetime.BeginOfDay(time.Now().Add(-time.Hour * 24))
	y, m, d, _ := util.GetYMDH(day)
	dirs, _ := filepath.Glob(fmt.Sprintf("%s/%s/%s/*", path, y, m))
	compare := fileutil.Join(path, y, m, d)
	for _, dir := range dirs {
		if dir <= compare {
			err := util.Delete(dir, 3)
			if err != nil {
				glog.Error(err)
			}
		}
	}
}

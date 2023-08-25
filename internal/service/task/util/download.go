package util

import (
	"errors"
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	"github.com/khaosles/giz/fileutil"
	glog "github.com/khaosles/go-contrib/core/log"
)

/*
   @File: download.go
   @Author: khaosles
   @Time: 2023/8/22 16:36
   @Desc:
*/

func Download(url, tmpfile, savefile string, retry int) error {
	glog.Debugf("Download file ===> %s", url)
	_ = fileutil.MkdirP(tmpfile)
	_ = fileutil.MkdirP(savefile)
	// 文件存在
	if fileutil.IsFile(tmpfile) {
		_ = fileutil.Rm(tmpfile)
	}
	// 下载器
	client := grab.NewClient()
	req, err := grab.NewRequest(fileutil.Dirname(tmpfile), url)
	if err != nil {
		return err
	}
	req.Filename = tmpfile
	req.SkipExisting = false
	req.NoResume = false
	// 执行下载
	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		return err
	}
	// 获取当前文件大小
	expectedSize := resp.Size()
	// 下载文件大小
	actualSize := fileutil.FileSize(tmpfile)
	if fileutil.IsFile(tmpfile) && (actualSize == expectedSize) {
		// 移动文件
		_ = fileutil.Move(tmpfile, savefile)
	} else if retry > 0 {
		retry--
		return Download(url, tmpfile, savefile, retry)
	} else {
		return errors.New(fmt.Sprintf("file download failed=> %s  all: %d download: %d", tmpfile, expectedSize, actualSize))
	}
	return nil
}

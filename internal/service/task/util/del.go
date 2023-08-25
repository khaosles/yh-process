package util

import (
	"os"
	"path/filepath"

	"github.com/khaosles/giz/fileutil"
)

/*
   @File: del.go
   @Author: khaosles
   @Time: 2023/8/23 15:35
   @Desc:
*/

// Delete 删除文件并删除{tier}层空文件夹
func Delete(path string, tier int) error {
	if tier < 0 {
		tier = 0
	}
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	for i := 0; i < tier; i++ {
		path = fileutil.Dirname(path)
		result := removeEmptyDir(path)
		if !result {
			break
		}
	}
	return nil
}

func removeEmptyDir(dir string) bool {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return false
	}
	if len(files) == 0 {
		err := os.Remove(dir)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

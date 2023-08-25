package download

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/khaosles/giz/assert"
	"github.com/khaosles/giz/fileutil"
)

/*
   @File: download.go
   @Author: khaosles
   @Time: 2023/8/23 11:36
   @Desc:
*/

func Download(c context.Context, ctx *app.RequestContext) {
	url := ctx.Query("url")
	assert.IsBlank(url, 40401, "参数不能为空")
	assert.IsNotFile(url, 40400, fmt.Sprintf("file not exist: %s", url))
	// file information
	filesize := fileutil.FileSize(url)
	// file header
	fileContentDisposition := "attachment;filename=\"" + fileutil.Basename(url) + "\""
	ctx.Header("Content-Disposition", fileContentDisposition)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Length", strconv.FormatInt(filesize, 10))
	ctx.File(url)
}

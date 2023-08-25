package gfs

import (
	"fmt"
	"strings"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	"yh-process/internal/consts"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
	"yh-process/internal/vo"
)

/*
   @File: data_fetcher.go
   @Author: khaosles
   @Time: 2023/8/22 23:13
   @Desc:
*/

type GfsDataFetcher struct {
	node2.BaseNode
}

func NewGfsDataFetcher(id, proxy string, cache node2.Cache) *GfsDataFetcher {
	var n GfsDataFetcher
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *GfsDataFetcher) Execute() *node2.Result {

	var result node2.Result
	var err error
	var tmpdir, sourcedir string
	var filename, tmpfile, sourcefile string
	var url string
	var uris []*vo.URLInfo
	var jsonObject = gjson.NewJsonObject()

	// 解析id
	idsplit := strings.Split(n.Id, ":")
	// 获取url信息
	uris, err = redisutil.GetOfURLInfo(fmt.Sprintf(consts.RedisUri, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 临时处理路径
	tmpdir, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisTmpDir, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 源数据路径
	sourcedir, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisSourceDir, idsplit[0]))
	for _, uri := range uris {
		// 下载地址
		url = strings.Replace(uri.Url, "{FFF}", idsplit[2], -1)
		// 文件名称
		filename = strings.Replace(uri.Filename, "{FFF}", idsplit[2], -1)
		// 临时下载路径
		tmpfile = fileutil.Join(tmpdir, filename)
		// 最终保存路径
		sourcefile = fileutil.Join(sourcedir, filename)
		if fileutil.IsFile(sourcefile) {
			jsonObject.Put(uri.Id, sourcefile)
			continue
		}
		// 下载
		err = util.Download(url, tmpfile, sourcefile, 3)
		if err != nil {
			goto returntag
		}
		jsonObject.Put(uri.Id, sourcefile)
	}

returntag:
	if err != nil {
		// 错误处理
		result.Success = false
		result.Err = err.Error()
		n.Status = node2.Failed
	} else {
		result.Success = true
		result.Msg = jsonObject
		n.Status = node2.Succeed
	}
	result.Id = n.Id
	n.Result = &result
	return &result
}

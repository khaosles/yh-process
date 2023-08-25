package base

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/model"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
)

/*
   @File: bin_processor.go
   @Author: khaosles
   @Time: 2023/8/20 00:12
   @Desc:
*/

// BinDataProcessor 生成bin模块
type BinDataProcessor struct {
	node2.BaseNode
}

func NewBinDataProcessor(id, proxy string, cache node2.Cache) *BinDataProcessor {
	var n BinDataProcessor
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *BinDataProcessor) Execute() *node2.Result {

	var err error
	var nodata float32
	var dataType string
	var metadata *model.SysMetaData
	var tifPath, binPath, binfile string
	var filelist []string
	var jsonObject = gjson.NewJsonObject()
	var result node2.Result

	glog.Infof("execute tif to bin ===> %s", n.Id)

	// 解析id
	idsplit := strings.Split(n.Id, ":")

	// 保存路径
	tifPath, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisTifDir, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 无效值
	nodata, err = redisutil.GetOfF32(fmt.Sprintf(consts.RedisNodata, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 数据类型
	dataType, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisDataType, idsplit[0]))
	if err != nil {
		goto returntag
	}
	metadata = dao.SysMetaDataDao.FindByQ(dataType, idsplit[2])

	// 保存路径
	binPath, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisBinDir, idsplit[0]))
	if err != nil {
		goto returntag
	}

	// 获取文件
	filelist, err = filepath.Glob(fileutil.Join(tifPath, idsplit[2], fmt.Sprintf("*%s*.tif", idsplit[2])))
	if err != nil {
		goto returntag
	}
	binfile = fileutil.Join(binPath, fmt.Sprintf("%s.bin", idsplit[2]))
	// 转为bin
	err = util.Convert(filelist, binfile, metadata.Gain, metadata.Offset, nodata)
	if err != nil {
		goto returntag
	}
	jsonObject.Put(idsplit[2], binfile)

returntag:
	if err != nil {
		// 失败
		result.Success = false
		result.Err = err.Error()
		n.Status = node2.Failed
	} else {
		// 成功
		result.Success = true
		result.Msg = jsonObject
		n.Status = node2.Succeed
	}
	result.Id = n.Id
	n.Result = &result
	return &result
}

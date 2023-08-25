package stofs

import (
	"fmt"
	"strings"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/model"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util/redisutil"
)

/*
   @File: tiff_processor.go
   @Author: khaosles
   @Time: 2023/8/25 00:06
   @Desc:
*/

type GfsStofsTiffDataProcessor struct {
	node2.BaseNode
}

func NewGfsStofsTiffDataProcessor(id, proxy string, cache node2.Cache) node2.Node {
	var n GfsStofsTiffDataProcessor
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *GfsStofsTiffDataProcessor) Execute() *node2.Result {

	var err error
	var ncFile, tifPath string
	var result node2.Result
	var nodata float32
	var jsonObject = gjson.NewJsonObject()
	var elements []*model.SysElementInfo

	glog.Infof("execute nc to tif ===> %s", n.Id)

	// 解析id
	idsplit := strings.Split(n.Id, ":")

	dependentId := n.Dependents.Values()[0]
	parent, ok := n.NodeCahce.Get(dependentId)
	if !ok {
		err = fmt.Errorf("node cache not exist ===> %s", dependentId)
		goto returntag
	}
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
	// nc文件路径
	ncFile, err = parent.GetResult().Msg.GetString("stofs")
	if err != nil {
		goto returntag
	}
	if !fileutil.IsFile(ncFile) {
		err = fmt.Errorf("file not exist ===> %s", ncFile)
		goto returntag
	}

	// 获取需要提取的要素
	elements, err = redisutil.GetOfElements(fmt.Sprintf(consts.RedisFieldBase, idsplit[0]))
	if err != nil {
		goto returntag
	}
	for _, element := range elements {
		if element.Elemname != idsplit[2] {
			continue
		}
		outpath := fileutil.Join(tifPath, element.Elemname)
		err = ToGrid(ncFile, outpath, element.Elemname, float64(element.Gain), float64(element.Offset), float64(nodata))
		if err != nil {
			goto returntag
		}
	}

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

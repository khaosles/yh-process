package base

import (
	"fmt"
	"sort"
	"strings"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/gjson"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
)

/*
   @File: interp_processor.go.go
   @Author: khaosles
   @Time: 2023/8/20 00:11
   @Desc:
*/

// InterpDataProcessor 时序插值模块
type InterpDataProcessor struct {
	node2.BaseNode
}

func NewInterpDataProcessor(id, proxy string, cache node2.Cache) *InterpDataProcessor {
	var n InterpDataProcessor
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *InterpDataProcessor) Execute() *node2.Result {

	var err error
	var nodata float32
	var tifPath string
	var jsonObject = gjson.NewJsonObject()
	var result node2.Result
	var hours []string

	glog.Infof("execute interpation ===> %s", n.Id)

	// 解析id
	idsplit := strings.Split(n.Id, ":")

	dependentIds := n.Dependents.Values()
	// 获取预测时刻
	for _, id := range dependentIds {
		hours = append(hours, strings.Split(id, ":")[2])
	}
	sort.Strings(hours)

	// 无效值
	nodata, err = redisutil.GetOfF32(fmt.Sprintf(consts.RedisNodata, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 保存路径
	tifPath, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisTifDir, idsplit[0]))
	if err != nil {
		goto returntag
	}
	// 时序插值
	err = util.Inter(fileutil.Join(tifPath, idsplit[2]), idsplit[2], hours, float64(nodata))

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

package base

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/redis"
	"yh-process/internal/consts"
	node2 "yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/util/redisutil"
	"yh-process/internal/vo"
)

/*
   @File: data_manager.go
   @Author: khaosles
   @Time: 2023/8/20 00:02
   @Desc:
*/

// DataManager 数据管理
type DataManager struct {
	node2.BaseNode
}

func NewDataManager(id, proxy string, cache node2.Cache) *DataManager {
	var n DataManager
	n.Id = id
	n.Proxy = proxy
	n.Status = node2.Ready
	n.NodeCahce = cache
	n.Dependents = safe.NewOrderedSet[string]()
	return &n
}

func (n *DataManager) Execute() *node2.Result {
	var err error
	var updateTimes, reportTimes string
	var filepath, dataType string
	var updateTime, reportTime time.Time
	var version int64
	var updateVo vo.UpdateVo
	var result node2.Result
	var msgByte []byte

	glog.Infof("execute data update ===> %s", n.Id)
	// 解析id
	idsplit := strings.Split(n.Id, ":")
	dependentId := n.GetDependents()[0]
	dependent, ok := n.NodeCahce.Get(dependentId)
	if !ok {
		err = fmt.Errorf("节点信息不存在 ===> %s", dependentId)
		goto returntag
	}
	if filepath, err = dependent.GetResult().Msg.GetString(idsplit[2]); err != nil {
		goto returntag
	}
	if updateTimes, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisUpdateTime, idsplit[0])); err != nil {
		goto returntag
	}
	if reportTimes, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisReportTime, idsplit[0])); err != nil {
		goto returntag
	}
	if version, err = redisutil.GetOfI64(fmt.Sprintf(consts.RedisDataVersion, idsplit[0])); err != nil {
		goto returntag
	}
	if dataType, err = redisutil.GetOfString(fmt.Sprintf(consts.RedisDataType, idsplit[0])); err != nil {
		goto returntag
	}
	if updateTime, err = time.Parse(time.DateTime, updateTimes); err != nil {
		goto returntag
	}
	if reportTime, err = time.Parse(time.DateTime, reportTimes); err != nil {
		goto returntag
	}

	updateVo = vo.UpdateVo{
		DataType:   dataType,
		DataTime:   updateTime,
		ReportTime: reportTime,
		Version:    version,
		Element:    idsplit[2],
		Filepath:   filepath,
	}
	if msgByte, err = sonic.Marshal(updateVo); err != nil {
		goto returntag
	}
	redis.Publish(consts.RedisUpdate, string(msgByte))

returntag:
	if err != nil {
		// 失败
		result.Success = false
		result.Err = err.Error()
		n.Status = node2.Failed
	} else {
		// 成功
		result.Success = true
		n.Status = node2.Succeed
	}
	result.Id = n.Id
	n.Result = &result
	return &result
}

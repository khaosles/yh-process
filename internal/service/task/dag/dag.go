package dag

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/khaosles/giz/safe"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/redis"
	"yh-process/internal/consts"
	node2 "yh-process/internal/service/task/dag/node"
)

/*
   @File: dag.go
   @Author: khaosles
   @Time: 2023/8/20 00:20
   @Desc:
*/

type DAG struct {

	// 节点
	nodes *safe.OrderedSet[string]
	// 超时时间
	timeout time.Duration
	// 重试间隔，每个次间隔，就会刷新执行失败的节点
	retry time.Duration
	// 图名称，每个图唯一
	Name string
	// 节点缓存
	NodeCache node2.Cache
	// 代理缓存
	ProxyCahce *node2.ProxyCache

	// 更新信号
	singnal chan struct{}

	// 上下文
	ctx context.Context
	// 结束函数
	cancel context.CancelFunc
}

func NewDAG(timeout, retry time.Duration, Name string) *DAG {
	ctx, cannel := context.WithCancel(context.Background())
	return &DAG{
		Name:       Name,
		NodeCache:  node2.NewLocalCache(),
		ProxyCahce: node2.NewProxyCache(),
		timeout:    timeout,
		retry:      retry,
		singnal:    make(chan struct{}, 1),
		nodes:      safe.NewOrderedSet[string](),
		ctx:        ctx,
		cancel:     cannel,
	}
}

// AddNode 新增节点
func (d *DAG) AddNode(id string, n node2.Node) {
	// 添加缓存
	d.NodeCache.Put(id, n)
	// 添加到node
	d.nodes.Add(id)
}

// DeleteNode 删除节点
func (d *DAG) DeleteNode(id string) {
	// 删除节点
	d.nodes.Remove(id)
	// 删除缓存
	d.NodeCache.Delete(id)
}

// AddEdge 添加边
func (d *DAG) AddEdge(tail, head string) error {
	// 校验节点
	_, ok := d.NodeCache.Get(tail)
	if !ok {
		return fmt.Errorf("tail node not exist")
	}
	headNode, ok := d.NodeCache.Get(head)
	if !ok {
		return fmt.Errorf("head node not exist")
	}

	// 节点添加依赖
	headNode.AddDependent(tail)
	return nil
}

func (d *DAG) Run() {

	// 错误重试
	go d.update()
	// 结果输出
	d.watch()

	d.singnal <- struct{}{}

	// 定时器
	overallTick := time.NewTicker(d.timeout)
	for {
		select {
		case <-overallTick.C:
			d.cancel()
		case <-d.ctx.Done():
			d.close()
			return
		case <-d.singnal:
			d.run()
		default:
			continue
		}
	}
}

func (d *DAG) run() {
	for _, id := range d.nodes.Values() {
		d.item(id)
	}
}

func (d *DAG) item(id string) {
	n, _ := d.NodeCache.Get(id)
	for _, i := range n.GetDependents() {
		d.item(i)
	}
	idsplit := strings.Split(id, ":")
	proxy, _ := d.ProxyCahce.Get(idsplit[1])
	proxy.Apply(n)
}

// watch 日志监测 gorountine
func (d *DAG) watch() {
	d.ProxyCahce.Foreach(func(proxy *node2.TaskProxy) {
		proxyCp := proxy
		go func(proxy *node2.TaskProxy) {
			for {
				select {
				case result := <-proxy.Rch:
					if result == nil {
						continue
					}
					if result.Success {
						glog.Infof("node { %s } run success", result.Id)
					} else {
						glog.Errorf("node { %s } run failed, err: %s", result.Id, result.Err)
					}
				case <-d.ctx.Done():
					proxy.Close()
					return
				}
			}
		}(proxyCp)
	})
}

func (d *DAG) close() {
	// 清理缓存
	d.clearRedis()
	d.NodeCache.Clear()
	d.ProxyCahce.Clear()
}

// clearRedis 删除redis缓存
func (d *DAG) clearRedis() {
	glog.Infof("clear redis ===> %s", d.Name)
	script := `
		local keys = redis.call('KEYS', ARGV[1])
		for i, key in ipairs(keys) do
			redis.call('DEL', key)
		end
		return keys
	`
	err := redis.Eval(script, nil, fmt.Sprintf("%s%s:*", consts.RedisPrefix, d.Name))
	if err != nil {
		glog.Error(err.Error())
	}
}

// update 更新与退出监测 gorountine
func (d *DAG) update() {
	ticker := time.NewTicker(d.retry)
	for {
		select {
		case <-ticker.C:
			isExit := true
			d.NodeCache.Foreach(func(n node2.Node) {
				if n.GetStatus().IsFailed() {
					n.SetStatus(node2.Ready)
				}
				if !n.GetStatus().IsSucceed() {
					isExit = false
				}
			})
			// 退出
			if isExit {
				d.cancel()
			} else {
				d.singnal <- struct{}{}
			}
		}
	}
}

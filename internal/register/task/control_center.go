package task

import (
	"context"
	"sync"

	"github.com/bytedance/sonic"
	r8 "github.com/go-redis/redis/v8"
	"github.com/khaosles/giz/safe"
	"github.com/khaosles/giz/timer"
	glog "github.com/khaosles/go-contrib/core/log"
	"github.com/khaosles/go-contrib/etcd"
	"github.com/khaosles/go-contrib/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"yh-process/internal/consts"
)

/*
   @File: control_center.go
   @Author: khaosles
   @Time: 2023/8/6 02:28
   @Desc:
*/

type ControlCenter interface {
	Run()
	Close()
}

func NewControlCenter() ControlCenter {
	return &controlCenter{target: consts.EtcdDataDownload}
}

type Option int

const (
	Add Option = iota
	Delete
)

type update struct {
	Op     Option
	Key    string
	Detail Detail
}

type wch chan []*update

type controlCenter struct {
	target string      // key
	tm     timer.Timer // 定时器
	wch    wch         // 监测数据通道
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	ts     *safe.Map[string, *r8.PubSub] // 记录存储
}

func (cc *controlCenter) Run() {
	cc.ctx, cc.cancel = context.WithCancel(context.Background())
	cc.tm = timer.NewTimerTask()
	cc.ts = safe.NewMap[string, *r8.PubSub]()
	cc.wch = make(wch, 5)
	// 获取所有注册信息
	resp, err := etcd.Client.Get(cc.ctx, cc.target, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		glog.Warn("unmarshal get failed", zap.Error(err))
	}
	// 添加定时任务
	for _, kv := range resp.Kvs {
		var detail Detail
		err := sonic.Unmarshal(kv.Value, &detail)
		if err != nil {
			glog.Error(err)
			return
		}
		cc.addTask(detail.Name, detail.Spec) // 添加定时任务
	}
	// 监测变化
	go cc.watch(resp.Header.Revision + 1)
	cc.wg.Add(1)
	// 更新
	go cc.update()
}

func (cc *controlCenter) watch(rev int64) {
	defer close(cc.wch)
	opts := []clientv3.OpOption{clientv3.WithRev(rev), clientv3.WithPrefix()}
	wch := etcd.Client.Watch(cc.ctx, cc.target, opts...) // 监测key
	for {
		select {
		case <-cc.ctx.Done():
			return
		case wresp, ok := <-wch: // key发送变化
			if !ok {
				glog.Warn("watch closed", zap.String("target", cc.target))
				return
			}
			if wresp.Err() != nil {
				glog.Warn("watch failed", zap.String("target", cc.target), zap.Error(wresp.Err()))
				return
			}
			// 变化事件
			deltaUps := make([]*update, 0, len(wresp.Events))
			for _, e := range wresp.Events {
				var op Option
				switch e.Type {
				// 新增
				case clientv3.EventTypePut:
					op = Add
				// 删除
				case clientv3.EventTypeDelete:
					op = Delete
				default:
					continue
				}
				// 解析值
				var detail Detail
				err := sonic.Unmarshal(e.Kv.Value, &detail)
				if err != nil {
					glog.Warn("unmarshal endpoint update failed", zap.String("key", string(e.Kv.Key)), zap.Error(err))
					continue
				}
				up := &update{Op: op, Key: string(e.Kv.Key), Detail: detail}
				deltaUps = append(deltaUps, up)
			}
			if len(deltaUps) > 0 {
				// 发送数据
				cc.wch <- deltaUps
			}
		}
	}
}

func (cc *controlCenter) update() {
	defer cc.wg.Done()
	for { // 监测节点更新
		select {
		case <-cc.ctx.Done():
			return
		case ups, ok := <-cc.wch:
			if !ok {
				return
			}
			for _, up := range ups {
				switch up.Op {
				case Add:
					cc.addTask(up.Detail.Name, up.Detail.Spec) // 添加定时任务
				case Delete:
					cc.rmTask(up.Detail.Name) // 删除定时任务
				default:
					continue
				}
			}
		}
	}
}

func (cc *controlCenter) Close() {
	cc.cancel()
	cc.ts.Range(func(key string, value *r8.PubSub) {
		cc.rmTask(key)
	})
	cc.wg.Wait()
}

func (cc *controlCenter) addTask(name, spec string) {
	// 添加定时任务
	_, err := cc.tm.AddTaskByFunc(name, spec, func() {
		// 消息发布
		redis.Publish(consts.RedisRelease+name, "start!")
		glog.Info("Exec => ", name)
	})
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Info("Add task => ", name)
	// 订阅消息回调
	pubSub := redis.Subscribe(consts.RedisCallback+name, func(s string) {
		glog.Info("Recive => ", s)
	})
	cc.ts.Put(name, pubSub)
}

func (cc *controlCenter) rmTask(name string) {
	// 删除记录
	pubSub, has := cc.ts.Get(name)
	if has {
		// 停止任务
		cc.tm.StopTask(name)
		// 清除任务
		cc.tm.Clear(name)
		// 取消回调订阅
		_ = pubSub.Unsubscribe(context.Background(), consts.RedisCallback+name)
		// 删除记录
		cc.ts.Del(name)
	}
	glog.Info("Remove task => ", name)
}

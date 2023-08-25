package task

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/khaosles/giz/g"
	"github.com/khaosles/go-contrib/etcd"
	"github.com/khaosles/go-contrib/redis"
	"yh-process/internal/consts"
	"yh-process/internal/model"
)

/*
   @File: register.go
   @Author: khaosles
   @Time: 2023/8/4 00:15
   @Desc:
*/

func Register(task *model.TimedTask, ttl int64, cb func(string)) error {
	key := fmt.Sprintf("%s/%s", consts.EtcdDataDownload, task.TaskName)
	// 订阅消息
	go redis.Subscribe(consts.RedisRelease+task.Id, cb)
	// 退出监测
	go g.Exit(func() {
		_ = etcd.Del(key)
	})
	// 格式化数据
	detail := Detail{
		Name:         task.TaskName,
		Spec:         task.Spec,
		RegisterTime: time.Now(),
	}
	// 数据编码
	data, err := sonic.Marshal(detail)
	if err != nil {
		return err
	}
	err = etcd.PutWithLease(key, string(data), ttl)
	if err != nil {
		return err
	}
	return nil
}

func Unregister(name string) {
	_ = etcd.Del(fmt.Sprintf("%s/%s", consts.EtcdDataDownload, name))
}

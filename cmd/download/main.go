package main

import (
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/registry/etcd"
	"github.com/khaosles/go-contrib/core/config"
	glog "github.com/khaosles/go-contrib/core/log"
	"yh-process/internal/consts"
	"yh-process/internal/middleware"
	"yh-process/internal/service/download"
)

/*
   @File: main.go
   @Author: khaosles
   @Time: 2023/8/23 11:47
   @Desc:
*/

func main() {
	etcdCfg := config.GCfg.Etcd
	r, err := etcd.NewEtcdRegistry(
		etcdCfg.Nodes,
		etcd.WithAuthOpt(etcdCfg.Username, etcdCfg.Password),
	)
	if err != nil {
		glog.Fatal(err)
	}
	addr := config.GetAddr(consts.ServerDownload)
	h := server.Default(
		server.WithIdleTimeout(time.Hour),
		server.WithHostPorts(addr),
		server.WithTransport(standard.NewTransporter),
		server.WithRegistry(r, &registry.Info{
			ServiceName: consts.ServerDownload,
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      10,
			Tags:        nil,
		}),
		server.WithStreamBody(true),
	)
	h.Use(middleware.GlobalErrorHander)
	h.GET("/download", download.Download)
	h.Spin()
}

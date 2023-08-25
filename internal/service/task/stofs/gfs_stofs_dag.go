package stofs

import (
	"fmt"
	"time"

	"github.com/khaosles/giz/fileutil"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/model"
	"yh-process/internal/service/task/base"
	"yh-process/internal/service/task/dag"
	"yh-process/internal/service/task/dag/node"
	"yh-process/internal/service/task/gfs"
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
	"yh-process/internal/vo"
)

/*
   @File: gfs_stofs_dag.go
   @Author: khaosles
   @Time: 2023/8/25 00:17
   @Desc:
*/

const (
	Timeout = 5 * time.Hour
	Retry   = 1 * time.Minute
)

func ExecGFSStofsDAG() {

	// 更新时间
	updateTime := gfs.GetDataTime()
	updateTime = time.Date(2023, 8, 21, 00, 0, 0, 0, time.UTC)
	// 生成图名称
	dagName := fmt.Sprintf("%s-%s", consts.DataGfsStofs, updateTime.Format("06010215"))

	urlInfos := GenURL()                        // url
	gfs.TimeFormatURLInfo(urlInfos, updateTime) // 格式化时间
	// 查询元数据
	order := dao.SysOrderInfoDao.GetOrderInfo(consts.DataGfsStofs)
	elements := dao.SysElementInfoDao.GetElementInfoByOrderId(order.Id)

	// 初始化
	nit(dagName, updateTime, elements, urlInfos)
	// 新建图
	dagGraph := dag.NewDAG(Timeout, Retry, dagName)
	threads := []int{2, 1, 1, 1}
	// 设置代理缓存
	dagGraph.ProxyCahce.Put(consts.Graph1, node.NewTaskProxy(threads[0]))
	dagGraph.ProxyCahce.Put(consts.Graph2, node.NewTaskProxy(threads[1]))
	dagGraph.ProxyCahce.Put(consts.Graph4, node.NewTaskProxy(threads[2]))
	dagGraph.ProxyCahce.Put(consts.Graph5, node.NewTaskProxy(threads[3]))

	// 构建图
	for _, dataType := range []string{"cwl", "htp"} {
		// 下载
		g1Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph1, dataType)
		dataFetcher := NewGfsStofsDataFetcher(g1Id, consts.Graph1, dagGraph.NodeCache)
		dagGraph.AddNode(g1Id, dataFetcher)
		g2Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph2, dataType)

		// 提取tiff
		tiffDataProcessor := NewGfsStofsTiffDataProcessor(g2Id, consts.Graph2, dagGraph.NodeCache)
		dagGraph.AddNode(g2Id, tiffDataProcessor)
		_ = dagGraph.AddEdge(g1Id, g2Id)

		// 转为bytes
		g4Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph4, dataType)
		binProcessor := base.NewBinDataProcessor(g4Id, consts.Graph4, dagGraph.NodeCache)
		dagGraph.AddNode(g4Id, binProcessor)
		_ = dagGraph.AddEdge(g2Id, g4Id)

		// 更新
		g5Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph5, dataType)
		dataManager := base.NewDataManager(g5Id, consts.Graph5, dagGraph.NodeCache)
		dagGraph.AddNode(g5Id, dataManager)
		_ = dagGraph.AddEdge(g4Id, g5Id)
	}

	dagGraph.Run()
}

func nit(dagName string, updateTime time.Time, elements []*model.SysElementInfo, urlInfos []*vo.URLInfo) {
	rootPath := dao.SysConfigDao.GetValue(consts.RootPath)
	nodata := dao.SysConfigDao.GetValue(consts.Nodata)
	// 文件下载处理路径
	y, m, d, h := util.GetYMDH(updateTime)
	sourcePath := fileutil.Join(rootPath, consts.Source, "gfs", y, m, d, h)
	tmpPath := fileutil.Join(rootPath, consts.Tmp, "gfs")
	tifPath := fileutil.Join(rootPath, consts.Product, "gfs", y, m, d, h, "tif")
	binPath := fileutil.Join(rootPath, consts.Product, "gfs", y, m, d, h, "bin")

	// 存入redis
	_ = redisutil.SetOfURLInfo(fmt.Sprintf(consts.RedisUri, dagName), urlInfos, consts.ExpireTime)
	_ = redisutil.SetOfElements(fmt.Sprintf(consts.RedisFieldBase, dagName), elements, consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisUpdateTime, dagName), updateTime.Format(time.DateTime), consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisReportTime, dagName), updateTime.Add(-time.Hour*5).Format(time.DateTime), consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisDataVersion, dagName), fmt.Sprintf("%d", updateTime.Unix()), consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisDataType, dagName), consts.TypeGfs, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisNodata, dagName), nodata, consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisSourceDir, dagName), sourcePath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisTmpDir, dagName), tmpPath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisTifDir, dagName), tifPath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisBinDir, dagName), binPath, consts.ExpireTime)
}

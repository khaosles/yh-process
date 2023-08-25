package gfs

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
	"yh-process/internal/service/task/util"
	"yh-process/internal/service/task/util/redisutil"
	"yh-process/internal/vo"
)

/*
   @File: gfs_dag.go
   @Author: khaosles
   @Time: 2023/8/23 10:50
   @Desc:
*/

const (
	Timeout = 5 * time.Hour
	Retry   = 1 * time.Minute
)

func ExecGFSDAG(
	baseName string, genUrlCallBack func() []*vo.URLInfo,
	tiffProcessorCallBack func(string, string, node.Cache) node.Node) {

	// 更新时间
	updateTime := GetDataTime()
	//updateTime = time.Date(2023, 8, 21, 00, 0, 0, 0, time.UTC)
	// 预测的时刻
	hours := GetHours()
	// 生成图名称
	dagName := fmt.Sprintf("%s-%s", baseName, updateTime.Format("06010215"))

	urlInfos := genUrlCallBack()            // url
	TimeFormatURLInfo(urlInfos, updateTime) // 格式化时间
	// 查询元数据
	order := dao.SysOrderInfoDao.GetOrderInfo(baseName)
	elements := dao.SysElementInfoDao.GetElementInfoByOrderId(order.Id)
	uvs := dao.SysElementInfoDao.GetWindyByOrderId(order.Id)

	// 统计所要要转为bin的要素
	var bins []*model.SysElementInfo
	for _, element := range elements {
		if element.IsBin.True() {
			bins = append(bins, element)
		}
	}
	for _, uv := range uvs {
		if uv.Speed.IsBin.True() {
			bins = append(bins, uv.Speed)
		}
		if uv.Direction.IsBin.True() {
			bins = append(bins, uv.Direction)
		}
	}

	// 初始化
	nit(dagName, updateTime, uvs, elements, urlInfos)
	// 新建图
	dagGraph := dag.NewDAG(Timeout, Retry, dagName)
	threads := util.GetThreads(consts.ThreadGfs)
	if len(threads) != 5 {
		threads = []int{5, 5, 1, 1, 1}
	}
	// 设置代理缓存
	dagGraph.ProxyCahce.Put(consts.Graph1, node.NewTaskProxy(threads[0]))
	dagGraph.ProxyCahce.Put(consts.Graph2, node.NewTaskProxy(threads[1]))
	dagGraph.ProxyCahce.Put(consts.Graph3, node.NewTaskProxy(threads[2]))
	dagGraph.ProxyCahce.Put(consts.Graph4, node.NewTaskProxy(threads[3]))
	dagGraph.ProxyCahce.Put(consts.Graph5, node.NewTaskProxy(threads[4]))

	// 构建图
	var g2List []string
	for _, hour := range hours {
		// 下载
		g1Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph1, hour)
		dataFetcher := NewGfsDataFetcher(g1Id, consts.Graph1, dagGraph.NodeCache)
		dagGraph.AddNode(g1Id, dataFetcher)
		g2Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph2, hour)

		// 提取tiff
		tiffDataProcessor := tiffProcessorCallBack(g2Id, consts.Graph2, dagGraph.NodeCache)
		dagGraph.AddNode(g2Id, tiffDataProcessor)
		_ = dagGraph.AddEdge(g1Id, g2Id)
		g2List = append(g2List, g2Id)
	}

	for _, element := range bins {
		// 时序插值
		g3Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph3, element.Elemname)
		interpProocessor := base.NewInterpDataProcessor(g3Id, consts.Graph3, dagGraph.NodeCache)
		dagGraph.AddNode(g3Id, interpProocessor)
		for _, g2 := range g2List {
			_ = dagGraph.AddEdge(g2, g3Id)
		}

		// 转为bytes
		g4Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph4, element.Elemname)
		binProcessor := base.NewBinDataProcessor(g4Id, consts.Graph4, dagGraph.NodeCache)
		dagGraph.AddNode(g4Id, binProcessor)
		_ = dagGraph.AddEdge(g3Id, g4Id)

		// 更新
		g5Id := fmt.Sprintf("%s:%s:%s", dagName, consts.Graph5, element.Elemname)
		dataManager := base.NewDataManager(g5Id, consts.Graph5, dagGraph.NodeCache)
		dagGraph.AddNode(g5Id, dataManager)
		_ = dagGraph.AddEdge(g4Id, g5Id)
	}

	dagGraph.Run()
}

func nit(dagName string, updateTime time.Time, uvs []*vo.UV, elements []*model.SysElementInfo, urlInfos []*vo.URLInfo) {
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
	_ = redisutil.SetOfUV(fmt.Sprintf(consts.RedisFieldUV, dagName), uvs, consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisUpdateTime, dagName), updateTime.Format(time.DateTime), consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisReportTime, dagName), updateTime.Format(time.DateTime), consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisDataVersion, dagName), fmt.Sprintf("%d", updateTime.Unix()), consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisDataType, dagName), consts.TypeGfs, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisNodata, dagName), nodata, consts.ExpireTime)

	_ = redisutil.Set(fmt.Sprintf(consts.RedisSourceDir, dagName), sourcePath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisTmpDir, dagName), tmpPath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisTifDir, dagName), tifPath, consts.ExpireTime)
	_ = redisutil.Set(fmt.Sprintf(consts.RedisBinDir, dagName), binPath, consts.ExpireTime)
}

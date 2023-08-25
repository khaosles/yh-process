package atmos

import (
	"fmt"

	"github.com/khaosles/giz/slice"
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/vo"
)

/*
   @File: genurl.go
   @Author: khaosles
   @Time: 2023/8/21 10:13
   @Desc:
*/

func GenURL() []*vo.URLInfo {
	// 订单信息
	orderInfo := dao.SysOrderInfoDao.GetOrderInfo(consts.DataGfsAtmos)
	// 要素信息
	elementsInfo := dao.SysElementInfoDao.GetElementInfoByOrderId(orderInfo.Id)
	// 基础url
	baseUrl := orderInfo.BaseUrl

	// 添加下载变量
	var urlTmp []string
	for _, product := range elementsInfo {
		if !slice.Contain(urlTmp, product.Var) {
			urlTmp = append(urlTmp, product.Var)
		}
		if !slice.Contain(urlTmp, product.Lev) {
			urlTmp = append(urlTmp, product.Lev)
		}
	}

	// 拼接路径
	for _, tmp := range urlTmp {
		baseUrl += fmt.Sprintf("&%v=on", tmp)
	}
	return []*vo.URLInfo{
		{
			Id:       "gribfile",
			Url:      baseUrl,
			Filename: orderInfo.FileName,
		},
	}
}

package stofs

import (
	"yh-process/internal/consts"
	"yh-process/internal/dao"
	"yh-process/internal/vo"
)

/*
   @File: genurl.go
   @Author: khaosles
   @Time: 2023/8/21 10:33
   @Desc:
*/

func GenURL() []*vo.URLInfo {
	// 订单信息
	orderInfo := dao.SysOrderInfoDao.GetOrderInfo(consts.DataGfsStofs)

	//cwlUrl := strings.Replace(orderInfo.BaseUrl, "{TTT}", "cwl", -1)
	//cwlFilename := strings.Replace(orderInfo.FileName, "{TTT}", "cwl", -1)
	//htpUrl := strings.Replace(orderInfo.BaseUrl, "{TTT}", "htp", -1)
	//htpFilename := strings.Replace(orderInfo.FileName, "{TTT}", "htp", -1)

	return []*vo.URLInfo{
		{"stofs", orderInfo.BaseUrl, orderInfo.FileName},
	}
}

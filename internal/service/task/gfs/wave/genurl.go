package wave

import (
	"strings"

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
	orderInfo := dao.SysOrderInfoDao.GetOrderInfo(consts.DataGfsWave)

	url1 := strings.Replace(orderInfo.BaseUrl, "{RRRR}", "0p16", -1)
	url2 := strings.Replace(orderInfo.BaseUrl, "{RRRR}", "0p25", -1)
	filename1 := strings.Replace(orderInfo.FileName, "{RRRR}", "0p16", -1)
	filename2 := strings.Replace(orderInfo.FileName, "{RRRR}", "0p25", -1)

	return []*vo.URLInfo{
		{
			"0p16", url1, filename1,
		},
		{
			"0p25", url2, filename2,
		},
	}
}

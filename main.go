package main

import (
	"strings"

	"yh-process/internal/service/task/gfs"
	"yh-process/internal/service/task/util"
)

/*
   @File: main.go
   @Author: khaosles
   @Time: 2023/8/24 09:57
   @Desc:
*/

func main() {

	url := "https://nomads.ncep.noaa.gov/cgi-bin/filter_gfs_0p25_1hr.pl?dir=%2Fgfs.20230829%2F00%2Fatmos&file=gfs.t00z.pgrb2.0p25.f{FFF}&var_UGRD=on&var_VGRD=on&lev_10_m_above_ground=on"
	baseName := "/Users/yherguot/doraemon/data/gfs/tmp/gfs.t00z.pgrb2.0p25.f{FFF}"

	hours := gfs.GetHours()
	for _, hour := range hours {
		baseUrl := strings.Replace(url, "{FFF}", hour, -1)
		fileName := strings.Replace(baseName, "{FFF}", hour, -1)
		util.Download(baseUrl, fileName+".tmp", fileName, 3)
	}

	//now := time.Now()
	//filepath := "/Users/yherguot/doraemon/data/yh/stofs_2d_glo.t00z.fields.cwl.nc"
	//outpath := "/Users/yherguot/doraemon/data/gfs/product/stofs/cwl1"
	//stofs.ToGrid(filepath, outpath, "cwl", 100, 0, -9999)
	//glog.Infof("t: %f", time.Now().Sub(now).Seconds())
	//stofs.ExecGFSStofsDAG()
	//path := "/Users/yherguot/doraemon/data/gfs/product/gfs/2023/08/21/00/tif/cwl/*.tif"
	//out := "/Users/yherguot/doraemon/data/gfs/product/gfs/2023/08/21/00/tif/cwl.bin"
	//fmt.Println(path)
	//fp, err := os.Open(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//bytes, err := io.ReadAll(fp)
	//fmt.Println(len(bytes))
	//for i := 0; i < len(bytes); i += 2 {
	//	val, _ := util.ParseByte(bytes[i : i+2])
	//	fmt.Println(val)
	//}
	//matches, _ := filepath.Glob(path)
	//sort.Strings(matches)
	//data := make([][][]float32, len(matches))
	//for i, match := range matches {
	//	fmt.Println(i)
	//	img, err := util.ImgRead[float32](match, false, false)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	data[i] = img.Data
	//}
	//fmt.Println(data)
	//bs := util.ToBytes(data, 100, 0, -9999)
	//err := os.WriteFile(out, bs, os.ModePerm)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

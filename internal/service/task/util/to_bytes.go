package util

import (
	"bytes"
	"encoding/binary"
	"os"
	"sort"

	"github.com/khaosles/giz/fileutil"
	"github.com/khaosles/giz/g"
)

/*
   @File: to_bytes.go
   @Author: khaosles
   @Time: 2023/8/6 18:35
   @Desc:
*/

func Convert(filelist []string, outfile string, gain, offset, nodata float32) error {
	sort.Strings(filelist)
	_ = fileutil.MkdirP(outfile)

	var data [][][]float32
	for _, file := range filelist {
		img, err := ImgRead[float32](file, false, false)
		if err != nil {
			return err
		}
		data = append(data, img.Data)
	}
	bs := ToBytes(data, gain, offset, nodata)
	return os.WriteFile(outfile, bs, 0777)
}

// ToBytes 数组转 bytes
// data 维度 时间 纬度 经度
// gain 缩放系数
// offset 偏移系数
// 计算公式: dst = src * gain + offset
// nodata 数据无效值 无效值处不进行缩放偏移
func ToBytes[T g.Numeric](data [][][]T, gain, offset, nodata T) (bs []byte) {
	for i := 0; i < len(data); i++ { // 时间
		for j := 0; j < len(data[0]); j++ { // 纬度
			for k := 0; k < len(data[0][0]); k++ { // 经度
				v := data[i][j][k]
				if v != nodata {
					v = v/gain - offset // 进行缩放偏移
				}
				bs = append(bs, toByte(int16(v))...) // 转byte
			}
		}
	}
	return
}

func toByte(val int16) []byte {
	// 创建一个长度为 2 的 byte 数组
	byteArray := make([]byte, 2)
	// 将 int16 转换为 byte 数组
	binary.LittleEndian.PutUint16(byteArray, uint16(val))
	return byteArray
}

func ParseByte(bs []byte) (val int16, err error) {
	err = binary.Read(bytes.NewReader(bs), binary.LittleEndian, &val)
	return
}

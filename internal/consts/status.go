package consts

/*
   @File: status.go
   @Author: khaosles
   @Time: 2023/8/6 17:56
   @Desc:
*/

type Status int8

const (
	Failed  Status = iota - 1 // 失败
	Ing                       // 进行中
	Succeed                   // 成功
)

func (s Status) IsFailed() bool {
	return s == Failed
}

func (s Status) IsING() bool {
	return s == Ing
}

func (s Status) IsSucceed() bool {
	return s == Succeed
}

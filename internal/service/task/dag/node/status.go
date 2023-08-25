package node

/*
   @File: status.go
   @Author: khaosles
   @Time: 2023/8/19 17:28
   @Desc:
*/

// Status Node processing status
type Status uint

const (
	// Ready Pending
	Ready Status = iota

	// Running In progress
	Running

	// Succeed Successfully processed
	Succeed

	// Failed Processing failed
	Failed
)

func (s Status) IsReady() bool {
	return s == Ready
}

func (s Status) IsRunning() bool {
	return s == Running
}

func (s Status) IsSucceed() bool {
	return s == Succeed
}

func (s Status) IsFailed() bool {
	return s == Failed
}

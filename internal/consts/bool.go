package consts

/*
   @File: bool.go
   @Author: khaosles
   @Time: 2023/8/22 16:49
   @Desc:
*/

type Bool uint

const (
	False Bool = iota
	True
)

func (b Bool) True() bool {
	return b == True
}

func (b Bool) False() bool {
	return b == False
}

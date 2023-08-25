package node

import (
	"github.com/panjf2000/ants"
)

/*
   @File: task_proxy.go
   @Author: khaosles
   @Time: 2023/8/19 23:29
   @Desc:
*/

// TaskProxy is used to control the number of concurrent threads for nodes and execute them.
type TaskProxy struct {

	// Size Number of parallel executions
	Size int

	// pool Goroutine pool
	pool *ants.Pool

	// Rch is the channel used for result transmission.
	Rch chan *Result
}

// NewTaskProxy creates a new NodeProxy instance with the specified size.
// If size is greater than 1, it creates a pool; otherwise, no pool is created.
func NewTaskProxy(size int) *TaskProxy {
	var pool *ants.Pool
	if size > 1 {
		pool, _ = ants.NewPool(size)
	}
	return &TaskProxy{
		Size: size,
		pool: pool,
		Rch:  make(chan *Result, 100),
	}
}

// Apply automatically determines whether to execute node's validation
// and processing in parallel or single-threaded based on the pool.
func (p *TaskProxy) Apply(n Node) {

	if p.pool != nil {
		// Execute using goroutines
		_ = p.pool.Submit(func() {
			result := p.apply(n)
			p.Rch <- result
		})
	} else {
		// Execute using a single thread
		result := p.apply(n)
		p.Rch <- result
	}
}

// apply executes the validation and processing of the node.
func (p *TaskProxy) apply(n Node) *Result {
	var result *Result

	// 校验父节点是否处理完成或者该节点是否处于准备状态
	if !n.CheckDependents() || !n.GetStatus().IsReady() {
		return nil
	}

	// If validation fails, exit processing and return the error.
	if result = n.Validate(); result != nil && !result.Success {
		return result
	}

	// 设置节点状态
	n.SetStatus(Running)

	// Execute the processing.
	return n.Execute()
}

func (p *TaskProxy) Close() {
	close(p.Rch)
}

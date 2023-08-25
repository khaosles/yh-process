package node

/*
   @File: node.go
   @Author: khaosles
   @Time: 2023/8/19 16:44
   @Desc:
*/

// Node is the interface that any node in the DAG graph must implement.
// Nodes communicate with each other using Redis for message passing.
type Node interface {

	// Validate performs validation on the data received by the node.
	Validate() *Result

	// Execute carries out the execution process.
	Execute() *Result

	// GetId retrieves the node's id.
	GetId() string

	// GetProxy retrieves the node's proxy.
	GetProxy() string

	// GetResult retrieves the execution result of the node.
	GetResult() *Result

	// GetDependents retrieves the node's dependencies.
	GetDependents() []string

	// CheckDependents 校验依赖是否准备完整
	CheckDependents() bool

	// AddDependent adds a node's dependency.
	AddDependent(string)

	// GetStatus retrieves the node's status.
	GetStatus() Status

	// SetStatus set the node's status.
	SetStatus(Status)

	// ToString converts node information to a string.
	ToString() string
}

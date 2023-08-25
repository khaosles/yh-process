package node

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/khaosles/giz/safe"
)

/*
   @File: base_node.go
   @Author: khaosles
   @Time: 2023/8/19 16:59
   @Desc:
*/

// BaseNode contains the basic attributes that every node must have.
// Each struct based on this structure needs to implement the `Execute` method,
// and the `Validate` method can be implemented as needed.
// In a distributed data processing framework, the node's key in Redis is {prefix}:{dag_name}:{id}.
type BaseNode struct {

	// Id is the unique identifier for the node. {dag}:{node}:{type}
	Id string

	// Proxy Name of the proxy being processed
	Proxy string

	// Status represents the node's status.
	// There are four possible states: `Ready`, `Running`, `Succeed`, and `Failed`.
	// Ready: Indicates that the node is ready to trigger the `Validate` and `Execute` methods.
	// Running: Indicates that the node is currently running.
	// Succeed: Indicates that the node has successfully completed its execution.
	// Failed: Indicates that the node's execution has failed and needs to be reset to the `Ready` state for processing.
	Status Status

	// Dependents is a list of node IDs that this node depends on.
	Dependents *safe.OrderedSet[string]

	// Result is used to store the node's execution status and data information.
	// The `Validate` method should update this if it fails
	// The `Execute` method should update this regardless of success or failure.
	Result *Result

	// NodeCahce
	NodeCahce Cache
}

func (n *BaseNode) Validate() *Result {
	return nil
}

// Execute !!! This method needs to be specifically implemented.
func (n *BaseNode) Execute() *Result {
	panic(fmt.Errorf("`Execute` not implemented"))
}

func (n *BaseNode) ToString() string {
	bytes, err := sonic.Marshal(n)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (n *BaseNode) CheckDependents() bool {
	for _, id := range n.Dependents.Values() {
		node, ok := n.NodeCahce.Get(id)
		if !ok {
			return false
		}
		if !node.GetStatus().IsSucceed() {
			return false
		}
	}
	return true
}

func (n *BaseNode) GetResult() *Result {
	return n.Result
}

func (n *BaseNode) AddDependent(id string) {
	n.Dependents.Add(id)
}

func (n *BaseNode) GetDependents() []string {
	return n.Dependents.Values()
}

func (n *BaseNode) GetId() string {
	return n.Id
}

func (n *BaseNode) GetProxy() string {
	return n.Proxy
}

func (n *BaseNode) GetStatus() Status {
	return n.Status
}

func (n *BaseNode) SetStatus(status Status) {
	n.Status = status
}

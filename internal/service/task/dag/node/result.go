package node

import "github.com/khaosles/giz/gjson"

/*
   @File: result.go
   @Author: khaosles
   @Time: 2023/8/19 17:22
   @Desc:
*/

// Result represents the execution result of each node.
type Result struct {

	// Id represents the unique identifier of the node.
	Id string `json:"id"`

	// Err holds the error message, if any; otherwise, it's an empty string.
	Err string `json:"err"`

	// Success indicates whether the execution was successful.
	Success bool `json:"success"`

	// Msg holds data transmission between nodes in JSON format.
	// Nodes dependent on this data need to perform JSON deserialization and agree on the expected object structure.
	Msg *gjson.JsonObject `json:"msg"`
}

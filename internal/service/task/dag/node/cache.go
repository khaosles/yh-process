package node

/*
   @File: cache.go
   @Author: khaosles
   @Time: 2023/8/20 00:15
   @Desc:
*/

// Cache is the interface for node caching.
type Cache interface {

	// Put adds an item to the cache.
	Put(string, Node)

	// Get retrieves an item from the cache.
	Get(string) (Node, bool)

	// Delete removes an item from the cache.
	Delete(string) bool

	// Clear removes all items from the cache.
	Clear()

	// Foreach applies a modification function to each value in the cache.
	Foreach(cb func(Node))

	// Len returns the length of the cache.
	Len() int
}

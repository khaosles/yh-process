package node

/*
   @File: local_cache.go
   @Author: khaosles
   @Time: 2023/8/20 00:20
   @Desc:
*/

// LocalCache represents local cache.
type LocalCache struct {
	LRUCache[string, Node]
}

// NewLocalCache creates a new LocalCache instance.
func NewLocalCache() *LocalCache {
	var localCache LocalCache
	localCache.cache = make(map[string]*lruNode[string, Node])
	localCache.head = nil
	localCache.tail = nil
	localCache.capacity = -1
	localCache.length = 0
	return &localCache
}

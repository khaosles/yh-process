package node

/*
   @File: proxy_cache.go
   @Author: khaosles
   @Time: 2023/8/20 00:40
   @Desc:
*/

// ProxyCache is a cache for storing TaskProxy instances.
type ProxyCache struct {
	LRUCache[string, *TaskProxy]
}

// NewProxyCache creates a new ProxyCache instance.
func NewProxyCache() *ProxyCache {
	var proxyCache ProxyCache
	proxyCache.cache = make(map[string]*lruNode[string, *TaskProxy])
	proxyCache.head = nil
	proxyCache.tail = nil
	proxyCache.capacity = -1
	proxyCache.length = 0
	return &proxyCache
}

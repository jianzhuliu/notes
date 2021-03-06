/*
对 lru 的一个封装，并发控制
为何不适用 RWMutex ，主要是 lru.Get 时会更新当前节点到链表表头

*/

package gocache

import (
	"gocache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
	OnRemoved  lru.CallbackFunc //键被移除时的回调函数
}

func NewCache(cacheBytes int64, fn ...lru.CallbackFunc) *cache {
	var onRemoved lru.CallbackFunc
	if len(fn) > 0 {
		onRemoved = fn[0]
	}

	return &cache{cacheBytes: cacheBytes, OnRemoved: onRemoved}
}

//添加
func (c *cache) add(key string, value ByteView) {
	//加锁
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		//延迟初始化 Lazy Initialization
		c.lru = lru.New(c.cacheBytes, c.OnRemoved)
	}

	c.lru.Add(key, value)
}

//获取
func (c *cache) get(key string) (value ByteView, ok bool) {
	//加锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//还没有添加过数据
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

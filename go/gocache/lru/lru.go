/*
缓存淘汰策略
LRU 	Least Recently Used	最近最少使用


*/
package lru

import (
	"container/list"
	"fmt"
)

type Cache struct {
	maxBytes  int64                    //允许使用的最大内存, <0 则不限制内存大小
	nbytes    int64                    //当前已经使用的内存
	ll        *list.List               //双向链表
	cache     map[string]*list.Element //键对应双向列表中节点指针
	OnRemoved CallbackFunc             //键被移除时的回调函数
}

//所有存储的值需要实现此接口，用于计算占用的内存大小
type Value interface {
	Len() int
}

//存入的键值数据
type entry struct {
	key   string //保留键，可用于删除时，通过key 快速删除 cache里面对应的记录
	value Value
}

type CallbackFunc func(string, Value)

//构造函数，初始化
func New(maxBytes int64, onRemoved CallbackFunc) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		OnRemoved: onRemoved,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
	}
}

//查找
func (c *Cache) Get(key string) (Value, bool) {
	if e, ok := c.cache[key]; ok {
		//如果缓存中存在，就移动到链表头部
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		return kv.value, true
	}

	return nil, false
}

//添加
func (c *Cache) Add(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		//如果缓存中已经存在
		kv := e.Value.(*entry)

		//更新内存大小，加新值，减原值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())

		//更新 value
		kv.value = value

		//同时移动到链表头部
		c.ll.MoveToFront(e)
	} else {
		//不存在，则新增
		kv := &entry{key: key, value: value}

		//添加到链表头部
		e := c.ll.PushFront(kv)

		//添加到缓存中
		c.cache[key] = e

		//更新内存大小
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//检测内存大小，删除链表尾部 key
	for c.maxBytes > 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//删除链表尾部 key, 移除最近最少访问的记录
func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		kv := e.Value.(*entry)

		//删除缓存中记录
		delete(c.cache, kv.key)

		//更新内存大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		//删除链表记录
		c.ll.Remove(e)

		//如果有回调函数，则调用
		if c.OnRemoved != nil {
			c.OnRemoved(kv.key, kv.value)
		}
	}
}

//获取已经存储了多少条记录
func (c *Cache) Len() int {
	return c.ll.Len()
}

//缓存内基本信息展示
func (c *Cache) Info() string {
	return fmt.Sprintf("maxBytes:%d, nbytes:%d, len:%d", c.maxBytes, c.nbytes, c.ll.Len())
}

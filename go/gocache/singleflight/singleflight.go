/*
缓存雪崩
	1、缓存在同一时候全部失效，造成瞬间DB请求增大，增加压力，造成雪崩。
	2、造成原因，服务器宕机，设置了相同的过期时间

缓存击穿
	1、某个存在的key，在过期时刻，有大量请求
	2、处理方式，查询数据库时，只有一个请求到达DB，其他请求等待，或者到期后失败

缓存穿透
	1、同一时刻，有大量请求一个不存在的 key
	2、处理方案，不存在的key，首次查询后也加一个缓存过期时间

*/
package singleflight

import (
	"sync"
)

type request struct {
	wg    sync.WaitGroup //分组，用于大量请求，等待处理结果
	value interface{}    //请求获得的结果
	err   error          //错误信息
}

type Entry struct {
	mu       sync.Mutex
	requests map[string]*request
}

//处理请求
func (e *Entry) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	e.mu.Lock()
	if e.requests == nil {
		e.requests = make(map[string]*request)
	}

	//同一时间，其他请求，处于等待中
	if r, ok := e.requests[key]; ok {
		r.wg.Wait()
		return r.value, r.err
	}

	r := new(request)

	//先记录请求
	r.wg.Add(1)
	e.requests[key] = r
	e.mu.Unlock()

	r.value, r.err = fn()
	r.wg.Done()

	//处理结束，删除记录，用于后期获取新数据
	e.mu.Lock()
	delete(e.requests, key)
	e.mu.Unlock()

	return r.value, r.err
}

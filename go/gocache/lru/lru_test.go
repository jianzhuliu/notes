/*
测试用例
*/
package lru

import (
	"reflect"
	"testing"
)

//定义一个实现了 Len(), 满足缓存中可以存储的条件。
type lruValue string

func (d lruValue) Len() int {
	return len(d)
}

//表格驱动测试
var lruTests = []struct {
	key   string
	value string
}{
	{"key1", "value1"},
	{"key2", "value2"},
	{"key3", "value3"},
}

func TestAdd(t *testing.T) {
	//创建一个不限制内存大小，且删除键没有回调函数的对象
	lru := New(int64(0), nil)

	//添加
	for _, kv := range lruTests {
		lru.Add(kv.key, lruValue(kv.value))
	}

	//t.Log(lru.Info())

	if lru.Len() != len(lruTests) {
		t.Fatalf("cache add fail")
	}

	//添加已经存在的
	lru.Add("newk1", lruValue("valuek1"))
	//t.Log(lru.Info())
	if v, ok := lru.Get("newk1"); !ok || string(v.(lruValue)) != "valuek1" {
		t.Fatalf("cache get newk1 fail")
	}

	lru.Add("newk1", lruValue("newvaluek1"))
	//t.Log(lru.Info())
	if v, ok := lru.Get("newk1"); !ok || string(v.(lruValue)) != "newvaluek1" {
		t.Fatalf("cache get newk1 fail")
	}
}

func TestGet(t *testing.T) {
	//创建一个不限制内存大小，且删除键没有回调函数的对象
	lru := New(int64(0), nil)

	//添加
	for _, kv := range lruTests {
		lru.Add(kv.key, lruValue(kv.value))
	}

	//打印缓存内基本数据
	//t.Log(lru.Info())

	//查找
	for _, kv := range lruTests {
		if v, ok := lru.Get(kv.key); !ok || string(v.(lruValue)) != kv.value {
			t.Fatalf("cache hit %s=%s fail", kv.key, kv.value)
		}
	}

	//查找不存在的
	if v, ok := lru.Get("key999"); ok {
		t.Fatalf("cache hit not cached, key999=%s", string(v.(lruValue)))
	}
}

func TestRemoveOldest(t *testing.T) {
	//构建容量，再新增
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	maxBytes := len(k1 + v1 + k2 + v2)
	lru := New(int64(maxBytes), nil)
	lru.Add(k1, lruValue(v1))
	lru.Add(k2, lruValue(v2))
	//t.Log(lru.Info())

	lru.Add(k3, lruValue(v3))
	//t.Log(lru.Info())

	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldest %s fail", k1)
	}
}

func TestOnRemoved(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		//t.Logf("callback %s=%v",key, value)
		keys = append(keys, key)
	}

	k1, k2, k3, k4 := "key1", "key2", "key3", "key4"
	v1, v2, v3, v4 := "value1", "value2", "value3", "value4"
	maxBytes := len(k1 + v1 + k2 + v2)
	lru := New(int64(maxBytes), callback)
	lru.Add(k1, lruValue(v1))
	lru.Add(k2, lruValue(v2))
	//t.Log(lru.Info())

	lru.Add(k3, lruValue(v3))
	//t.Log(lru.Info())

	lru.Add(k4, lruValue(v4))
	//t.Log(lru.Info())

	expect := []string{k1, k2}
	//t.Log(expect,keys)

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("call OnRemoved failed, expect keys %s, but got %s", expect, keys)
	}
}

package gocache

import (
	"fmt"
	"reflect"
	"testing"
)

// 8bit(位) = 1Byte(字节)
// 1024Byte = 1KB
// 1024KB = 1MB
//2^10 = 1024 		=> 1kb
//2^20 = 1048576	=> 1M

var cacheBytes int64 = 1 << 20

var db = map[string]string{
	"k1": "value1",
	"k2": "value2",
	"k3": "value3",
}

func TestGetter(t *testing.T) {
	//匿名函数 转换为 GetterFunc
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("getter failed, expect %s, but got %s", expect, v)
	}
}

func TestGet(t *testing.T) {
	loadCount := make(map[string]int)
	g := NewGroup("db", cacheBytes, GetterFunc(func(key string) ([]byte, error) {
		t.Log("[SlowDB] search key ", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			}

			loadCount[key]++
			return []byte(v), nil
		}

		return nil, fmt.Errorf("%s not exists", key)
	}))

	for k, v := range db {
		//读取缓存
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("[one] fail to get %s", k)
		}

		//再次读取，走缓存
		if _, err := g.Get(k); err != nil || loadCount[k] != 1 {
			t.Fatalf("[one] cache %s miss", k)
		}
	}

	//读取不存在的 key
	if view, err := g.Get("unknow"); err == nil {
		t.Fatalf("unknow key is not exists, but %s got", view.String())
	}
}

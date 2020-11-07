package main

import (
	"flag"
	"fmt"
	"gocache"
	"log"
	"net/http"
)

var (
	cacheBytes int64 = 1 << 20
	host       string
	port       int
)

var db = map[string]string{
	"k1": "value1",
	"k2": "value2",
	"k3": "value3",
}

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "setting the web host")
	flag.IntVar(&port, "port", 9001, "setting the web port")
}

func SetGroup() {
	gocache.NewGroup("db", cacheBytes, gocache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key ", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}

		return nil, fmt.Errorf("%s not exists", key)
	}))
}

//curl 127.0.0.1:9001/gocache/db/k1
func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", host, port)

	//设置缓存分组
	SetGroup()

	httpPool := gocache.NewHTTPPool(addr)
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, httpPool))
}

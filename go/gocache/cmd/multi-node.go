/*
win10 下执行步骤
1、编译，生成二进制 go build cmd/multi-node.go -o bin/
2、cd bin

//开启缓存服务器
3、 multi-node.exe -port 8001
4、 multi-node.exe -port 8002
5、 multi-node.exe -port 8003

//开启api对外服务器
6、 multi-node.exe -port -api

//执行测试
7、curl 127.0.0.1:9001/api?key=k4
8、curl 127.0.0.1:9001/api?key=k7
9、curl 127.0.0.1:9001/api?key=k6

*/
package main

import (
	"flag"
	"fmt"
	"gocache"
	"log"
	"net/http"
	"strconv"
)

var (
	cacheBytes int64 = 1 << 20
	dbMaxNum         = 100
	host       string
	port       int
	isApi      bool
)

var db map[string]string

//节点列表
var urlMap = map[int]string{
	8001: "http://127.0.0.1:8001",
	8002: "http://127.0.0.1:8002",
	8003: "http://127.0.0.1:8003",
}

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "setting the web host")
	flag.IntVar(&port, "port", 9001, "setting the web port")
	flag.BoolVar(&isApi, "api", false, "is api server")

	//生成测试数据
	genDB()
}

func genDB() {
	db = make(map[string]string, dbMaxNum)
	for i := 1; i <= dbMaxNum; i++ {
		istr := strconv.Itoa(i)
		db["k"+istr] = "value" + istr
	}
}

func createGroup() *gocache.Group {
	g := gocache.NewGroup("db", cacheBytes, gocache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}

		return nil, fmt.Errorf("key(%s) not exists", key)
	}))

	return g
}

func main() {
	flag.Parse()

	//创建缓存分组
	g := createGroup()

	//分布式缓存节点
	nodes := make([]string, len(urlMap))
	for _, baseUrl := range urlMap {
		nodes = append(nodes, baseUrl)
	}

	baseUrl := fmt.Sprintf("http://%s:%d", host, port)

	if !isApi {
		var ok bool
		if baseUrl, ok = urlMap[port]; !ok {
			panic(fmt.Sprintf("no node config at port=%d", port))
		}
	}

	//创建节点选择器
	nodePicker := gocache.NewHTTPPool(baseUrl)

	//设置远程节点
	nodePicker.SetNodes(nodes...)

	//注入节点选择器
	g.RegisterNodePicker(nodePicker)

	addr := baseUrl[7:]
	log.Println("server is start at ", addr)

	//启动服务器
	if isApi {
		http.HandleFunc("/api", func(rw http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(rw, "bad request", http.StatusBadRequest)
				return
			}

			view, err := g.Get(key)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}

			rw.Header().Set("Content-Type", "application/octet-stream")
			rw.WriteHeader(http.StatusOK)
			rw.Write(view.ByteSlice())
		})
		http.ListenAndServe(addr, nil)
	} else {
		//缓存服务器
		http.ListenAndServe(addr, nodePicker)
	}

}

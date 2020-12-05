package main

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"net/http"

	"gitee.com/jianzhuliu/common"
	"gitee.com/jianzhuliu/common/conf"
)

const (
	//并发数
	ParallelNum = 1000000
)

var expInt *expvar.Int

//初始化，配置核心数，以免运行过程中，系统卡死
func init() {
	numCPU := runtime.NumCPU() * 3 / 4
	if numCPU < 1 {
		numCPU = 1
	}

	runtime.GOMAXPROCS(numCPU)

	//初始化记录数为0
	expInt = expvar.NewInt("total")
	expInt.Set(0)

	conf.FlagInit(conf.Fweb)
}

//处理计数
func handleInt() {
	defer func(beginTime time.Time) {
		log.Printf("done, spend:%v, total:%v\n", time.Since(beginTime), expInt.Value())
	}(time.Now())

	var wg sync.WaitGroup
	log.Println("begin to order")
	//并发执行
	for i := 0; i < ParallelNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			expInt.Add(1)
		}()
	}

	wg.Wait()
}

func main() {
	flag.Parse()

	handleInt()

	addr := conf.FlagWebAddr()

	//浏览器自动打开
	go func() {
		common.OpenBrowser(fmt.Sprintf("%s/debug/vars", addr))
	}()

	http.ListenAndServe(addr, nil)
}

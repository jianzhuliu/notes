package main

import (
	"encoding/json"
	"flag"
	"fmt"
	bc "go_blockchain"
	"net/http"
	"strings"

	"gitee.com/jianzhuliu/tools/browser"
	"gitee.com/jianzhuliu/tools/common"
	"gitee.com/jianzhuliu/tools/conf"
	"gitee.com/jianzhuliu/tools/logger"
)

var blockChain *bc.BlockChain

func init() {
	//初始化定义相关目录，比如日志文件夹
	err := common.DefaultInit()
	if err != nil {
		common.ExitErr("common.DefaultInit()", err)
	}

	flag.Parse()

	//区块链初始化
	blockChain = bc.NewBlockChain()
}

func main() {
	addr := fmt.Sprintf("%s:%d", conf.V_host, conf.V_port)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)

	//另起 goroutine 用于浏览器自动打开
	go browser.OpenWithNotice(addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.PrintfWithTime("start serve[%s] failed,err=%v", addr, err.Error())
	}
}

//获取区块链列表
func handleIndex(rw http.ResponseWriter, r *http.Request) {
	blocks := blockChain.GetBlocks()
	bytes, err := json.Marshal(blocks)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		logger.PrintWithTime("json.Marshal() fail==", err)
		return
	}

	header := rw.Header()
	header.Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(bytes)
}

//添加区块
func handleAdd(rw http.ResponseWriter, r *http.Request) {
	data := strings.TrimSpace(r.URL.Query().Get("data"))
	if len(data) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	blockChain.AddBlock(data)
	http.Redirect(rw, r, "/", http.StatusFound)
}

package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"gitee.com/jianzhuliu/common"
	"gitee.com/jianzhuliu/common/conf"
)

func init() {
	//初始化 db 命令行参数
	conf.FlagInit(conf.Fweb)
}

func main() {
	//命令解析
	flag.Parse()

	//获取web参数配置
	addr := conf.FlagWebAddr()
	fmt.Println("addr:", addr)

	gin.DisableConsoleColor()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!!")
	})

	// curl -v 127.0.0.1:8080/ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("foo")
	})

	go common.OpenBrowser(addr)

	r.Run(addr)
}

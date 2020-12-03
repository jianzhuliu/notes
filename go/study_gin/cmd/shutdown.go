package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	log.Println("addr:", addr)

	gin.DisableConsoleColor()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Hello World!!!")
	})

	serv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	//自动打开浏览器访问
	go common.OpenBrowser(addr)

	//监听信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	log.Println("going to shutdown...")

	//优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Println("shutdown err:", err)
	}

	log.Println("server exit")

}

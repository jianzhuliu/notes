package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"

	"gitee.com/jianzhuliu/common/conf"
	"gitee.com/jianzhuliu/common/logger"
)

var activeConn sync.Map

func init() {
	//初始化命令行参数
	conf.FlagInit(conf.Fweb)

	//初始化日志
	logger.InitLogger()
}

func main() {
	//解析参数
	flag.Parse()

	//获取web命令参数对应addr
	addr := conf.FlagWebAddr()
	logger.Printf("addr: %s", addr)

	listenAndServe(addr)
}

func listenAndServe(addr string) {
	//绑定监听地址
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.ExitOnError("fail to listen,addr=%s,err:%v", addr, err)
	}

	defer l.Close()

	//开启信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())

	//另外开启 goroutine 处理监听信号
	go func() {
		sig := <-c
		switch sig {
		case os.Interrupt, os.Kill:
			logger.Printf("going to shutdown...")

			//先关闭listener阻止新连接进入
			_ = l.Close()

			cancel()

			//关闭活动连接
			activeConn.Range(func(key, value interface{}) bool {
				conn := key.(net.Conn)
				conn.Close()
				return true
			})
		}
	}()

	var wg sync.WaitGroup

LOOPFOR:
	for {
		select {
		case <-ctx.Done():
			//关闭连接
			logger.Printf("waiting disconnnect...")
			break LOOPFOR
		default:
		}

		//一直阻塞，直到有新的连接建立，或者断开连接
		conn, err := l.Accept()

		if err != nil {
			logger.Printf("fail to accept,err:%v", err)
			continue
		}

		logger.Printf("accept link")
		go func() {
			//监听是否关闭
			select {
			case <-ctx.Done():
				//关闭连接
				logger.Printf("conn disconnnect...")
				conn.Close()
				return
			default:
			}

			defer wg.Done()
			wg.Add(1)
			handleConn(ctx, conn)
		}()
	}

	wg.Wait()
}

//处理连接
func handleConn(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReader(conn)
	activeConn.Store(conn, 1)

	for {
		//以换行符为读取间隔，读取结果包含换行符
		msg, err := reader.ReadBytes('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("connection closed")
				activeConn.Delete(conn)
			} else {
				//连接被关闭，或者其他错误
				logger.Printf("fail to read,err:%v", err)
			}

			return
		}

		//原样返回
		conn.Write(msg)
	}
}

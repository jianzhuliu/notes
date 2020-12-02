package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"net"

	"gitee.com/jianzhuliu/common/conf"
	"gitee.com/jianzhuliu/common/logger"
)

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

	for {
		//一直阻塞，直到有新的连接建立，或者断开连接
		conn, err := l.Accept()

		if err != nil {
			logger.ExitOnError("fail to accept,err:%v", err)
		}

		go handleConn(conn)
	}
}

//处理连接
func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		//以换行符为读取间隔，读取结果包含换行符
		msg, err := reader.ReadBytes('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Printf("connection closed")
			} else {
				logger.Printf("fail to read,err:%v", err)
			}

			return
		}

		//原样返回
		conn.Write(msg)
	}
}

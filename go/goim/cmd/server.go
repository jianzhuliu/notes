package main

import (
	"flag"
	"goim"
)

var (
	//定义服务器创建需要的参数
	host string
	port int
)

func init() {
	//命令行配置参数
	flag.StringVar(&host, "host", "127.0.0.1", "服务器ip,默认(127.0.0.1)")
	flag.IntVar(&port, "port", 8888, "服务器端口号,默认(8888)")
}

func main() {
	//命令行解析参数
	flag.Parse()

	server := goim.NewServer(host, port)

	//启动服务器
	server.Run()
}

package main

import (
	"fmt"
	"myrpc"
	"net"
	"net/rpc/jsonrpc"
)

type Hello struct {
}

func (h *Hello) Say(name string, reply *string) error {
	*reply = "hello " + name
	return nil
}

func main() {
	//rpc注册服务
	myrpc.RegisterService(new(Hello))

	//监听端口
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("[net.Listen] err,", err)
		return
	}

	defer listener.Close()
	fmt.Println("[server] is begin to accept")

	//获取连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("[listener.Accept] err,", err)
		return
	}

	defer conn.Close()
	fmt.Println("[server] accept a conn ", conn.RemoteAddr().String())

	//rpc启动服务
	jsonrpc.ServeConn(conn)
}

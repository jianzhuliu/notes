package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

/*
满足 rpc 调用规则
1、方法可导出， 大写字母 Add
2、2个参数，都是可导出类型
3、第二个参数为指针
4、放回值为 error 类型
*/
func (a *Arith) Add(args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

var (
	Host = "127.0.0.1"
	Port = 8888
	Addr = fmt.Sprintf("%s:%d", Host, Port)

	ready chan struct{} = make(chan struct{})
)

//开启服务端
func startServer() {
	defer close(ready)

	//1、注册rpc服务
	rpc.Register(new(Arith))

	//2、开启监听
	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		fmt.Printf("[net.Listen] %s err=%v\n", Addr, err)
		return
	}

	defer listener.Close()

	////////通知可以启动客户端了
	fmt.Println("[server] begin to accept")
	ready <- struct{}{}

	//3、获取连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("[listener.Accept] err=%v\n", err)
		return
	}

	defer conn.Close()
	fmt.Println("[server] accept a conn ", conn.RemoteAddr().String())

	//4、启动rpc服务
	rpc.ServeConn(conn)
}

//启动客户端
func startClient() {
	//1、rpc拨号
	client, err := rpc.Dial("tcp", Addr)
	if err != nil {
		fmt.Printf("[rpc.Dial] %s err=%v\n", Addr, err)
		return
	}

	defer client.Close()
	fmt.Println("[client] ready to call remote service")

	//2、调用远程服务
	args := Args{A: 3, B: 5}
	reply := &Reply{}
	err = client.Call("Arith.Add", args, reply)
	if err != nil {
		fmt.Printf("[client.Call] err=%v", err)
		return
	}
	//3、输出
	fmt.Println("[client] call Arith.Add reply=", reply.C)
}

func main() {
	//启动服务器
	go startServer()

	_, ok := <-ready

	if !ok {
		fmt.Println("服务启动失败")
		return
	}

	//启动客户端
	startClient()
}

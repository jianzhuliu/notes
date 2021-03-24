package myrpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

///////////////////////////////////////////服务端 begin
//定义rpc服务需要的接口，以免注册rpc服务提供了不满足条件签名的方法
type HelloInterface interface {
	Say(string, *string) error
}

//并统一提供一个注入服务的入口
func RegisterService(target HelloInterface) {
	rpc.RegisterName("hello", target)
}

///////////////////////////////////////////服务端 end

///////////////////////////////////////////客户端 begin
type HelloClient struct {
	c *rpc.Client
}

//构建一个创建对象
func NewHelloClient(addr string) (*HelloClient, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HelloClient{c: client}, nil
}

//以jsonrpc形式创建对象
func NewHelloJsonClient(addr string) (*HelloClient, error) {
	client, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &HelloClient{c: client}, nil
}

//关闭方法
func (h *HelloClient) Close() error {
	return h.c.Close()
}

//封装调用远程方法
func (h *HelloClient) Say(name string, reply *string) error {
	return h.c.Call("hello.Say", name, reply)
}

///////////////////////////////////////////客户端 end

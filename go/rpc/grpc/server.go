package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myrpc/grpc/pb"
	"net"
)

//定义rpc服务对象
type helloServer struct {
	pb.UnimplementedHelloServer
}

//实现对象服务方法
func (h *helloServer) Say(ctx context.Context, student *pb.Student) (*pb.Student, error) {
	student.Scores = []int32{98, 84, 96}
	return student, nil
}

func main() {
	//1、注册grpc服务
	server := grpc.NewServer()
	pb.RegisterHelloServer(server, new(helloServer))

	//2、开启监听
	listener, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("[server] net.Listen err,", err)
		return
	}

	defer listener.Close()
	fmt.Println("[server] begin to accpet conn")

	//3、启动grpc服务
	server.Serve(listener)
}

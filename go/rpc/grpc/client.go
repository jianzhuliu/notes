package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myrpc/grpc/pb"
	"time"
)

func main() {
	//1、grpc 拨号连接
	conn, err := grpc.Dial("127.0.0.1:8989", grpc.WithInsecure())
	if err != nil {
		fmt.Println("[client] grpc.Dial err,", err)
		return
	}

	defer conn.Close()
	fmt.Println("[client] begin to call remote service")

	//2、创建 grpc client
	client := pb.NewHelloClient(conn)

	//3、调用远程方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Say(ctx, &pb.Student{Name: "jianzhu"})
	if err != nil {
		fmt.Println("[client] call Say err,", err)
		return
	}

	//4、打印结果
	fmt.Println("[client] name=", r.GetName(), ", scores=", r.GetScores())
}

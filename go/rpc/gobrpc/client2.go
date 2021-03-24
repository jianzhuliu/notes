package main

import (
	"fmt"
	"myrpc"
)

func main() {
	//rpc拨号
	client, err := myrpc.NewHelloClient("127.0.0.1:8000")
	if err != nil {
		fmt.Println("[client] myrpc.NewHelloClient err,", err)
		return
	}

	defer client.Close()
	fmt.Println("[client] ready to call remote service")

	//调用远程接口
	var reply string
	err = client.Say("jianzhu", &reply)
	if err != nil {
		fmt.Println("[client] Call err,", err)
		return
	}

	//打印输出
	fmt.Println("[client] reply,", reply)
}

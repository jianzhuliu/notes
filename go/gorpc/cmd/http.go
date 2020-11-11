package main

import (
	//"fmt"
	"gorpc"
	//"gorpc/codec"
	"context"
	"log"
	"net"
	"sync"
	"time"
	//"encoding/json"
	"net/http"
)

type Foo int

type Args struct {
	Num1 int
	Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

//开启服务器端
func startServer(addr chan<- string) {
	var foo Foo
	if err := gorpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	l, err := net.Listen("tcp", ":9898")
	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on", l.Addr())
	gorpc.HandleHTTP()
	addr <- l.Addr().String()
	//gorpc.Accept(l)
	http.Serve(l, nil)
}

func call(addrCh <-chan string) {
	addr := <-addrCh

	//client, err := gorpc.DialHTTP("tcp", addr)
	client, err := gorpc.XDial("http@" + addr)
	//client, err := gorpc.DialHTTP("tcp", addr, &gorpc.Option{CodecType: codec.JsonType})

	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer func() {
		_ = client.Close()
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			//args := fmt.Sprintf("gorpc req %d", i)
			args := &Args{Num1: i, Num2: i * i}
			var reply int

			if err := client.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("seq:%d, %d + %d = %d", i, args.Num1, args.Num2, reply)
		}(i)
	}

	wg.Wait()

}

func main() {
	addrCh := make(chan string)
	go call(addrCh)
	startServer(addrCh)

	for {
	}
}

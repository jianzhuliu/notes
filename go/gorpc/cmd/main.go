package main

import (
	//"fmt"
	"context"
	"gorpc"
	"gorpc/codec"
	"log"
	"net"
	"sync"
	"time"
	//"encoding/json"
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

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	gorpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	//client, _ := net.Dial("tcp", <-addr)
	//client, _ := gorpc.Dial("tcp", <-addr)
	client, _ := gorpc.Dial("tcp", <-addr, &gorpc.Option{CodecType: codec.JsonType})
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

	/*
		//_ = json.NewEncoder(client).Encode(gorpc.DefaultOption)
		//cc := codec.NewGobCodec(client)

		jsonOption := gorpc.Option{
			MagicNumber: gorpc.MagicNumber,
			CodecType: codec.JsonType,
		}
		_ = json.NewEncoder(client).Encode(jsonOption)
		cc := codec.NewJsonCodec(client)

		for i :=0; i< 5; i++ {
			h := &codec.Header{
				ServiceMethod:"Arith.Multiply",
				Seq: uint64(i),
			}

			_ = cc.Write(h, fmt.Sprintf("gorpc req %d", h.Seq))
			_ = cc.ReadHeader(h)
			var reply string
			_ = cc.ReadBody(&reply)
			log.Println("reply:", reply)
		}
		//*/
}

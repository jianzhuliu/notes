package main

import (
	"context"
	"gorpc"
	"gorpc/xclient"
	"log"
	"net"
	"sync"
	"time"
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

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

//开启服务器端
func startServer(addrCh chan<- string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}

	//需要启动多个服务，所以不适用默认的
	server := gorpc.NewServer()

	var foo Foo
	if err := server.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	log.Println("start rpc server on", l.Addr())
	addrCh <- l.Addr().String()
	server.Accept(l)
}

func main() {
	addrCh1 := make(chan string)
	addrCh2 := make(chan string)

	//开启2个服务
	go startServer(addrCh1)
	go startServer(addrCh2)

	addr1 := <-addrCh1
	addr2 := <-addrCh2

	time.Sleep(time.Second)

	rpcAddrs := []string{"tcp@" + addr1, "tcp@" + addr2}
	call(rpcAddrs...)
	log.Println("-----------------------")
	broadcast(rpcAddrs...)
}

func action(xc *xclient.XClient, ctx context.Context, actionType, serviceMethod string, args Args) {
	var reply int
	var err error

	switch actionType {
	case "call":
		err = xc.Call(ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}

	if err != nil {
		log.Printf("%s error %s(%d, %d),%v", actionType, serviceMethod, args.Num1, args.Num2, err)
	} else {
		log.Printf("%s success %s(%d, %d) = %d", actionType, serviceMethod, args.Num1, args.Num2, reply)
	}
}

//单次调用
func call(rpcAddr ...string) {
	d := xclient.NewMultiServersDiscovery(rpcAddr...)
	mode := xclient.RandomSelect

	xc := xclient.NewXClient(d, mode, nil)
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			action(xc, ctx, "call", "Foo.Sum", Args{Num1: i, Num2: i * i})
		}(i)
	}

	wg.Wait()
}

//广播
func broadcast(rpcAddr ...string) {
	d := xclient.NewMultiServersDiscovery(rpcAddr...)
	mode := xclient.RandomSelect

	xc := xclient.NewXClient(d, mode, nil)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			action(xc, context.Background(), "broadcast", "Foo.Sum", Args{Num1: i, Num2: i * i})

			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			action(xc, ctx, "broadcast", "Foo.Sleep", Args{Num1: i, Num2: i * i})

		}(i)
	}

	wg.Wait()
}

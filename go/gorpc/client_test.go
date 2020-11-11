package gorpc

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

func TestClient_dialTimeout(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("net.Listen() failed %s", err)
	}

	addr := l.Addr().String()
	f := func(conn net.Conn, opt *Option) (*Client, error) {
		_ = conn.Close()
		time.Sleep(time.Second * 2)
		return nil, nil
	}

	t.Run("timeout", func(t *testing.T) {
		_, err := dialTimeout(f, "tcp", addr, &Option{ConnectionTimeout: time.Second})
		if err == nil {
			t.Fatal("expect a timeout error")
		}

		//t.Log("err",err)
		if !strings.Contains(err.Error(), "connect timeout") {
			t.Fatal("expect a timeout error")
		}
	})

	t.Run("notimeout", func(t *testing.T) {
		_, err := dialTimeout(f, "tcp", addr, &Option{ConnectionTimeout: 0})
		if err != nil {
			t.Fatal("should not be nil, 0 means no limit")
		}
	})
}

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = Register(&b)
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	Accept(l)
}

func TestClient_Call(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)

	t.Run("client timeout", func(t *testing.T) {
		client, _ := Dial("tcp", addr)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		var reply int
		err := client.Call(ctx, "Bar.Timeout", 1, &reply)
		if err == nil {
			t.Fatal("expect a timeout error")
		}
		//t.Logf("client timeout,err:%s, ctx.Err:%s", err, ctx.Err())
		if !strings.Contains(err.Error(), ctx.Err().Error()) {
			t.Fatal("expect a timeout error")
		}
	})

	t.Run("server handle timeout", func(t *testing.T) {
		client, _ := Dial("tcp", addr, &Option{
			HandleTimeout: time.Second,
		})

		var reply int
		err := client.Call(context.Background(), "Bar.Timeout", 1, &reply)
		if err == nil {
			t.Fatal("expect a timeout error")
		}
		//t.Log("server handle timeout,err:", err)
		if !strings.Contains(err.Error(), "handle timeout") {
			t.Fatal("expect a timeout error")
		}
	})

}

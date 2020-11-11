package xclient

import (
	"context"
	"fmt"
	. "gorpc"
	"io"
	"log"
	"reflect"
	"sync"
)

type XClient struct {
	mu   sync.Mutex
	d    Discovery  //服务发现实例
	mode SelectMode //服务发现选择策略
	opt  *Option    //协议选项

	clients map[string]*Client //请求缓存，可复用
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*Client),
	}
}

//出现错误，或者异常后，需要关闭所有的连接
func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()

	for k, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, k)
	}

	return nil
}

//
func (xc *XClient) dial(rpcAddr string) (*Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()

	//首先查看缓存是否可用
	c, ok := xc.clients[rpcAddr]
	if ok {
		if c.IsAvailable() {
			log.Println("hit cache ", rpcAddr)
			return c, nil
		}
		_ = c.Close()

		//不可用，则删除
		delete(xc.clients, rpcAddr)
	}

	c, err := XDial(rpcAddr, xc.opt)
	if err != nil {
		return nil, err
	}

	log.Println("save cache ", rpcAddr)
	xc.clients[rpcAddr] = c
	return c, err
}

func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}

	err = client.Call(ctx, serviceMethod, args, reply)
	return err
}

//对外提供调用接口
func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}

	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)

}

//广播
//将请求广播到所有的服务实例，如果任意一个实例发生错误，则返回其中一个错误
//如果调用成功，则返回其中一个的结果
//为了提升性能，请求是并发的。
//并发情况下需要使用互斥锁保证 error 和 reply 能被正确赋值
//借助 context.WithCancel 确保有错误发生时，快速失败
func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	//获取所有服务列表
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}

	if len(servers) == 0 {
		return fmt.Errorf("this is no server yet")
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	var mu sync.Mutex

	//reply 为 nil, 则不需要返回值
	var replyDone = reply == nil
	var e error

	for _, rpcAddr := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var replyi interface{}
			if reply != nil {
				replyi = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}

			err = xc.call(rpcAddr, ctx, serviceMethod, args, replyi)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				//只要有个一个出错，则关闭剩余未完成的
				cancel()

				//如果历史赋值过，就不重复赋值，只取首次错误值
				if e == nil {
					e = err
				}

				return
			}

			//没有错误情况下，只取首次赋值
			if !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(replyi).Elem())
				replyDone = true
			}
		}()
	}

	wg.Wait()

	return e
}

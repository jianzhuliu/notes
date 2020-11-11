/*
客户端需要处理超时的情况
1、与服务器建立连接， net.DialTimeout
2、发送请求到服务器端，写报文导致的超时， NewClient
3、等待服务器处理
4、接收服务器响应 context.Context

*/
package gorpc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorpc/codec"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Call struct {
	Seq           uint64      //序列号
	ServiceMethod string      //请求服务及方法，同结构体对应方法，如"Arith.Multiply"
	Args          interface{} //请求参数
	Reply         interface{} //回复参数
	Error         error       //错误
	Done          chan *Call  //结束时通知
}

func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	cc       codec.Codec      //消息解码器，用来序列化将要发送出去的请求，反序列化接收到的数据
	opt      *Option          //协议交换配置
	sending  sync.Mutex       //发送消息加锁
	header   codec.Header     //头部信息
	mu       sync.Mutex       //字段信息修改加锁
	seq      uint64           //发送请求变化
	pending  map[uint64]*Call //未处理完的请求，键时编号
	closing  bool             //用户主动关闭
	shutdown bool             //错误发生时关闭

}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shut down")

//实现 io.Closer 接口
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()

	//已经关闭了,并发访问
	if client.closing {
		return ErrShutdown
	}

	client.closing = true
	return client.cc.Close()
}

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()

	return !client.shutdown && !client.closing
}

//注册一个请求
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}

	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

//移除一个请求
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

//终止一个请求，一般是出错时
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()

	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

//接收
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}

		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			//call不存在，可能是请求没有发送完整，获取其它原因被取消了
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			//服务器出错
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body:" + err.Error())
			}
			call.done()
		}
	}

	//出错情况下，终止所有未处理的请求
	client.terminateCalls(err)
}

//创建客户端
func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	//获取编码器方式
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		return nil, err
	}

	//发生协议交换配置
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error:", err)
		_ = conn.Close()
		return nil, err
	}

	return newClientCodec(f(conn), opt), nil
}

func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1,
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}

	go client.receive()
	return client
}

func parseOpions(opts ...*Option) (*Option, error) {
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption, nil
	}

	if len(opts) != 1 {
		return nil, errors.New("number of option is more than 1")
	}

	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}

	return opt, nil
}

func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewClient, network, address, opts...)
	/*
		opt, err := parseOpions(opts...)
		if err != nil {
			return nil, err
		}
		conn, err := net.Dial(network, address)
		if err != nil {
			return nil, err
		}

		defer func() {
			if client == nil {
				_ = conn.Close()
			}
		}()

		return NewClient(conn, opt)
		//*/
}

type clientResult struct {
	client *Client
	err    error
}

type clientFunc func(conn net.Conn, opt *Option) (*Client, error)

//方便测试，采用依赖注入
func dialTimeout(f clientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOpions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	ch := make(chan clientResult)

	go func() {
		client, err := f(conn, opt)
		ch <- clientResult{client: client, err: err}
	}()

	if opt.ConnectionTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}

	//超时处理
	select {
	case <-time.After(opt.ConnectionTimeout):
		return nil, fmt.Errorf("rpc client: connect timeout:expect within %s", opt.ConnectionTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}

func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()

	//注册请求
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)

		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

/////对外暴露的调用方法
//异步模式
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}

	client.send(call)
	return call

}

//同步模式,阻塞，等待结果
func (client *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	call := client.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-ctx.Done():
		client.removeCall(call.Seq)
		return fmt.Errorf("rpc client call failed: %s", ctx.Err().Error())
	case call := <-call.Done:
		return call.Error
	}
}

//////////////////////http支持
//http代理
func NewHTTPClient(conn net.Conn, opt *Option) (*Client, error) {
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", defaultRPCPath))

	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	//log.Println("client:resp.Status", resp.Status)
	if err == nil && resp.Status == connected {
		return NewClient(conn, opt)
	}

	if err == nil {
		err = fmt.Errorf("unexpected HTTP response:%s", resp.Status)
	}

	return nil, err
}

func DialHTTP(network, address string, opts ...*Option) (*Client, error) {
	return dialTimeout(NewHTTPClient, network, address, opts...)
}

//通用调用方式， rpcAddr 满足格式 protocol@addr
//比如 http@127.0.0.1:8001 , tcp@127.0.0.1:90001, unix@/tmp/gorpc.socket
func XDial(rpcAddr string, opts ...*Option) (*Client, error) {
	parts := strings.Split(rpcAddr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("rpc client err :wrong format '%s', expect protocol@addr", rpcAddr)
	}

	protocol, addr := parts[0], parts[1]
	switch protocol {
	case "http":
		return DialHTTP("tcp", addr, opts...)
	default:
		return Dial(protocol, addr, opts...)
	}

}

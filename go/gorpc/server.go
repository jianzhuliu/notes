/*
客户端与服务端的通信需要协商一些内容，例如 HTTP 报文，分为 header 和 body 2 部分，
body 的格式和长度通过 header 中的 Content-Type 和 Content-Length 指定，
服务端通过解析 header 就能够知道如何从 body 中读取需要的信息。

对于 RPC 协议来说，这部分协商是需要自主设计的。

为了提升性能，一般在报文的最开始会规划固定的字节，来协商相关的信息。
比如第1个字节用来表示序列化方式，第2个字节表示压缩方式，第3-6字节表示 header 的长度，7-10 字节表示 body 的长度

服务器处理超时的情况
1、读取客户端请求
2、发生响应报文
3、调用映射服务的方法

*/
package gorpc

import (
	"encoding/json"
	"fmt"
	"gorpc/codec"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x17bcbe9

//协议选项
type Option struct {
	MagicNumber int        //用于校验单次rpc 请求
	CodecType   codec.Type //通讯解编码类型

	ConnectionTimeout time.Duration
	HandleTimeout     time.Duration
}

//默认配置，方便客户端调用
var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         codec.GobType,
	ConnectionTimeout: time.Second * 10,
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

//接收连接请求
func (server *Server) Accept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("rpc server: accept error :", err)
			return
		}

		//log.Println("server going to serve", l.Addr().String())
		go server.ServeConn(conn)
	}
}

func Accept(l net.Listener) {
	DefaultServer.Accept(l)
}

//服务器请求
//客户端固定采用 JSON 编码 Option，后续的 header 和 body 的编码方式由 Option 中的 CodeType 指定，
//服务端首先使用 JSON 解码 Option，然后通过 Option 得 CodeType 解码剩余的内容
func (server *Server) ServeConn(conn io.ReadWriteCloser) {

	var opt Option
	//参数校验
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}

	server.serveCodec(f(conn), &opt)
}

//发生错误时，默认 body
var invalidRequest = struct{}{}

//处理编码信息
func (server *Server) serveCodec(cc codec.Codec, opt *Option) {
	//并发处理请求，回包必须逐个发送，加锁处理
	mu := new(sync.Mutex)     //不支持复制，指针传递
	wg := new(sync.WaitGroup) //确保全部请求处理完毕

	//一次连接，可以处理多个请求，即可以出现多个 header 和 body
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}

			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, mu)
			continue
		}

		wg.Add(1)
		go server.handleRequest(cc, req, mu, wg, opt.HandleTimeout)
	}

	wg.Wait()
	_ = cc.Close()
}

type request struct {
	h            *codec.Header //头部
	argv, replyv reflect.Value //请求的参数及回复数据

	mtype *methodType
	svc   *service
}

//读取头部
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}

	return &h, nil
}

//读取请求
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	//读取头部
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	//读取body
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return nil, err
	}

	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	if err := cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}

	return req, nil
}

//请求回包
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server:write response error:", err)
	}
}

//处理请求
func (server *Server) handleRequest(cc codec.Codec, req *request, mu *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	called := make(chan struct{})
	send := make(chan struct{})

	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, mu)
			send <- struct{}{}
			return
		}

		//log.Printf("handleRequest %s.%s(%v,%v) = %v, NumCalls=%d", req.svc.name, req.mtype.method.Name,req.argv,req.replyv.Type(), req.replyv.Elem().Interface(), req.mtype.NumCalls())
		server.sendResponse(cc, req.h, req.replyv.Interface(), mu)
		send <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-send
		return
	}

	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server:request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, mu)
	case <-called:
		<-send
	}
}

//注册 service
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, loaded := server.serviceMap.LoadOrStore(s.name, s); loaded {
		return fmt.Errorf("rpc:service already defined:%s", s.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

//查找 service
func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	serviceMethodArr := strings.Split(serviceMethod, ".")
	if len(serviceMethodArr) != 2 {
		err = fmt.Errorf("rpc server: service/method request ill-formed: %s", serviceMethod)
		return
	}

	serviceName, methodName := serviceMethodArr[0], serviceMethodArr[1]
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = fmt.Errorf("rpc server: can't find service %s", serviceName)
		return
	}

	svc = svci.(*service)
	mtype = svc.method[methodName]

	if mtype == nil {
		err = fmt.Errorf("rpc server: can't find method %s", methodName)
		return
	}

	return
}

/////////////支持http
const (
	connected        string = "200 Connected to go rpc"
	defaultRPCPath   string = "/gorpc"
	defaultDebugPath string = "/debug"
)

// RPC 的消息格式与标准的 HTTP 协议并不兼容，在这种情况下，就需要一个协议的转换过程。
//HTTP 协议的 CONNECT 方法恰好提供了这个能力，CONNECT 一般用于代理服务
func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		rw.Header().Set("Content-Type", "text/plain;charset=utf-8")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("405 must CONNECT"))
		return
	}

	conn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		log.Printf("rpc hijacking %s : %s", r.RemoteAddr, err.Error())
		return
	}

	log.Printf("server: RemoteAddr:%s, Host:%s, RequestURI:%s, url:%s", r.RemoteAddr, r.Host, r.RequestURI, r.URL.String())
	_, _ = io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	server.ServeConn(conn)
}

func (server *Server) HandleHTTP() {
	http.Handle(defaultRPCPath, server)
	http.Handle(defaultDebugPath, debugHTTP{server})
	log.Println("rpc server debug path:", defaultDebugPath)
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}

package goim

import (
	"fmt"
	"io"
	"net"
	"time"
)

//服务器
type Server struct {
	Host string //ip
	Port int    //端口

	//用户管理对象
	UserMgr *UserMgr

	MsgChan chan string //消息通道
}

//创建服务器对象
func NewServer(host string, port int) *Server {
	server := &Server{
		Host:    host,
		Port:    port,
		UserMgr: NewUserMgr(),
		MsgChan: make(chan string),
	}

	return server
}

//监听消息管道，如果有消息，就通知用户管理，广播消息
func (s *Server) ListenMsgChan() {
	for {
		msg := <-s.MsgChan
		s.UserMgr.Broadcast(msg)
	}
}

//通知服务器，广播消息
func (s *Server) Notify(msg string) {
	s.MsgChan <- msg
}

//处理连接
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	//根据连接创建用户
	user := NewUser(conn, s)

	//用户上线
	user.Online()

	aliveChan := make(chan struct{})

	//读取连接请求
	go func() {
		buf := make([]byte, C_READ_BUF_SIZE)

		for {
			n, err := conn.Read(buf)

			if n == 0 {
				//读取不到值了，说明已经下线了
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Printf("[%s] read msg err,%v", user, err)
				return
			}

			//提取消息，去除最后的换行符 \n
			msg := string(buf[:n-1])

			aliveChan <- struct{}{}

			user.HandleMsg(msg)
		}
	}()

	for {
		select {
		case <-aliveChan:
		case <-time.After(C_CONN_TIMEOUT):
			user.Timeout()
			return
		}
	}
}

//启动服务器
func (s *Server) Run() {
	//创建套接字监听对象
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		fmt.Printf("创建监听失败，%s:%d \n", s.Host, s.Port)
		return
	}

	fmt.Printf("[%s] 启动 \n", s)

	//开启监听消息管道
	go s.ListenMsgChan()

	//循环监听
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("监听异常, err=", err)
			continue
		}

		//处理连接
		go s.handleConn(conn)
	}
}

//打印服务器信息
func (s *Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

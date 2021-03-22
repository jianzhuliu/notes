package goim

import (
	"fmt"
	"net"
	"strings"
)

//用户
type User struct {
	Name string   //用户名
	Addr string   //登录地址
	conn net.Conn //连接对象

	server  *Server     //归属服务器对象
	MsgChan chan string //接收消息管道
}

//创建用户对象
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Addr:    addr,
		conn:    conn,
		server:  server,
		MsgChan: make(chan string),
	}

	//服务器用户管理添加用户
	server.UserMgr.AddUser(user)

	//开启监听消息
	go user.ListenMsgChan()

	return user
}

//监听消息管道，如果有消息，就发送
func (u *User) ListenMsgChan() {
	for {
		msg := <-u.MsgChan
		u.SendMsg(msg)
	}
}

//通知用户，接收消息
func (u *User) Notify(msg string) {
	u.MsgChan <- msg
}

//打印用户信息
func (u *User) String() string {
	return fmt.Sprintf("%s@%s", u.Name, u.Addr)
}

//用户上线
func (u *User) Online() {
	msg := fmt.Sprintf("[%s] 上线了", u)
	u.server.Notify(msg)
}

//用户下线
func (u *User) Offline() {
	//移除服务器用户管理器中信息
	u.server.UserMgr.RemoveUser(u)

	msg := fmt.Sprintf("[%s] 下线了", u)
	u.server.Notify(msg)

}

//用户超时退出
func (u *User) Timeout() {
	//移除服务器用户管理器中信息
	u.server.UserMgr.RemoveUser(u)

	msg := fmt.Sprintf("[%s] 超时下线了", u)
	u.server.Notify(msg)

}

//给用户当前连接发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg + "\n"))
}

//根据用户名，获取用户信息
func (u *User) GetUserByName(name string) (*User, bool) {
	return u.server.UserMgr.GetUserByName(name)
}

//用户处理消息逻辑
func (u *User) HandleMsg(msg string) {
	if msg == "who" {
		//查看所有在线用户
		users := u.server.UserMgr.GetAllUsers()
		for _, user := range users {
			u.SendMsg(fmt.Sprintf("[%s] 在线", user.Name))
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		name := strings.Split(msg, "|")[1]

		//查询用户是否已经存在
		if _, ok := u.GetUserByName(name); ok {
			u.SendMsg(fmt.Sprintf("此用户名(%s)已经存在了", name))
			return
		}

		u.server.UserMgr.RemoveUser(u)
		u.Name = name
		u.server.UserMgr.AddUser(u)
		u.SendMsg("改名成功,新名字为" + name)
	} else if len(msg) > 4 && msg[:3] == "to|" {
		msgArr := strings.Split(msg, "|")
		name := msgArr[1]
		if user, ok := u.GetUserByName(name); ok {
			user.SendMsg(fmt.Sprintf("[%s] 对你说:%s", name, msgArr[2]))
		} else {
			u.SendMsg("用户不存在，请退出再来")
		}
	} else {
		//群聊模式
		msg = fmt.Sprintf("[%s]:%s", u, msg)
		u.server.Notify(msg)
	}
}

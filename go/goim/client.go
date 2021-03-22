package goim

import (
	"fmt"
	"io"
	"net"
	"os"
)

//客户端
type Client struct {
	//服务器端配置
	ServerHost string
	ServerPort int

	conn net.Conn //连接对象
}

func NewClient(serverHost string, serverPort int) *Client {
	client := &Client{
		ServerHost: serverHost,
		ServerPort: serverPort,
	}

	return client
}

//打印配置信息
func (c *Client) String() string {
	return fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort)
}

//启动客户端
func (c *Client) Run() {
	//建立连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerHost, c.ServerPort))
	if err != nil {
		fmt.Printf("net.Dial [%s] err=%v", c, err)
		return
	}

	c.conn = conn

	//开启读取
	go c.startReader()

	var menuValue int

	for {
		c.showMenu()
		fmt.Scanln(&menuValue)

		if menuValue >= 0 && menuValue <= 3 {
			switch menuValue {
			case 1:
				//群聊模式
				c.worldChat()
			case 2:
				c.privateChat()
			case 3:
				c.rename()
			default:
				return
			}
		} else {
			fmt.Println("请输入正确类型")
		}
	}
}

//开启读数据
func (c *Client) startReader() {
	//如果有消息直接读取显示
	io.Copy(os.Stdout, c.conn)
}

//显示首层菜单
func (c *Client) showMenu() {
	fmt.Println("请选择输入类型")
	fmt.Println("1:群聊模式")
	fmt.Println("2:单聊模式")
	fmt.Println("3:改名")
	fmt.Println("0:表示退出")

}

//发送消息
func (c *Client) SendMsg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

//群聊
func (c *Client) worldChat() {
	fmt.Println("目前群聊服务中，q|Q退出")
	var msg string
	for {
		fmt.Println("请输入群聊内容，q|Q退出")
		fmt.Scanln(&msg)
		if msg == "q" || msg == "Q" {
			return
		}

		c.SendMsg(msg)
		msg = ""
	}
}

//获取在线玩家信息
func (c *Client) showUsers() {
	c.SendMsg("who")
}

//私聊
func (c *Client) privateChat() {
	fmt.Println("目前私聊服务中，q|Q退出")
	c.showUsers()
	var name string
	for {
		fmt.Println("请输入用户名")
		fmt.Scanln(&name)
		if name == "q" || name == "Q" {
			return
		}

		var msg string
		for {
			fmt.Printf("请输入对%s说的话, q|Q退出\n", name)
			fmt.Scanln(&msg)
			if msg == "q" || msg == "Q" {
				return
			}

			c.SendMsg("to|" + name + "|" + msg)
			msg = ""
		}

	}
}

//改名
func (c *Client) rename() {
	fmt.Println("请输入新的用户名")
	var msg string
	fmt.Scanln(&msg)
	c.SendMsg("rename|" + msg)
}

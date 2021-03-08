/*
错误处理
*/

package codepattern

import "fmt"

//数据库配置信息
type MysqlConf struct {
	Host    string
	Port    int
	User    string
	Passwd  string
	Dbname  string
	Charset string
}

type Option func(*MysqlConf)

func WithPort(port int) Option {
	return func(c *MysqlConf) {
		c.Port = port
	}
}

func WithUser(user string) Option {
	return func(c *MysqlConf) {
		c.User = user
	}
}

func WithPasswd(passwd string) Option {
	return func(c *MysqlConf) {
		c.Passwd = passwd
	}
}

func WithDbname(dbname string) Option {
	return func(c *MysqlConf) {
		c.Dbname = dbname
	}
}

func WithCharset(charset string) Option {
	return func(c *MysqlConf) {
		c.Charset = charset
	}
}

func NewMysqlConf(host string, options ...Option) *MysqlConf {
	s := &MysqlConf{
		Host:    host,
		Port:    3306,
		User:    "demo",
		Passwd:  "demo123",
		Dbname:  "test",
		Charset: "utf8mb4",
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *MysqlConf) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=%s",
		s.User, s.Passwd, s.Host, s.Port, s.Dbname, s.Charset,
	)
}

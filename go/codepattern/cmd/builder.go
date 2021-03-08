package main

import (
	"codepattern"
	"fmt"
)

func main() {
	sb := codepattern.MysqlConfBuilder{}
	s := sb.Create("127.0.0.1").
		WithUser("jianzhu").
		WithPasswd("123456").
		WithCharset("utf8mb4").
		WithDbname("demo")

	fmt.Println(s.Dsn())
}

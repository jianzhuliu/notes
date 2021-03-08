package main

import (
	"codepattern"
	"fmt"
)

func main() {
	s := codepattern.NewMysqlConf("127.0.0.1",
		codepattern.WithUser("jianzhu"), codepattern.WithPasswd("123456"), codepattern.WithCharset("utf8mb4"), codepattern.WithDbname("demo"),
	)

	fmt.Println(s.Dsn())
}

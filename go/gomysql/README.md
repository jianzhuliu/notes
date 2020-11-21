
# 命令行形式，读取 mysql 数据库信息

> - demo main.go

```go
package main

import (
	"fmt"
	"gomysql/command"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Println("Run()|fail|", err)
	}
}
```

> - 显示所有支持的命令列表

```bash 
go run main.go | go run main.go -h
```

> - 查看数据库版本号

```bash 
go run main.go version
```

> - 显示所有数据库名

```bash 
go run main.go databases
```

> - 显示单个数据库名下所有表

```bash 
go run main.go tables -database information_schema
```

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

> - 创建子命令，生成文件在目录 command 下

```bash 
go run main.go command -name demo -desc="demo description"
```

> - 根据数据库表，生成 struct 结构，生成目录在 models 下，文件名规则，数据库名 + 表名 + .go, 比如 gomysql_columns.go
> - 并在 cmd 目录下生成对应测试文件，格式 table_ + 表名 + .go, 比如  table_columns.go

```bash 
go run main.go command -name demo -desc="demo description"
```
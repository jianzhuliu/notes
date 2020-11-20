package gomysql

import (
	"fmt"
	"strings"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("databases", " list the databases")

	//子命令配置执行函数
	subCommand.SetRun(RunDatabases)

	//添加子命令
	AddCommand(subCommand)
}

//显示数据库列表
func RunDatabases() error {
	databases, err := GetDatabases()
	if err != nil {
		return err
	}

	if len(databases) == 0 {
		fmt.Println("the database is empty")
	} else {
		fmt.Println("=============the databases list:")
		fmt.Println(strings.Join(databases, "\r\n"))
		fmt.Println("=============total:", len(databases))
	}
	return nil
}

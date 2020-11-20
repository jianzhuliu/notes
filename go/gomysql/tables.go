package gomysql

import (
	"fmt"
	"strings"
)

var tblname string

func init() {
	//新建子命令
	subCommand := NewSubCommand("tables", " list the tables of a database")

	//子命令配置执行函数
	subCommand.SetRun(RunTables)

	//添加子命令
	AddCommand(subCommand)
}

//显示数据库列表
func RunTables() error {
	if len(V_db_name) == 0 {
		return fmt.Errorf("please set the params -database")
	}

	tables, err := GetTables()
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		fmt.Printf("the tables of %s is empty \n", V_db_name)
	} else {
		fmt.Printf("=============the tables of %s list:\n", V_db_name)
		fmt.Println(strings.Join(tables, "\r\n"))
		fmt.Println("=============total:", len(tables))
	}
	return nil
}

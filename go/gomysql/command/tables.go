package command

import (
	"fmt"
	"strings"

	"gomysql/conf"
	"gomysql/db"
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
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	if len(conf.V_db_name) == 0 {
		return fmt.Errorf("please set the params -database")
	}

	tables, err := Idb.Tables()
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		fmt.Printf("the tables of %s is empty \n", conf.V_db_name)
	} else {
		fmt.Printf("=============the tables of %s list:\n", conf.V_db_name)
		fmt.Println(strings.Join(tables, "\r\n"))
		fmt.Println("=============total:", len(tables))
	}
	return nil
}

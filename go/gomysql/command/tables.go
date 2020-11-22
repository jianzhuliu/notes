package command

import (
	"fmt"
	"strings"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("tables", " list the tables of a database")

	//子命令配置执行函数
	subCommand.SetRun(RunTables)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseTables)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseTables(sub *SubCommand) error {
	sub.SetFlagValue("check_table", "false")
	return nil
}

//显示数据库列表
func RunTables() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	tables, err := Idb.Tables()
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		fmt.Printf("the tables of %s is empty \n", conf.V_db_database)
	} else {
		fmt.Printf("=============the tables of %s list:\n", conf.V_db_database)
		fmt.Println(strings.Join(tables, "\r\n"))
		fmt.Println("=============total:", len(tables))
	}
	return nil
}

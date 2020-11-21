package command

import (
	"fmt"
	"strings"

	"gomysql/conf"
	"gomysql/db"
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
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	databases, err := Idb.Databases()
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

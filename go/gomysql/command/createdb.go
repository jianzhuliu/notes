//the sub command "createdb", created at "2020-11-30 17:38:48"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

var createDbName string

func init() {
	//新建子命令
	subCommand := NewSubCommand("createdb", "create database")

	//跳过db 校验及初始化
	//subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunCreatedb)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseCreatedb)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseCreatedb(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	sub.StringVar(&createDbName, "name", "", "the database name to be created")

	return nil
}

//查看数据库版本号
func RunCreatedb() error {
	//参数校验
	if len(createDbName) == 0 {
		return fmt.Errorf("please set the database name to be created, -name")
	}

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	db := Idb.Db()

	//创建数据库
	sql := fmt.Sprintf("create database %s default character set utf8mb4 collate utf8mb4_0900_ai_ci", createDbName)
	_, err := db.Exec(sql)

	if err != nil {
		return err
	}

	fmt.Printf("the database %s has created\n", createDbName)

	return nil
}

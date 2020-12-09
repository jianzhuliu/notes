//the sub command "droptable", created at "2020-12-09 17:12:51"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("droptable", "drop table")

	//跳过db 校验及初始化
	//subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunDroptable)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseDroptable)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseDroptable(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	/*
		//取消验证表名
		sub.SetFlagValue("check_table", "false")
		//*/

	return nil
}

//查看数据库版本号
func RunDroptable() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	//version, err := Idb.Version()

	db := Idb.Db()

	sql := fmt.Sprintf("drop table if exists %s", conf.V_db_table)
	_, err := db.Exec(sql)

	if err != nil {
		return err
	}

	fmt.Printf("%s --done\n", sql)

	return nil
}

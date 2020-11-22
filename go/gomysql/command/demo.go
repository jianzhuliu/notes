//the sub command "demo", created at "2020-11-22 17:32:47"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("demo", "demo")

	//子命令配置执行函数
	subCommand.SetRun(RunDemo)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseDemo)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseDemo(sub *SubCommand) error {
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
func RunDemo() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	//version, err := Idb.Version()

	db := Idb.Db()

	sql := fmt.Sprintf("select version()")
	row := db.QueryRow(sql)

	var version string
	if err := row.Scan(&version); err != nil {
		return err
	}

	fmt.Printf("the mysql version is %s \n", version)

	return nil
}

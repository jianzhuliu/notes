package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("version", " show the db version")

	//子命令配置执行函数
	subCommand.SetRun(RunVersion)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseVersion)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseVersion(sub *SubCommand) error {
	sub.SetFlagValue("check_database", "false")
	sub.SetFlagValue("check_table", "false")
	return nil
}

//查看数据库版本号
func RunVersion() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	version, err := Idb.Version()
	if err != nil {
		return err
	}

	fmt.Printf("the mysql version is %s \n", version)
	return nil
}

//the sub command "demo", created at "2020-11-21 15:33:46"
package command

import (
	"fmt"
	"time"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("demo", "demo desc")

	//子命令配置执行函数
	subCommand.SetRun(RunDemo)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunDemo() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	version, err := Idb.Version()
	if err != nil {
		return err
	}

	fmt.Printf("[%s] the mysql version is %s \n", time.Now().Format(conf.C_time_layout), version)
	return nil
}

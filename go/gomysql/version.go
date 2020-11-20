package gomysql

import (
	"fmt"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("version", " show the db version")

	//子命令配置执行函数
	subCommand.SetRun(RunVersion)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunVersion() error {
	version, err := GetVersion()
	if err != nil {
		return err
	}

	fmt.Printf("the mysql version is %s \n", version)
	return nil
}

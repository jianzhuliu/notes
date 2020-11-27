//the sub command "web", created at "2020-11-27 10:57:56"
package command

import (
	"fmt"

	"gomysql/web"
)

var (
	webHost string
	webPort int
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("web", "web restfull api")

	//跳过db 校验及初始化
	//subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunWeb)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseWeb)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseWeb(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	sub.StringVar(&webHost, "web_host", "127.0.0.1", "web server host")
	sub.IntVar(&webPort, "web_port", 7777, "web server port")

	return nil
}

//查看数据库版本号
func RunWeb() error {
	addr := fmt.Sprintf("%s:%d", webHost, webPort)
	web.HandleWeb(addr)

	return nil
}

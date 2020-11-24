//the sub command "start", created at "2020-11-24 10:10:14"
package command

import (
	"fmt"
	"os/exec"
	"runtime"

	"gomysql/conf"
	//"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("start", "start mysql in windows")

	//跳过db 校验及初始化
	subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunStart)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseStart)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseStart(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	return nil
}

//查看数据库版本号
func RunStart() error {
	var args []string
	var commandName string
	switch runtime.GOOS {
	case "windows":
		commandName = "cmd"
		args = []string{"/c", conf.C_win_cmd_start_mysql}
	default:
		return fmt.Errorf("not done yet on platform %s", runtime.GOOS)
	}

	cmd := exec.Command(commandName, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n%s\n", cmd.String(), stdoutStderr)

	return nil
}

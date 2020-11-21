package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"io/ioutil"
	"path/filepath"

	"gomysql/conf"
)

var (
	commandName string //子命令名称
	commandDesc string //子命令描述
	forceFlag   bool   //是否强制生成,如果已经存在
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("command", " create a sub command template")

	//跳过db 校验及初始化
	subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunCreateCommand)

	//添加命令参数
	subCommand.StringVar(&commandName, "name", "", "nameing of a sub command")
	subCommand.StringVar(&commandDesc, "desc", "", "description of a sub command")
	subCommand.BoolVar(&forceFlag, "f", false, "force to create file if exist")

	//添加子命令
	AddCommand(subCommand)
}

//创建一个子命令模板
func RunCreateCommand() error {
	//参数校验
	if len(commandName) == 0 {
		return fmt.Errorf("please set the name of a sub command, -name")
	}

	curPath, err := os.Getwd()
	if err != nil {
		return err
	}

	//不存在，则创建目录
	commandPath := filepath.Join(curPath, "command")
	_, err = os.Stat(commandPath)
	if err != nil {
		if !os.IsExist(err) {
			if err = os.Mkdir(commandPath, os.ModePerm); err != nil {
				return err
			}
		}
	}

	//生成文件
	commandName := strings.ToLower(commandName)
	commandFileName := fmt.Sprintf("%s.go", commandName)
	commandFile := filepath.Join(commandPath, commandFileName)

	//判断文件是否存在
	_, err = os.Stat(commandFile)
	if err == nil || os.IsExist(err) {
		if !forceFlag {
			return fmt.Errorf("the command file has exist, %s", commandFileName)
		}
	}

	commandNameTitle := strings.Title(commandName)
	curTime := time.Now().Format(conf.C_time_layout)
	fileData := fmt.Sprintf(commandTemplate, commandName, commandDesc, commandNameTitle, curTime)

	//写入文件
	err = ioutil.WriteFile(commandFile, []byte(fileData), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

//模板文件
var commandTemplate string = `//the sub command %[1]q, created at %[4]q
package command

import (
	"fmt"
	"time"
	
	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("%[1]s", "%[2]s")

	//子命令配置执行函数
	subCommand.SetRun(Run%[3]s)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func Run%[3]s() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%%s has not registered", conf.V_db_driver)
	}

	version, err := Idb.Version()
	if err != nil {
		return err
	}
	
	fmt.Printf("[%%s] the mysql version is %%s \n", time.Now().Format(conf.C_time_layout), version)
	return nil
}

`

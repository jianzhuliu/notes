package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"io/ioutil"
	"path/filepath"

	"gomysql/conf"
	"gomysql/utils"
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

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseCreateCommand)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseCreateCommand(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	//添加命令参数
	sub.StringVar(&commandName, "name", "", "nameing of a sub command")
	sub.StringVar(&commandDesc, "desc", "", "description of a sub command")
	sub.BoolVar(&forceFlag, "f", false, "force to create file if exist")

	return nil
}

//创建一个子命令模板
func RunCreateCommand() error {
	//参数校验

	if len(commandName) == 0 {
		return fmt.Errorf("please set the name of a sub command, -name")
	}

	//不存在，则创建目录
	commandPath, err := utils.GenFolder("command")
	if err != nil {
		return err
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

	//是否跳过检查数据库
	checkDatabase := ""
	if !conf.V_check_database {
		checkDatabase = "/"
	}

	//是否跳过检查表名
	checkTable := ""
	if !conf.V_check_table {
		checkTable = "/"
	}

	fileData := fmt.Sprintf(commandTemplate, commandName, commandDesc, commandNameTitle, curTime, checkDatabase, checkTable)

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
	
	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("%[1]s", "%[2]s")

	//跳过db 校验及初始化
	//subCommand.SetSkipDbInit(true)
	
	//子命令配置执行函数
	subCommand.SetRun(Run%[3]s)
	
	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParse%[3]s)
	
	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParse%[3]s (sub *SubCommand) error {
	%[5]s/*
	//取消验证数据库名
	sub.SetFlagValue("check_database", "false")
	//*/
	
	%[6]s/*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/
	
	return nil
}

//查看数据库版本号
func Run%[3]s() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%%s has not registered", conf.V_db_driver)
	}

	//version, err := Idb.Version()
	
	db := Idb.Db()

	sql := fmt.Sprintf("select version()")
	row := db.QueryRow(sql)
	
	var version string 
	if err := row.Scan(&version);err != nil {
		return err
	}
	
	fmt.Printf("the mysql version is %%s \n", version)
	
	
	return nil
}

`

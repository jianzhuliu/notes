//the sub command "import", created at "2020-12-09 15:34:40"
package command

import (
	"bufio"
	"fmt"
	"os"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("import", "import sql file")

	//跳过db 校验及初始化
	//subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunImport)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseImport)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseImport(sub *SubCommand) error {
	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	//sql 文件路径
	sub.StringVar(&conf.V_file, "file", "", "the sql file")
	sub.BoolVar(&conf.V_x, "x", false, "excute import sql file")

	return nil
}

//查看数据库版本号
func RunImport() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	if conf.V_file == "" {
		return fmt.Errorf("please input the sql file,by -file")
	}

	f, err := os.Open(conf.V_file)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	//以分号分割
	onSemiColon := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		for i := 0; i < len(data); i++ {
			if data[i] == ';' {
				return i + 1, data[:i], nil
			}
		}

		//结尾处
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}

	db := Idb.Db()
	scanner.Split(onSemiColon)

	for scanner.Scan() {
		sql := scanner.Text()
		if conf.V_x {
			_, err = db.Exec(sql)

			if err != nil {
				return err
			}

			fmt.Println(sql, "-----------done")
		} else {
			fmt.Println(sql)
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

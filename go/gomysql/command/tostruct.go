//the sub command "tostruct", created at "2020-11-22 17:37:21"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

var (
	all bool //是否处理所有表
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("tostruct", "format db table to struct")

	//子命令配置执行函数
	subCommand.SetRun(RunTostruct)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseTostruct)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseTostruct(sub *SubCommand) error {
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
func RunTostruct() error {
	//参数校验
	//*
	if !all && len(conf.V_db_table) == 0 {
		return fmt.Errorf("please set the table name, -table or gen all table to struct, -all")
	}
	//*/

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	if len(conf.V_db_table) > 0 {
		//处理单个表
		tableToKind, err := Idb.TableToKind(conf.V_db_database, conf.V_db_table)

		if err != nil {
			return err
		}

		fmt.Printf("%s of %s to kind:\n", conf.V_db_table, conf.V_db_database)
		for k, v := range tableToKind {
			fmt.Printf("%-20s ========> %-20s \n", k, v)
		}
		return nil
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

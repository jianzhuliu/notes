//the sub command "showcreate", created at "2020-11-21 20:07:08"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("showcreate", "show create table|database sql")

	//子命令配置执行函数
	subCommand.SetRun(RunShowcreate)

	//设置命令参数固定值
	//subCommand.SetFlagValue("check_database", "false")
	//subCommand.SetFlagValue("check_table", "false")

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunShowcreate() error {
	//参数校验
	///*
	if len(conf.V_db_database) == 0 && len(conf.V_db_table) == 0 {
		if len(conf.V_db_database) == 0 {
			return fmt.Errorf("please set the database name, -database")
		}

		if len(conf.V_db_table) == 0 {
			return fmt.Errorf("please set the table name, -table")
		}
	}
	//*/

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	//优先查看表创建 sql
	if len(conf.V_db_table) > 0 {
		create_table_sql, err := Idb.ShowCreateTableSql(conf.V_db_table)
		if err != nil {
			return err
		}

		fmt.Printf("[%s] create_table_sql:\n", conf.V_db_table)
		fmt.Println(create_table_sql)

		return nil
	}

	//查看数据库创建 sql
	if len(conf.V_db_database) > 0 {
		create_database_sql, err := Idb.ShowCreateDatabaseSql(conf.V_db_database)
		if err != nil {
			return err
		}

		fmt.Printf("[ %s] create_database_sql:\n", conf.V_db_database)
		fmt.Println(create_database_sql)

		return nil
	}

	return nil
}

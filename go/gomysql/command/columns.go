//the sub command "columns", created at "2020-11-21 17:20:29"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("columns", "show column info")

	//子命令配置执行函数
	subCommand.SetRun(RunColumns)

	//添加命令参数
	subCommand.StringVar(&conf.V_db_table, "table", "", "show fields from table")

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunColumns() error {
	//参数校验
	if len(conf.V_db_table) == 0 {
		return fmt.Errorf("please set the table name, -table")
	}

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	columnTypes, err := Idb.Columns(conf.V_db_table)
	if err != nil {
		return err
	}

	formatTmp := "|%-35s|%-15s\n"
	fmt.Printf(formatTmp, "ColumnName", "DataType")
	fmt.Println()
	for _, columnType := range columnTypes {
		fmt.Printf(formatTmp, columnType.Name(), columnType.DatabaseTypeName())
	}

	return nil
}

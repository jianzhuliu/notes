//the sub command "fields", created at "2020-11-21 15:49:45"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

var (
	tblname string
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("fields", "show fields from table")

	//子命令配置执行函数
	subCommand.SetRun(RunFields)

	//添加命令参数
	subCommand.StringVar(&tblname, "table", "", "show fields from table")

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunFields() error {
	//参数校验
	if len(tblname) == 0 {
		return fmt.Errorf("please set the table name, -table")
	}

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	fields, err := Idb.Fields(conf.V_db_name, tblname)
	if err != nil {
		return err
	}

	formatTmp := "|%-35s|%-15s|%-25s\n"
	fmt.Printf(formatTmp, "ColumnName", "DataType", "ColumnType")
	fmt.Println()
	for _, column := range fields {
		fmt.Printf(formatTmp, column.ColumnName, column.DataType, column.ColumnType)
	}

	return nil
}

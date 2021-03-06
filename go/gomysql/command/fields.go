//the sub command "fields", created at "2020-11-21 15:49:45"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("fields", "show fields from table")

	//子命令配置执行函数
	subCommand.SetRun(RunFields)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunFields() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	fields, err := Idb.Fields(conf.V_db_database, conf.V_db_table)
	if err != nil {
		return err
	}

	formatTmp := "|%-10v|%-35s|%-15s|%-10v|%-15s|%-35s|%-50s \n"
	fmt.Printf(formatTmp, "DbOrder", "ColumnName", "KindStr", "KindSize", "DataType", "ColumnType", "ColumnComment")
	fmt.Println()
	for _, column := range fields {
		fmt.Printf(formatTmp, column.DbOrder, column.ColumnName, column.KindStr, column.KindSize, column.DataType, column.ColumnType, column.ColumnComment)
	}

	return nil
}

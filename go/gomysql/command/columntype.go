//the sub command "columntype", created at "2020-11-21 19:54:44"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("columntype", "test the column type")

	//子命令配置执行函数
	subCommand.SetRun(RunColumntype)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunColumntype() error {
	//参数校验
	///*
	if len(conf.V_db_table) == 0 {
		return fmt.Errorf("please set the table name, -table")
	}
	//*/

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	db := Idb.Db()

	sql := fmt.Sprintf("select * from `%s` limit 1", conf.V_db_table)
	rows, err := db.Query(sql)

	if err != nil {
		return err
	}

	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	formatTmp := "|%-35s|%-15s|%-20s|%-15s|%-20s|%-30s\n"
	fmt.Printf(formatTmp, "ColumnName", "DataType", "ScanTypeName", "nullable, ok", "length,ok", "precision, scale,ok")
	fmt.Println()

	for _, columnType := range columnTypes {
		scanTypeName := columnType.ScanType().Name()
		nullable, ok := columnType.Nullable()
		nullStr := fmt.Sprintf("%t,%t", nullable, ok)

		length, ok := columnType.Length()
		lengthStr := fmt.Sprintf("%[1]T,%[1]v,%t", length, ok)

		precision, scale, ok := columnType.DecimalSize()
		precisionScaleStr := fmt.Sprintf("%[1]T,%[1]v,%[2]T,%[2]v,%t", precision, scale, ok)

		fmt.Printf(formatTmp, columnType.Name(), columnType.DatabaseTypeName(), scanTypeName,
			nullStr, lengthStr, precisionScaleStr)
	}

	return nil
}

//the sub command "test", created at "2020-11-21 19:25:32"
package command

import (
	"fmt"

	"gomysql/conf"
	"gomysql/db"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("test", "test")

	//子命令配置执行函数
	subCommand.SetRun(RunTest)

	//添加子命令
	AddCommand(subCommand)
}

//查看数据库版本号
func RunTest() error {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	db := Idb.Db()

	sql := fmt.Sprintf("show create table `%s`", conf.V_db_table)
	row := db.QueryRow(sql)

	var db_tblname, db_create_table string
	if err := row.Scan(&db_tblname, &db_create_table); err != nil {
		return err
	}

	_ = db_tblname
	fmt.Println("create_table_sql:")
	fmt.Println(db_create_table)

	return nil
}

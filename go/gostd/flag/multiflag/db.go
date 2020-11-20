package multiflag

import (
	"fmt"
)

//私有参数
var tblname string
var dbCommand = &SubCommand{Name: "db", Short: "数据库配置"}

func init() {
	baseCommand.Init(dbCommand)
	dbCommand.StringVar(&tblname, "table", "", "set db table name")
	dbCommand.Run = RunDB
}

func RunDB() error {
	fmt.Printf("db_host=%s, db_port=%d, db_user=%s, db_passwd=%s \n", db_host, db_port, db_user, db_passwd)
	fmt.Printf("tblname=%s \n", tblname)
	return nil
}

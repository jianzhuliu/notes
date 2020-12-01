/*
变量配置
*/
package conf

import (
	"fmt"
)

var (
	//参数配置
	V_db_host     string
	V_db_port     int
	V_db_user     string
	V_db_passwd   string
	V_db_database string
	V_db_driver   string

	V_db_table       string
	V_check_database bool
	V_check_table    bool

	V_helpFlag bool
	V_version  bool
)

//默认mysql dsn
var V_default_mysql_dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4",
	C_db_user, C_db_passwd, C_db_host, C_db_port, C_db_database)

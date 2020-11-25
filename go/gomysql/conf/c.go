/*
默认常量配置
*/
package conf

const (
	//mysql 默认数据库连接配置参数
	C_db_host     string = "127.0.0.1"
	C_db_port     int    = 3306
	C_db_user     string = "root"
	C_db_passwd   string = ""
	C_db_database string = "gomysql" // or information_schema
	C_db_driver   string = "mysql"

	//时间格式化配置
	C_time_layout string = "2006-01-02 15:04:05"

	//windows 下 开启及关闭 mysql 命令
	C_win_cmd_start_mysql = `D:\program\mysql\bin\mysqld --defaults-file=D:\program\mysql_config\my.ini --console`
	C_win_cmd_stop_mysql  = `D:\program\mysql\bin\mysqladmin -u root shutdown`

	//主键标识
	C_primary_key = "id"
)

/*
flag 常用变量
*/

package conf

import (
	"flag"
	"fmt"
)

//定义 flag 初始化类型
const (
	Fweb = 1 << iota
	Fdb
	Fstd = Fweb | Fdb
)

var (
	//web相关
	V_web_host string
	V_web_port int

	//数据库相关
	V_db_host   string
	V_db_port   int
	V_db_user   string
	V_db_passwd string
	V_db_dbname string
	V_db_driver string
)

//初始化 flag 参数
func FlagInit(flagType int) {
	if flagType&Fweb != 0 {
		flag.StringVar(&V_web_host, "host", "127.0.0.1", "setting the web host")
		flag.IntVar(&V_web_port, "port", 8080, "setting the wet port")
	}

	if flagType&Fdb != 0 {
		flag.StringVar(&V_db_host, "h", "127.0.0.1", "setting the db host")
		flag.IntVar(&V_db_port, "P", 3306, "setting the db port")
		flag.StringVar(&V_db_user, "u", "root", "setting the db user")
		flag.StringVar(&V_db_passwd, "p", "", "setting the db password")
		flag.StringVar(&V_db_dbname, "d", "", "setting the db database name")
		flag.StringVar(&V_db_driver, "D", "mysql", "setting the db driver")
	}
}

//获取 db dsn (data source name)
func FlagDbDsn() string {
	//[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=utf8",
		V_db_user, V_db_passwd, V_db_host, V_db_port, V_db_dbname,
	)

	return dsn
}

//获取 web addr
func FlagWebAddr() string {
	addr := fmt.Sprintf("%s:%d", V_web_host, V_web_port)
	return addr
}

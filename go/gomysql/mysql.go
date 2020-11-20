package gomysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Db *sql.DB

//设置db
func setDb() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8",
		V_db_user, V_db_passwd, V_db_host, V_db_port, V_db_name)
	log.Println(dns)
	Db, err = sql.Open("mysql", dns)
	if err != nil {
		return
	}

	//ping
	err = Db.Ping()

	return
}

//获取数据库版本号
func GetVersion() (version string, err error) {
	Db.QueryRow("select version()").Scan(&version)
	return
}

//获取数据库列表
func GetDatabases() ([]string, error) {
	rows, err := Db.Query("show databases")
	if err != nil {
		return nil, err
	}

	databases := []string{}

	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}

		databases = append(databases, database)
	}

	return databases, nil
}

//获取表列表
func GetTables() ([]string, error) {
	rows, err := Db.Query("show tables")
	if err != nil {
		return nil, err
	}

	tables := []string{}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

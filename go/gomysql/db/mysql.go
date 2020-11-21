package db

import (
	_ "github.com/go-sql-driver/mysql"
)

//添加直接获取mysql对应 Idb
func GetMysqlDb() Idb {
	Idb, ok := GetDb(DriverMysql)
	if !ok {
		return nil
	}

	return Idb
}

//mysql 对象
type DbMysql struct {
	DbBase
}

//确保实现了 Idb 所有接口
var _ Idb = (*DbMysql)(nil)

//初始化，注册 mysql
func init() {
	Register(DriverMysql, &DbMysql{})
}

//打开链接
func (mysql *DbMysql) Open(dsn string) error {
	return mysql.open(DriverMysql, dsn)
}

//获取版本号
func (mysql *DbMysql) Version() (string, error) {
	var version string
	err := mysql.Db().QueryRow("select version()").Scan(&version)
	if err != nil {
		return "", err
	}

	return version, err
}

//获取数据库列表
func (mysql *DbMysql) Databases() ([]string, error) {
	rows, err := mysql.Db().Query("show databases")
	if err != nil {
		return nil, err
	}

	var result []string
	var database string
	for rows.Next() {
		if err = rows.Scan(&database); err != nil {
			return nil, err
		}

		result = append(result, database)
	}
	return result, nil
}

//获取所有表
func (mysql *DbMysql) Tables() ([]string, error) {
	rows, err := mysql.Db().Query("show tables")
	if err != nil {
		return nil, err
	}

	var result []string
	var tblname string
	for rows.Next() {
		if err = rows.Scan(&tblname); err != nil {
			return nil, err
		}

		result = append(result, tblname)
	}
	return result, nil
}

package db

import (
	"database/sql"
	"fmt"
	"strings"

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

//关闭链接
func (mysql *DbMysql) Close() error {
	return mysql.Db().Close()
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

	defer rows.Close()

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

	defer rows.Close()

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

//查看某个表的字段列表
func (mysql *DbMysql) Fields(database, tblname string) ([]TableColumn, error) {
	sql := fmt.Sprintf("select column_name,column_type,data_type from information_schema.columns where table_schema=? and table_name=?")
	rows, err := mysql.Db().Query(sql, database, tblname)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []TableColumn
	var tmpData = TableColumn{}
	var orderby uint8

	for rows.Next() {
		var columnName, columnType, dateType string
		if err = rows.Scan(&columnName, &columnType, &dateType); err != nil {
			return nil, err
		}

		orderby++
		tmpData.DbOrder = orderby
		tmpData.ColumnName = columnName
		tmpData.ColumnType = strings.ToLower(columnType)
		tmpData.DataType = strings.ToLower(dateType)
	
		//匹配对应的 kind 及 size
		tmpData.KindStr, _ = mysql.KindOfDataType(tmpData)
		tmpData.KindSize,_ = KindSizeMap[tmpData.KindStr]

		result = append(result, tmpData)
	}

	return result, nil
}

//根据表名，获取表字段信息
func (mysql *DbMysql) Columns(tblname string) (columnTypes []*sql.ColumnType, err error) {
	sql := fmt.Sprintf("select * from `%s` limit 1", tblname)
	rows, err := mysql.Db().Query(sql)
	if err != nil {
		return
	}

	defer rows.Close()

	columnTypes, err = rows.ColumnTypes()
	return
}

//查看数据库创建 sql 语句
func (mysql *DbMysql) ShowCreateDatabaseSql(database string) (create_database_sql string, err error) {
	sql := fmt.Sprintf("show create database `%s`", database)
	row := mysql.Db().QueryRow(sql)
	var db_name string
	err = row.Scan(&db_name, &create_database_sql)
	_ = db_name

	return
}

//查看表创建 sql 语句
func (mysql *DbMysql) ShowCreateTableSql(tblname string) (create_table_sql string, err error) {
	sql := fmt.Sprintf("show create table `%s`", tblname)
	row := mysql.Db().QueryRow(sql)
	var db_tblname string
	err = row.Scan(&db_tblname, &create_table_sql)
	_ = db_tblname

	return
}

//获取表字段对应 go 类型
func (mysql *DbMysql) KindOfDataType(tableColumn TableColumn) (string, bool) {
	nameLen := len(tableColumn.DataType)
	if nameLen < 3 {
		return "", false
	}

	//表字段数据类型
	dataType := tableColumn.DataType

	if nameLen == 3 {
		//长度不足4个，凑够4个
		dataType += " "
	}

	var searchMap map[string]string = TypeMysqlToKind
	if strings.Contains(tableColumn.ColumnType, "unsigned") {
		searchMap = UnsignedTypeMysqlToKind
	}

	//首先根据前4个字符来查询
	namePreThree := strings.TrimSpace(dataType[0:4])
	if str, ok := searchMap[namePreThree]; ok {
		return str, true
	}

	//全字段匹配
	if str, ok := WholeTypeMysqlToKind[tableColumn.DataType]; ok {
		return str, ok
	}

	return "", false
}

/*
接口定义
*/
package db

import (
	"database/sql"
)

type Idb interface {
	//创建链接
	Open(dsn string) error

	//关闭链接
	Close() error

	//获取 db 对象
	Db() *sql.DB

	//显示数据库版本号
	Version() (string, error)

	//获取数据库列表
	Databases() ([]string, error)

	//获取所有表
	Tables() ([]string, error)

	//根据数据名及表明，获取表结构信息
	Fields(string, string) (TableColumnSice, error)

	//根据表名，获取表结构信息
	Columns(string) ([]*sql.ColumnType, error)

	//查看数据库创建 sql 语句
	ShowCreateDatabaseSql(string) (string, error)

	//查看表创建 sql 语句
	ShowCreateTableSql(string) (string, error)

	//根据字段信息，获取字段对应的 go 类型
	KindOfDataType(TableColumn) (string, bool)
}

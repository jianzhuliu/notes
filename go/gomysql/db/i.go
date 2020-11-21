/*
接口定义
*/
package db

import (
	"database/sql"
)

type Idb interface {
	Open(dsn string) error
	Db() *sql.DB
	Version() (string, error)
	Databases() ([]string, error)
	Tables() ([]string, error)
	Fields(string, string) ([]TableColumn, error)
}

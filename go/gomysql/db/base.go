package db

import (
	"database/sql"
)

//采用容器注入方式存储
var dbMap = map[string]Idb{}

//注入
func Register(driver string, db Idb) {
	dbMap[driver] = db
}

//获取
func GetDb(driver string) (db Idb, ok bool) {
	db, ok = dbMap[driver]
	return
}

//共用 db 对象
type DbBase struct {
	db *sql.DB
}

//创建数据库链接
func (base *DbBase) open(driver string, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	base.db = db

	return nil
}

//获取 db 对象
func (base *DbBase) Db() *sql.DB {
	return base.db
}

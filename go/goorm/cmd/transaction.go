package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"goorm/conf"
	"goorm/log"
)

func main() {
	db, err := sql.Open(conf.GetDsnByDriver(conf.DriverMysql))
	if db == nil || err != nil {
		log.Error("数据库连接失败", err)
		return
	}

	defer db.Close()

	_, _ = db.Exec("drop table if exists demo")
	_, err = db.Exec("create table if not exists demo(id int unsigned not null primary key,name varchar(50))")

	if err != nil {
		log.Error("建表失败")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error("开启事务失败", err)
	}

	_, err = tx.Exec("insert into demo(id,name) values(1,'name')")
	if err != nil {
		_ = tx.Rollback()
		log.Error("rollback", err)
		return
	}

	_ = tx.Commit()
	log.Info("commit")
}

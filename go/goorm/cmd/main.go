package main

import (
	_ "github.com/go-sql-driver/mysql"
	"goorm"
	"goorm/conf"
	"goorm/log"
)

func main() {
	engine := goorm.NewEngine(conf.GetDsnByDriver(conf.DriverMysql))
	if engine == nil {
		log.Error("数据库连接失败")
		return
	}

	session := engine.NewSession()
	_, err := session.Sql("drop table if exists user").Exec()
	if err != nil {
		log.Error(err)
		return
	}

	_, err = session.Sql("create table if not exists user(id int(10) unsigned not null primary key auto_increment,name varchar(30) not null default '')engine=innodb default charset=utf8mb4 ").Exec()
	if err != nil {
		log.Error(err)
		return
	}

	result, err := session.Sql("insert into user(name) values(?),(?)", "name1", "name2").Exec()
	if err != nil {
		log.Error(err)
		return
	}

	affectedNum, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("affectedNum %d", affectedNum)
}

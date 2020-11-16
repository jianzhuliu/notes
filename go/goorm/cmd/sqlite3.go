package main

import (
	"errors"
	"goorm"
	"goorm/conf"
	"goorm/log"
	"goorm/session"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `goorm:"PRIMARY KEY"`
	Age  int
}

func main() {
	engine, err := goorm.NewEngine(conf.GetDsnByDriver(conf.DriverSqlite))
	if engine == nil || err != nil {
		log.Error("数据库连接失败", err)
		return
	}

	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err = engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		log.Error("fail to rollback")
	}

	/*
		_, err = engine.Transaction(func(s *session.Session) (result interface{}, err error) {
			_ = s.Model(&User{}).CreateTable()
			_, err = s.Insert(&User{"Tom", 18})
			return
		})
		u := &User{}
		_ = s.First(u)
		if err != nil || u.Name != "Tom" {
			log.Error("failed to commit")
		}
		//*/
}

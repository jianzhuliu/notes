package goorm

import (
	"database/sql"
	"fmt"
	"goorm/dialect"
	"goorm/log"
	"goorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

//创建数据库连接
func NewEngine(driver, dsn string) (e *Engine, err error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Error(err)
		return
	}

	//确保连接可用
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	//数据库 dialect
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		msg := fmt.Sprintf("数据库引擎 %s 对应的 dialect 未找到", driver)
		log.Error(msg)
		err = fmt.Errorf(msg)
		return
	}

	log.Info("成功连接数据库")
	e = &Engine{db: db, dialect: dialect}
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
		return
	}
	log.Info("成功关闭数据库连接")
}

//创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

////事务

type TxFunc func(s *session.Session) (result interface{}, err error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}

	}()

	return f(s)
}

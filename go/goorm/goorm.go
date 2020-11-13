
package goorm

import (
	"database/sql"
	"goorm/log"
	"goorm/session"
)

type Engine struct{
	db *sql.DB
}

//创建数据库连接
func NewEngine(driver, dsn string) *Engine{
	db, err := sql.Open(driver, dsn)
	if err != nil  {
		log.Error(err)
		return nil 
	}
	
	//确保连接可用
	if err = db.Ping();err != nil {
		log.Error(err)
		return nil 
	}
	
	log.Info("成功连接数据库")
	e := &Engine{db:db}
	return e
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
		return 
	}
	log.Info("成功关闭数据库连接")
}

//创建会话
func (e *Engine) NewSession() *session.Session{
	return session.New(e.db)
}
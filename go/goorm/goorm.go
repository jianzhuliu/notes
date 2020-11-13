
package goorm

import (
	"fmt"
	"database/sql"
	"goorm/log"
	"goorm/session"
	"goorm/dialect"
)

type Engine struct{
	db *sql.DB
	dialect dialect.Dialect
}

//创建数据库连接
func NewEngine(driver, dsn string) (e *Engine, err error){
	db, err := sql.Open(driver, dsn)
	if err != nil  {
		log.Error(err)
		return 
	}
	
	//确保连接可用
	if err = db.Ping();err != nil {
		log.Error(err)
		return 
	}
	
	//数据库 dialect 
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		msg := fmt.Sprintf("数据库引擎 %s 对应的 dialect 未找到",driver)
		log.Error(msg)
		err = fmt.Errorf(msg)
		return 
	}
	
	log.Info("成功连接数据库")
	e = &Engine{db:db, dialect:dialect}
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
func (e *Engine) NewSession() *session.Session{
	return session.New(e.db, e.dialect)
}
/*
与数据库交互
*/
package session

import (
	"database/sql"
	"goorm/clause"
	"goorm/dialect"
	"goorm/log"
	"goorm/schema"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlArgs []interface{}

	//添加对结构体映射数据库后支持
	table   *schema.Schema
	dialect dialect.Dialect

	//语句构建
	clause clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

//获取db 对象
func (s *Session) DB() *sql.DB {
	return s.db
}

//清空sql语句及参数缓存
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlArgs = nil
	s.clause = clause.Clause{}
}

//构造 sql及参数
func (s *Session) Sql(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlArgs = append(s.sqlArgs, args...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlArgs)
	result, err = s.db.Exec(s.sql.String(), s.sqlArgs...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlArgs)
	rows, err = s.db.Query(s.sql.String(), s.sqlArgs...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlArgs)
	return s.db.QueryRow(s.sql.String(), s.sqlArgs...)
}

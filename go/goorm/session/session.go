/*
与数据库交互
*/
package session

import (
	"database/sql"
	"goorm/clause"
	"goorm/dialect"
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

	//事务
	tx *sql.Tx
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

type CommonDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

//获取db 对象
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}

//清空sql语句及参数缓存
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlArgs = nil
	s.clause = clause.Clause{}
}

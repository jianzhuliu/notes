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

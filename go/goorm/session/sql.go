/*
原生 sql 操作
*/
package session

import (
	"database/sql"
	"goorm/log"
)

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

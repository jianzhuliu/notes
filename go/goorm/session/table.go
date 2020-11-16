package session

import (
	"fmt"
	"goorm/log"
	"goorm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.table == nil || reflect.TypeOf(value) != reflect.TypeOf(s.table.Model) {
		s.table = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) Table() *schema.Schema {
	if s.table == nil {
		log.Error("Model is not set")
	}
	return s.table
}

//创建表
func (s *Session) CreateTable() error {
	table := s.Table()
	if table == nil {
		return fmt.Errorf("model is not set")
	}

	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	sql := fmt.Sprintf("create table %s (%s)", table.Name, strings.Join(columns, ","))
	_, err := s.Sql(sql).Exec()
	return err
}

//删除表
func (s *Session) DropTable() error {
	table := s.Table()
	if table == nil {
		return fmt.Errorf("model is not set")
	}

	sql := fmt.Sprintf("drop table %s", table.Name)
	_, err := s.Sql(sql).Exec()
	return err
}

//判断表是否存在
func (s *Session) HasTable() bool {
	table := s.Table()
	if table == nil {
		return false
	}
	tblname := table.Name
	sql, args := s.dialect.TableExistSQL(tblname)
	row := s.Sql(sql, args...).QueryRow()
	var db_tblname string
	_ = row.Scan(&db_tblname)

	return db_tblname == tblname
}

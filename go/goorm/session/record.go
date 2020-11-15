/*
crud
*/

package session

import (
	"errors"
	"goorm/clause"
	"reflect"
)

//插入，接收多个结构体对象，返回影响的函数及错误信息
// demo  s.Insert(&User{Id:1,Name:"name1"}, &User{Id:2,Name:"name2"})
func (s *Session) Insert(values ...interface{}) (int64, error) {
	insertValues := make([]interface{}, 0, len(values))
	for _, value := range values {
		table := s.Model(value).table
		s.clause.Set(clause.OpTypeInsert, table.Name, table.FieldNames)

		insertValues = append(insertValues, table.RecordValues(value))
	}

	s.clause.Set(clause.OpTypeValues, insertValues...)
	sql, sqlArgs := s.clause.Build()
	result, err := s.Sql(sql, sqlArgs...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//查找所有记录，传入空结构体对象（指针类型），填充查询后的数据
//demo  var user []User{}
//s.Find(&user)
func (s *Session) Find(values interface{}) (err error) {
	//查找切片对象对应的反射类型
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	table := s.Model(reflect.New(destType).Elem().Interface()).table

	s.clause.Set(clause.OpTypeSelect, table.Name, table.FieldNames)
	sql, sqlArgs := s.clause.Build()
	rows, err := s.Sql(sql, sqlArgs...).QueryRows()
	if err != nil {
		return
	}

	defer func() {
		err = rows.Close()
		return
	}()

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		values := make([]interface{}, 0, len(table.FieldNames))

		for _, fieldName := range table.FieldNames {
			values = append(values, dest.FieldByName(fieldName).Addr().Interface())
		}

		if err = rows.Scan(values...); err != nil {
			return
		}

		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return
}

//查找一条记录
// s.First(&User{})
func (s *Session) First(values interface{}) (err error) {
	destValue := reflect.Indirect(reflect.ValueOf(values))
	destSlice := reflect.New(reflect.SliceOf(destValue.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return errors.New("not found")
	}

	destValue.Set(destSlice.Index(0))
	return nil
}

//更新
func (s *Session) Update(values ...interface{}) (int64, error) {
	kv, ok := values[0].(map[string]interface{})
	if !ok {
		return 0, errors.New("params err, should be map[string]interface{}")
	}

	s.clause.Set(clause.OpTypeUpdate, s.table.Name, kv)
	sql, sqlArgs := s.clause.Build()
	result, err := s.Sql(sql, sqlArgs...).Exec()
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

//count
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.OpTypeCount, s.table.Name, "*")
	sql, sqlArgs := s.clause.Build()
	row := s.Sql(sql, sqlArgs...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

//delete
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.OpTypeDelete, s.table.Name)
	sql, sqlArgs := s.clause.Build()
	result, err := s.Sql(sql, sqlArgs...).Exec()
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

////////////////////链式调用
//s.Limit(10,5)
func (s *Session) Limit(args ...interface{}) *Session {
	s.clause.Set(clause.OpTypeLimit, args...)
	return s
}

//s.OrderBy("id asc", "name desc")
func (s *Session) OrderBy(args ...interface{}) *Session {
	s.clause.Set(clause.OpTypeOrderBy, args...)
	return s
}

//s.Where("id=?", 1)
func (s *Session) Where(desc string, args ...interface{}) *Session {
	values := make([]interface{}, 0, 1+len(args))
	values = append(values, desc)
	values = append(values, args...)

	s.clause.Set(clause.OpTypeWhere, values...)
	return s
}

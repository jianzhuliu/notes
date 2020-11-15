/*
crud
*/

package session

import (
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

//查找一条记录，传入空结构体对象（指针类型），填充查询后的数据
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

package session

import (
	"goorm/log"
	"reflect"
)

//定义钩子
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

//根据反射找到对应方法，并调用
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.Table().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	if fm.IsValid() {
		params := []reflect.Value{reflect.ValueOf(s)}
		if v := fm.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}

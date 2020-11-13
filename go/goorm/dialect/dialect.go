/*
不同数据库之间差异
*/
package dialect 

import (
	"reflect"
)

//数据库差异列表
var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(reflect.Value) string 	//go类型，映射为数据库类型
	TableExistSQL(string) (string,[]interface{}) //判断表是否存在的sql
}

//注册
func RegistryDialect(name string, dialect Dialect){
	dialectMap[name] = dialect
}

//获取
func GetDialect(name string) (dialect Dialect, ok bool){
	dialect, ok = dialectMap[name]
	return
}
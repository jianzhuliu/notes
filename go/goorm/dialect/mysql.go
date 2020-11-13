package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

//确保实现了 Dialect 对应接口
var _ Dialect = (*mysql)(nil)

func init() {
	RegistryDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		return "tinyint(1)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int unsigned"
	case reflect.Int64, reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := v.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", v.Type().Name(), v.Kind()))
}

func (m *mysql) TableExistSQL(tblname string) (string, []interface{}) {
	args := []interface{}{tblname}
	return "SELECT table_name FROM information_schema.tables WHERE table_schema=(select database()) and table_name = ?", args
}

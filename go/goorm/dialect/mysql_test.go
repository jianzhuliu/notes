package dialect

import (
	"reflect"
	"testing"
)

func TestMysqlDataTypeOf(t *testing.T) {
	dial := &mysql{}
	cases := []struct {
		Value interface{}
		Type  string
	}{
		{"Tom", "text"},
		{123, "int unsigned"},
		{1.2, "double"},
		{[]int{1, 2, 3}, "blob"},
	}

	for _, c := range cases {
		if typ := dial.DataTypeOf(reflect.ValueOf(c.Value)); typ != c.Type {
			t.Fatalf("expect %s, but got %s", c.Type, typ)
		}
	}
}

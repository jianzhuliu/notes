/*
go结构体与数据库表对应关系
*/

package schema

import (
	"reflect"
	"goorm/dialect"
	"go/ast"
)

//数据库字段信息
type Field struct{
	Name string 	//字段名
	Type string 	//字段对应数据库类型
	Tag string 		//标签
}

//数据库对应信息
type Schema struct{
	Model interface{} 	//处理原对象
	Name string 		//表名
	Fields []*Field 	//字段列表
	FieldNames []string //字段名列表
	fieldMap map[string]*Field //字段名对应字段信息
}

//通过字段名，获取字段信息
func (s *Schema) GetField(fieldName string) *Field{
	return s.fieldMap[fieldName]
}

type ITableName interface{
	TableName() string
}

//解析结构体为数据库对象
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var tblname string 
	t,ok := dest.(ITableName)
	//可以是映射对应表名，也可以是自定义表名
	if ok {
		tblname = t.TableName()
	} else {
		tblname = modelType.Name()
	}
	
	schema := &Schema{
		Model: dest,
		Name: tblname,
		fieldMap: make(map[string]*Field),
	}
	
	//解析表字段
	for i:=0; i< modelType.NumField();i++{
		structField := modelType.Field(i)
		fieldName := structField.Name
		//可导出且非嵌入式字段，可作为数据库字段
		if !structField.Anonymous && ast.IsExported(fieldName) {
			field := &Field{
				Name:fieldName,
				Type:d.DataTypeOf(reflect.Indirect(reflect.New(structField.Type))),
			}
			
			if v, ok := structField.Tag.Lookup("goorm");ok{
				field.Tag = v
			}
			
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, fieldName)
			schema.fieldMap[fieldName] = field
		}
	}
	
	return schema
}
package schema

import (
	"goorm/dialect"
	"testing"
)

type User struct{
	Id int `goorm:"primary key"`
	Name string 
}

//自定义表名
func (u *User) TableName() string{
	return "user"
}

var testDialect, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T){
	schema := Parse(&User{}, testDialect)
	//t.Log(schema.Name)
	//t.Log(schema.FieldNames)
	
	if schema.Name != "user" {
		t.Fatalf("db name should be user, but %s got",schema.Name)
	}
	
	if len(schema.Fields) != 2 {
		t.Fatalf("db fields num should be 2, but %d got", len(schema.Fields))
	}
	
	field := schema.GetField("Id")
	//t.Log(field)
	
	if field.Tag != "primary key" {
		t.Fatalf("parse tag fail")
	}
	
}
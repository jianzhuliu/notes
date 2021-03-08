package codepattern

import (
	"testing"
)


func TestBuilder(t *testing.T){
	sb := MysqlConfBuilder{}
	s := sb.Create("127.0.0.1").
		WithUser("jianzhu").
		WithPasswd("123456").
		WithCharset("utf8mb4").
		WithDbname("demo")

	dsn := s.Dsn()
	target := "jianzhu:123456@tcp(127.0.0.1:3306)/demo?parseTime=True&loc=Local&charset=utf8mb4"
	
	if dsn != target {
		t.Fatalf("\nshould get %q\nbut got %q", target, dsn)
	}
}
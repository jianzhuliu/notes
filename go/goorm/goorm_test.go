package goorm

import (
	_ "github.com/go-sql-driver/mysql"
	"goorm/conf"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := NewEngine(conf.GetDsnByDriver(conf.DriverMysql))

	if engine == nil || err != nil {
		t.Fatal("数据库连接失败", err)
	}

	defer engine.Close()
}

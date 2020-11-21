package db

import (
	"fmt"
	"strings"
	"testing"

	"gomysql/conf"
)

//测试使用mysql默认配置
var test_dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8",
	conf.C_db_user, conf.C_db_passwd, conf.C_db_host, conf.C_db_port, conf.C_db_name)

//获取一个db接口实现对象
func openMysql(t *testing.T) Idb {
	t.Helper()

	Idb, ok := GetDb(DriverMysql)
	if !ok {
		t.Fatal("mysql register fail")
	}

	err := Idb.Open(test_dsn)
	if err != nil {
		t.Fatal("fail to open mysql", err)
	}

	return Idb
}

//测试版本号
func TestVersion(t *testing.T) {
	Idb := openMysql(t)
	version, err := Idb.Version()
	if err != nil {
		t.Fatal("fail to get version", err)
	}

	t.Log("mysql version is ", version)
}

//测试数据库列表
func TestDatabases(t *testing.T) {
	Idb := openMysql(t)
	databases, err := Idb.Databases()
	if err != nil {
		t.Fatal("fail to get databases", err)
	}

	if len(databases) == 0 {
		t.Fatal("the databases is empty")
	}

	t.Log(strings.Join(databases, "\r\n"))
}

//测试所有表列表
func TestTables(t *testing.T) {
	Idb := openMysql(t)
	tables, err := Idb.Tables()
	if err != nil {
		t.Fatal("fail to get tables", err)
	}

	if len(tables) == 0 {
		t.Fatal("the tables is empty")
	}

	t.Log(strings.Join(tables, "\r\n"))
}

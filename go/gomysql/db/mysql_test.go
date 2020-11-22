package db

import (
	"fmt"
	"strings"
	"testing"

	"gomysql/conf"
)

//测试使用mysql默认配置
var test_dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8",
	conf.C_db_user, conf.C_db_passwd, conf.C_db_host, conf.C_db_port, conf.C_db_database)

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

type mysqlTypeTest struct {
	tableColumn TableColumn
	expect      string
	ok          bool
}

//表类型测试数据
var mysqlTypeTests = []mysqlTypeTest{
	{TableColumn{"name1", "int(11)", "int"}, "int", true},
	{TableColumn{"name2", "int(11) unsigned", "int"}, "uint", true},
	{TableColumn{"name3", "int", "int"}, "int", true},
	{TableColumn{"name4", "int unsigned", "int"}, "uint", true},

	{TableColumn{"name5", "char", "char"}, "string", true},
	{TableColumn{"name6", "char(10)", "char"}, "string", true},

	{TableColumn{"name7", "varchar(20)", "varchar"}, "string", true},
	{TableColumn{"name8", "char(10)", "char"}, "string", true},

	{TableColumn{"name9", "timestamp", "timestamp"}, "time.Time", true},

	{TableColumn{"name10", "float", "float"}, "float32", true},
	{TableColumn{"name11", "float(10,2)", "float"}, "float32", true},
}

//测试表类型映射
func TestKindOfDataType(t *testing.T) {
	Idb := openMysql(t)
	for _, info := range mysqlTypeTests {
		out, ok := Idb.KindOfDataType(info.tableColumn)
		if out != info.expect || ok != info.ok {
			t.Fatalf("%s expect %q, %t, but got %q, %t. tableColumn:%+v", info.tableColumn.ColumnName, info.expect, info.ok, out, ok, info.tableColumn)
		}
	}
}

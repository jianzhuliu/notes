package db

import (
	"strings"
	"testing"
)

//获取一个db接口实现对象
func openMysql(t *testing.T) Idb {
	t.Helper()

	Idb, err := GetMysqlIdb()
	if err != nil {
		t.Fatal(err)
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
	{TableColumn{ColumnName: "name1", ColumnType: "int(11)", DataType: "int"}, "int", true},
	{TableColumn{ColumnName: "name2", ColumnType: "int(11) unsigned", DataType: "int"}, "uint", true},
	{TableColumn{ColumnName: "name3", ColumnType: "int", DataType: "int"}, "int", true},
	{TableColumn{ColumnName: "name4", ColumnType: "int unsigned", DataType: "int"}, "uint", true},

	{TableColumn{ColumnName: "name5", ColumnType: "char", DataType: "char"}, "string", true},
	{TableColumn{ColumnName: "name6", ColumnType: "char(10)", DataType: "char"}, "string", true},

	{TableColumn{ColumnName: "name7", ColumnType: "varchar(20)", DataType: "varchar"}, "string", true},
	{TableColumn{ColumnName: "name8", ColumnType: "char(10)", DataType: "char"}, "string", true},

	{TableColumn{ColumnName: "name9", ColumnType: "timestamp", DataType: "timestamp"}, "TimeNormal", true},

	{TableColumn{ColumnName: "name10", ColumnType: "float", DataType: "float"}, "float32", true},
	{TableColumn{ColumnName: "name11", ColumnType: "float(10,2)", DataType: "float"}, "float32", true},
}

//测试表类型映射
func TestKindOfDataType(t *testing.T) {
	Idb := openMysql(t)
	for _, info := range mysqlTypeTests {
		out, ok := Idb.KindOfDataType(info.tableColumn)
		if out != info.expect || ok != info.ok {
			t.Fatalf("%s expect %q, %t, but got %q, %t ============= ColumnType:%s, DataType:%s", info.tableColumn.ColumnName, info.expect, info.ok, out, ok, info.tableColumn.ColumnType, info.tableColumn.DataType)
		}
	}
}

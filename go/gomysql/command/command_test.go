package command

import (
	"os"
	"path/filepath"
	"testing"
)

//统一测试调用入口
func runTest(args []string, t *testing.T) {
	t.Helper()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = args
	t.Log(os.Args)

	err := Run()
	if err != nil {
		t.Fatal(err)
	}
}

//测试命令版本号
func TestCommandVersion(t *testing.T) {
	args := []string{"cmd", "version"}
	runTest(args, t)
}

//测试命令数据库列表
func TestCommandDatabases(t *testing.T) {
	args := []string{"cmd", "databases"}
	runTest(args, t)
}

//测试命令展示所有表
func TestCommandTables(t *testing.T) {
	args := []string{"cmd", "tables"}
	runTest(args, t)
}

//测试命令创建新命令
func TestCommandCreate(t *testing.T) {
	args := []string{"cmd", "command", "-name", "demo", "-desc", "demo desc", "-f"}
	runTest(args, t)

	curPath, _ := os.Getwd()
	file := filepath.Join(curPath, "command", "demo.go")
	_, err := os.Stat(file)
	if err != nil {
		if !os.IsExist(err) {
			t.Fatal("fail to create new command file", err)
		}
	}
}

//测试命令显示所有字段信息
func TestCommandFields(t *testing.T) {
	args := []string{"cmd", "fields", "-table", "columns"}
	runTest(args, t)
}

//测试命令显示所有字段信息
func TestCommandColumns(t *testing.T) {
	args := []string{"cmd", "columns", "-table", "columns"}
	runTest(args, t)
}

//测试命令创建表 sql
func TestCommandCreateTableSql(t *testing.T) {
	args := []string{"cmd", "showcreate", "-database", "gomysql", "-table", "account_user"}
	runTest(args, t)
}

//测试命令创建数据库 sql
func TestCommandCreateDatabaseSql(t *testing.T) {
	args := []string{"cmd", "showcreate", "-database", "gomysql"}
	runTest(args, t)
}

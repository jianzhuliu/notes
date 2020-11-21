package command

import (
	"os"
	"testing"
)

//统一测试调用入口
func runTest(args []string, t *testing.T) {
	t.Helper()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = args

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

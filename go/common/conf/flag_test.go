package conf

import (
	"flag"
	"os"
	"testing"
)

//公用命令行参数解析
func flagInitTest(t *testing.T, flagType int, args []string) {
	t.Helper()

	//还原初始命令参数
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() { os.Args = oldArgs; flag.CommandLine = oldCommandLine }()

	os.Args = args

	//重置 flag, 以免多个测试用例共用相同参数，导致报错
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	//初始化
	FlagInit(flagType)

	//解析命令行参数
	flag.Parse()
}

//测试 db
func TestFlagDb(t *testing.T) {
	//构造参数命令， 参数解析是从 os.Args[1:]
	args := []string{"test flag db", "-h", "localhost", "-u", "jianzhu", "-P", "3306", "-d", "gocommon", "-p", "123456", "-D", "sqlite"}
	flagInitTest(t, Fdb, args)

	dsn := FlagDbDsn()
	expectDsn := "jianzhu:123456@tcp(localhost:3306)/gocommon?parseTime=True&loc=Local&charset=utf8"
	if dsn != expectDsn {
		t.Fatalf("expect %q, but %q got", expectDsn, dsn)
	}

	if V_db_driver != "sqlite" {
		t.Fatalf("db driver expect sqlite but %s got", V_db_driver)
	}
}

//测试 web
func TestFlagWeb(t *testing.T) {
	//构造参数命令， 参数解析是从 os.Args[1:]
	args := []string{"test flag web", "-host", "localhost", "-port", "9988"}
	flagInitTest(t, Fweb, args)

	addr := FlagWebAddr()
	expectAddr := "localhost:9988"
	if addr != expectAddr {
		t.Fatalf("expect %q, but %q got", expectAddr, addr)
	}
}

//测试 std
func TestFlagStd(t *testing.T) {
	//构造参数命令， 参数解析是从 os.Args[1:]
	args := []string{"test flag std", "-h", "localhost", "-u", "jianzhu", "-P", "3306", "-d", "gocommon", "-p", "123456", "-D", "sqlite", "-host", "localhost", "-port", "9988"}
	flagInitTest(t, Fstd, args)

	dsn := FlagDbDsn()
	expectDsn := "jianzhu:123456@tcp(localhost:3306)/gocommon?parseTime=True&loc=Local&charset=utf8"
	if dsn != expectDsn {
		t.Fatalf("expect %q, but %q got", expectDsn, dsn)
	}

	if V_db_driver != "sqlite" {
		t.Fatalf("db driver expect sqlite but %s got", V_db_driver)
	}

	addr := FlagWebAddr()
	expectAddr := "localhost:9988"
	if addr != expectAddr {
		t.Fatalf("expect %q, but %q got", expectAddr, addr)
	}
}

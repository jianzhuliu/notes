/*
create database goorm default character set utf8mb4 collate utf8mb4_bin
*/
package session

import (
	"database/sql"
	"fmt"
	"goorm/log"
	"os"
	"testing"
)

var (
	TestDB       *sql.DB
	TestDbDriver string = "mysql"
	TestDbHost   string = "127.0.0.1"
	TestDbPort   int    = 3306
	TestDbUser   string = "root"
	TestDbPasswd string = ""
	TestDbName   string = "goorm"
)

func exitErr(err error) {
	log.Error(err)
	os.Exit(2)
}

func setUp() {
	log.Info("setUp ================================")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		TestDbUser, TestDbPasswd, TestDbHost, TestDbPort, TestDbName)

	var err error
	TestDB, err = sql.Open(TestDbDriver, dsn)
	if err != nil {
		exitErr(err)
	}

	err = TestDB.Ping()
	if err != nil {
		exitErr(err)
	}
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func tearDown() {
	log.Info("tearDown ================================")
	TestDB.Close()
}

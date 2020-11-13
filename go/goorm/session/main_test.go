/*
create database goorm default character set utf8mb4 collate utf8mb4_bin
*/
package session

import (
	"database/sql"
	"fmt"
	"goorm/conf"
	"goorm/log"
	"os"
	"testing"
)

var (
	TestDB *sql.DB
)

func exitErr(err error) {
	log.Error(err)
	os.Exit(2)
}

func setUp() {
	log.Info("setUp ================================")

	var err error
	TestDB, err = sql.Open(conf.GetDsnByDriver(conf.DriverMysql))
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

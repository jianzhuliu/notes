package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goorm/log"
	"os"
	"strings"
	"time"
)

type opTypeFunc func() error

var (
	//数据库配置参数
	db_driver       string
	db_host         string
	db_port         int
	db_user         string
	db_passwd       string
	db_name         string
	db_ping_timeout time.Duration

	Db *sql.DB

	op_type string //操作类型

)

var ctx context.Context
var cancel context.CancelFunc

var op_type_func = map[string]opTypeFunc{
	"version":   version,
	"databases": databases,
	"tables":    tables,
}

func init() {
	//命令行参数配置
	flag.StringVar(&db_driver, "driver", "mysql", "the db driver")
	flag.StringVar(&db_host, "h", "127.0.0.1", "the db host")
	flag.IntVar(&db_port, "p", 3306, "the db port")
	flag.StringVar(&db_user, "u", "root", "the db user")
	flag.StringVar(&db_passwd, "P", "", "the db password")
	flag.StringVar(&db_name, "db_name", "", "the db name")

	flag.DurationVar(&db_ping_timeout, "timeout", 5*time.Second, "ping timeout")

	flag.StringVar(&op_type, "op", "version", "operation type")
}

//输出错误并退出
func exitErr(err error, msg string) {
	log.Error(msg, " --err:", err)
	os.Exit(2)
}

//初始化db
func setupDB() (err error) {
	//[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		db_user, db_passwd, db_host, db_port, db_name,
	)

	log.Info(dsn, db_ping_timeout)

	Db, err = sql.Open(db_driver, dsn)
	if err != nil {
		return err
	}

	//ping 一下

	err = Db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	ctx, cancel = context.WithTimeout(context.Background(), db_ping_timeout)
	defer cancel()

	err := setupDB()

	if err != nil {
		exitErr(err, "setupDB")
	}

	defer Db.Close()

	if fn, ok := op_type_func[op_type]; ok {
		err = fn()
	} else {
		err = fmt.Errorf("not defined")
	}

	if err != nil {
		exitErr(err, fmt.Sprintf("%s has err happended", op_type))
	}

	log.Info("done", op_type)
}

//显示版本号
func version() error {
	sql := "select version()"
	var version string
	row := Db.QueryRowContext(ctx, sql)
	if err := row.Scan(&version); err != nil {
		return err
	}

	log.Infof("%q == %s \n", sql, version)

	return nil
}

//显示数据库列表
func databases() error {
	sql := "show databases"

	rows, err := Db.QueryContext(ctx, sql)
	if err != nil {
		return err
	}

	defer rows.Close()

	db_map := []string{}

	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return err
		}

		db_map = append(db_map, database)
	}

	log.Info(sql, "=================", len(db_map))
	fmt.Println(strings.Join(db_map, "\n"))
	return nil
}

//显示数据库下所有表
func tables() error {
	if db_name == "" {
		return fmt.Errorf("please set the db name")
	}

	/*
		//选择数据库
		_, err := Db.ExecContext(ctx, fmt.Sprintf("use %s", db_name))
		if err != nil {
			return err
		}
		//*/

	sql := "show tables"
	rows, err := Db.QueryContext(ctx, sql)
	if err != nil {
		return err
	}

	defer rows.Close()

	tblname_map := []string{}

	for rows.Next() {
		var tblname string
		if err := rows.Scan(&tblname); err != nil {
			return err
		}

		tblname_map = append(tblname_map, tblname)
	}

	log.Info(sql, "=================", len(tblname_map))
	fmt.Println(strings.Join(tblname_map, "\n"))
	return nil
}

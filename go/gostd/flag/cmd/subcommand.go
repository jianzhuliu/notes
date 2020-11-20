package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//子命令 db
	fs_db := flag.NewFlagSet("db", flag.ExitOnError)
	
	//结果为指针类型
	var db_host *string = fs_db.String("host", "127.0.0.1", "set the db host")
	var db_port *int = fs_db.Int("port", 3306, "set the db port")
	var db_user *string = fs_db.String("user", "root", "set the db user")
	var db_passwd *string = fs_db.String("passwd", "", "set the db passwd")

	//子命令 web
	var fs_web flag.FlagSet
	fs_web.Init("web", flag.ExitOnError)
	
	var web_host string
	var web_port int
	//传递指针
	fs_web.StringVar(&web_host, "host", "localhost", "set the web host")
	fs_web.IntVar(&web_port, "port", 80, "set the web port")

	if len(os.Args) < 2 {
		fmt.Println("expect db or web subcommand")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "db":
		fs_db.Parse(os.Args[2:])
		fmt.Printf("db_host=%s, db_port=%d, db_user=%s, db_passwd=%s \n", *db_host, *db_port, *db_user, *db_passwd)
	case "web":
		fs_web.Parse(os.Args[2:])
		fmt.Printf("web_host=%s, web_port=%d \n", web_host, web_port)
	default:
		fmt.Printf("not define subcommand %s\n", os.Args[1])
	}
}

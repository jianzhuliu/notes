package command

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gomysql/conf"
	"gomysql/db"
)

//对外入口
func Run() error {
	return BaseCommand.Run()
}

//添加子命令
func AddCommand(sub *SubCommand) {
	BaseCommand.AddCommand(sub)
}

//所有命令的汇总
type Commands struct {
	fs              *flag.FlagSet
	subCommandSlice []*SubCommand
	subCommandMap   map[string]*SubCommand
}

func NewCommands() *Commands {
	c := &Commands{
		subCommandMap: make(map[string]*SubCommand),
		fs:            flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}

	//自定义说明
	c.fs.Usage = c.Usage
	c.setParams()

	return c
}

func (c *Commands) setParams() {
	//c.fs.BoolVar(&conf.V_helpFlag, "h", false, "show help information")
}

var BaseCommand = NewCommands()

//添加子命令
func (c *Commands) AddCommand(sub *SubCommand) {
	commandName := sub.name

	if sub.fs == nil {
		sub.fs = flag.NewFlagSet(commandName, flag.ExitOnError)
	}

	//设置公共参数
	setCommonParams(sub.fs)

	c.subCommandSlice = append(c.subCommandSlice, sub)
	c.subCommandMap[commandName] = sub
}

//设置公共参数
func setCommonParams(fs *flag.FlagSet) {
	fs.StringVar(&conf.V_db_host, "host", conf.C_db_host, "set the db host")
	fs.IntVar(&conf.V_db_port, "port", conf.C_db_port, "set the db port")
	fs.StringVar(&conf.V_db_user, "user", conf.C_db_user, "set the db user")
	fs.StringVar(&conf.V_db_passwd, "passwd", conf.C_db_passwd, "set the db passwd")
	fs.StringVar(&conf.V_db_database, "database", conf.C_db_database, "set the database name")
	fs.StringVar(&conf.V_db_driver, "driver", conf.C_db_driver, "set the db driver")
	fs.StringVar(&conf.V_db_table, "table", "", "set the db name")

	fs.BoolVar(&conf.V_check_database, "check_database", true, "show help information")
	fs.BoolVar(&conf.V_check_table, "check_table", true, "show help information")

	//fs.BoolVar(&conf.V_helpFlag, "h", false, "show help information")
}

func (c *Commands) Parse(args []string) error {
	return c.fs.Parse(args)
}

//执行
func (c *Commands) Run() (err error) {
	var args []string

	//没有参数
	if len(os.Args) < 2 {
		c.Usage()
		return
	}

	err = c.fs.Parse(os.Args[1:])

	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
	}

	commandName := os.Args[1]

	if commandName[0] == '-' || commandName[0] == 'h' {
		c.Usage()
		return
	}

	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	//走子命令步骤
	return c.RunCommand(commandName, args)
}

//执行子命令
func (c *Commands) RunCommand(commandName string, args []string) error {
	//匹配子命令
	if sub, ok := c.subCommandMap[commandName]; ok {
		if sub.BeforeParse != nil {
			err := sub.BeforeParse(sub)
			if err != nil {
				return err
			}
		}

		sub.parse(args)

		if conf.V_helpFlag {
			sub.fs.Usage()
			os.Exit(0)
		}

		//显示配置的参数
		sub.showParams()

		if !sub.skipDbInit {
			//db 相关校验及配置
			err := checkDb()
			if err != nil {
				return err
			}

			//关闭 db
			defer closeDb()
		}

		//执行子命令
		if sub.Run != nil {
			return sub.Run()
		} else {
			return fmt.Errorf("command [%s] has no run method", commandName)
		}
	} else {
		//没有匹配的子命令
		c.Usage()
	}

	return nil
}

//命令说明
func (c *Commands) Usage() {
	fmt.Println("Usage:")

	for _, sub := range c.subCommandSlice {
		fmt.Printf("  %s\t %s\n", sub.name, sub.usageLine)
	}

	c.fs.PrintDefaults()

	os.Exit(0)
}

//db 相关校验及配置
func checkDb() error {
	//参数校验
	if len(conf.V_db_driver) == 0 {
		return fmt.Errorf("please set the driver, -driver")
	}

	//验证数据库名
	if conf.V_check_database {
		if len(conf.V_db_database) == 0 {
			return fmt.Errorf("please set the database name, -database")
		}
	}

	//验证表
	if conf.V_check_table {
		if len(conf.V_db_table) == 0 {
			return fmt.Errorf("please set the table name, -table")
		}
	}

	//读取数据库引擎
	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	//data source name
	var dsn string
	switch conf.V_db_driver {
	case db.DriverMysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8",
			conf.V_db_user, conf.V_db_passwd, conf.V_db_host, conf.V_db_port, conf.V_db_database)
	}

	if dsn == "" {
		return fmt.Errorf("not defined driver=%s", conf.V_db_driver)
	}

	log.Println(dsn)

	//创建链接
	err := Idb.Open(dsn)
	if err != nil {
		return fmt.Errorf("fail to open db , driver=%s, dsn=%s, err=%s", conf.V_db_driver, dsn, err)
	}

	return nil
}

//关闭数据库连接
func closeDb() {
	Idb, ok := db.GetDb(conf.V_db_driver)
	if ok {
		Idb.Close()
	}
}

///////////////////////////////////////////子命令
//子命令
type SubCommand struct {
	name        string //子命令名称
	usageLine   string //一行描述
	fs          *flag.FlagSet
	Run         func() error            //执行入口,对外可配置
	BeforeParse func(*SubCommand) error //设置解析参数前处理，对外可配置

	skipDbInit bool //是否跳过 db 初始化
}

//构造子命令
func NewSubCommand(name string, usageLine string) *SubCommand {
	return &SubCommand{
		name:      name,
		usageLine: usageLine,
		fs:        flag.NewFlagSet(name, flag.ExitOnError),
	}
}

//新增参数配置
func (sub *SubCommand) SetSkipDbInit(skipDbInit bool) {
	sub.skipDbInit = skipDbInit
}

//对外可配置，说明文档
func (sub *SubCommand) SetUsage(fn func()) {
	sub.fs.Usage = fn
}

//对外可配置，设置解析参数前处理
func (sub *SubCommand) SetBeforeParse(fn func(*SubCommand) error) {
	sub.BeforeParse = fn
}

//对外可配置，执行入口
func (sub *SubCommand) SetRun(fn func() error) {
	sub.Run = fn
}

//设置参数值
func (sub *SubCommand) SetFlagValue(name, value string) {
	err := sub.fs.Set(name, value)
	if err != nil {
		panic(err)
	}
}

//解析参数
func (sub *SubCommand) parse(args []string) error {
	return sub.fs.Parse(args)
}

//显示已经配置了的参数
func (sub *SubCommand) showParams() {
	sub.fs.Visit(func(f *flag.Flag) {
		fmt.Printf("-%s=%s \n", f.Name, f.Value.String())
	})
}

//////////////////支持 Var 类型的参数配置
func (sub *SubCommand) StringVar(p *string, name string, value string, usage string) {
	sub.fs.StringVar(p, name, value, usage)
}

func (sub *SubCommand) BoolVar(p *bool, name string, value bool, usage string) {
	sub.fs.BoolVar(p, name, value, usage)
}

func (sub *SubCommand) IntVar(p *int, name string, value int, usage string) {
	sub.fs.IntVar(p, name, value, usage)
}

func (sub *SubCommand) UintVar(p *uint, name string, value uint, usage string) {
	sub.fs.UintVar(p, name, value, usage)
}

func (sub *SubCommand) Int64Var(p *int64, name string, value int64, usage string) {
	sub.fs.Int64Var(p, name, value, usage)
}

func (sub *SubCommand) Uint64Var(p *uint64, name string, value uint64, usage string) {
	sub.fs.Uint64Var(p, name, value, usage)
}

func (sub *SubCommand) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	sub.fs.DurationVar(p, name, value, usage)
}

func (sub *SubCommand) Float64Var(p *float64, name string, value float64, usage string) {
	sub.fs.Float64Var(p, name, value, usage)
}

func (sub *SubCommand) Var(value flag.Value, name string, usage string) {
	sub.fs.Var(value, name, usage)
}

package gomysql

import (
	"flag"
	"fmt"
	"os"
	"time"
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
	//c.fs.BoolVar(&V_helpFlag, "h", false, "show help information")
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
	fs.StringVar(&V_db_host, "host", "127.0.0.1", "set the db host")
	fs.IntVar(&V_db_port, "port", 3306, "set the db port")
	fs.StringVar(&V_db_user, "user", "root", "set the db user")
	fs.StringVar(&V_db_passwd, "passwd", "", "set the db passwd")
	fs.StringVar(&V_db_name, "database", "", "set the db name")

	//fs.BoolVar(&V_helpFlag, "h", false, "show help information")
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
		sub.parse(args)

		if V_helpFlag {
			sub.fs.Usage()
			os.Exit(0)
		}

		//设置 db
		err := setDb()
		if err != nil {
			return err
		}

		//显示配置的参数
		sub.showParams()

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

///////////////////////////////////////////子命令
//子命令
type SubCommand struct {
	name      string //子命令名称
	usageLine string //一行描述
	fs        *flag.FlagSet
	Run       func() error //执行入口,对外可配置
}

//构造子命令
func NewSubCommand(name string, usageLine string) *SubCommand {
	return &SubCommand{
		name:      name,
		usageLine: usageLine,
	}
}

//对外可配置，说明文档
func (sub *SubCommand) SetUsage(fn func()) {
	sub.fs.Usage = fn
}

//对外可配置，执行入口
func (sub *SubCommand) SetRun(fn func() error) {
	sub.Run = fn
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

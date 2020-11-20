package multiflag

import (
	"flag"
	"fmt"
	"os"
	"time"
)

//子命令
type SubCommand struct {
	Name    string //子命令名称
	Short   string //描述
	FlagSet *flag.FlagSet
	Run     func() error //执行入口
	Usage   func()       //使用说明
}

//解析参数
func (sub *SubCommand) parse(args []string) error {
	err := sub.FlagSet.Parse(args)
	if err != nil {
		return err
	}
	return nil
}

//////////////////支持 Var 类型的参数配置
func (sub *SubCommand) StringVar(p *string, name string, value string, usage string) {
	sub.FlagSet.StringVar(p, name, value, usage)
}

func (sub *SubCommand) BoolVar(p *bool, name string, value bool, usage string) {
	sub.FlagSet.BoolVar(p, name, value, usage)
}

func (sub *SubCommand) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	sub.FlagSet.DurationVar(p, name, value, usage)
}

func (sub *SubCommand) Float64Var(p *float64, name string, value float64, usage string) {
	sub.FlagSet.Float64Var(p, name, value, usage)
}

func (sub *SubCommand) Int64Var(p *int64, name string, value int64, usage string) {
	sub.FlagSet.Int64Var(p, name, value, usage)
}

func (sub *SubCommand) Uint64Var(p *uint64, name string, value uint64, usage string) {
	sub.FlagSet.Uint64Var(p, name, value, usage)
}

func (sub *SubCommand) UintVar(p *uint, name string, value uint, usage string) {
	sub.FlagSet.UintVar(p, name, value, usage)
}

func (sub *SubCommand) Var(value flag.Value, name string, usage string) {
	sub.FlagSet.Var(value, name, usage)
}

//所有命令的汇总
type MultiFlag struct {
	FlagSet         *flag.FlagSet //自带命令
	subCommandSlice []*SubCommand
	subCommandMap   map[string]*SubCommand
}

func NewMultiFlag() *MultiFlag {
	mf := &MultiFlag{
		subCommandMap: make(map[string]*SubCommand),
		FlagSet:       flag.NewFlagSet("", flag.ExitOnError),
	}

	//自定义说明
	mf.FlagSet.Usage = mf.Usage

	return mf
}

var baseCommand = NewMultiFlag()

func (m *MultiFlag) Init(s *SubCommand) {
	name := s.Name
	if s.FlagSet == nil {
		s.FlagSet = flag.NewFlagSet(name, flag.ExitOnError)
	}

	//设置参数
	//公共参数
	setCommonParams(s.FlagSet)

	m.subCommandSlice = append(m.subCommandSlice, s)
	m.subCommandMap[name] = s
}

func (m *MultiFlag) run(name string, args []string) error {
	//匹配子命令
	if subCommand, ok := m.subCommandMap[name]; ok {
		subCommand.parse(args)

		//子命令帮忙显示
		if helpFlag {
			if subCommand.Usage != nil {
				subCommand.Usage()
			} else {
				subCommand.FlagSet.Usage()
			}
			os.Exit(0)
		}

		//执行子命令
		if subCommand.Run != nil {
			subCommand.Run()
		}
	}

	return nil
}

//自定义说明
func (m *MultiFlag) Usage() {
	fmt.Println("Usage:")

	for _, sub := range m.subCommandSlice {
		fmt.Printf("  %s\t %s\n", sub.Name, sub.Short)
	}
}

func run(name string, args []string) error {
	return baseCommand.run(name, args)
}

//对外入口
func Run() error {
	var args []string
	var subCommandName string

	baseCommand.FlagSet.Parse(os.Args[1:])

	if len(os.Args) < 2 {
		subCommandName = "help"
	} else {
		subCommandName = os.Args[1]
		args = os.Args[2:]
	}

	err := run(subCommandName, args)
	if err != nil {
		return err
	}

	return nil
}

func Usage() {
	baseCommand.Usage()
}

var (
	db_host   string
	db_port   int
	db_user   string
	db_passwd string
	db_name   string

	helpFlag bool
)

func setCommonParams(fs *flag.FlagSet) {
	fs.StringVar(&db_host, "host", "127.0.0.1", "set the db host")
	fs.IntVar(&db_port, "port", 3306, "set the db port")
	fs.StringVar(&db_user, "user", "root", "set the db user")
	fs.StringVar(&db_passwd, "passwd", "", "set the db passwd")
	fs.StringVar(&db_name, "database", "", "set the db name")

	fs.BoolVar(&helpFlag, "help", false, "")
	fs.BoolVar(&helpFlag, "h", false, "")
}

package multiflag

import (
	"fmt"
)

//私有参数
var auto bool
var helpCommand = &SubCommand{Name: "help", Short: "帮助信息"}

func init() {
	baseCommand.Init(helpCommand)
	helpCommand.BoolVar(&auto, "auto", false, "set auto value")
	helpCommand.Run = RunHelp
}

func RunHelp() error {
	fmt.Printf("auto=%v \n", auto)
	return nil
}

//the sub command "sizeof", created at "2020-11-22 20:17:46"
package command

import (
	"fmt"
	"time"
	"unsafe"
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("sizeof", "test size of type")

	//跳过db 校验及初始化
	subCommand.SetSkipDbInit(true)

	//子命令配置执行函数
	subCommand.SetRun(RunSizeof)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseSizeof)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseSizeof(sub *SubCommand) error {
	//*
	//取消验证数据库名
	sub.SetFlagValue("check_database", "false")
	//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	return nil
}

//查看相关变量的 sizeof
func RunSizeof() error {
	var (
		v_int   int
		v_int8  int8
		v_int16 int16
		v_int32 int32

		v_uint   uint
		v_uint8  uint8
		v_uint16 uint16
		v_uint32 uint32

		v_float32 float32
		v_float64 float64

		v_time time.Time

		v_string string

		v_byte_slice  []byte
		v_uint8_slice []uint8
	)

	fmt.Printf("%-20s = %v \n", "sizeof(int)", unsafe.Sizeof(v_int))
	fmt.Printf("%-20s = %v \n", "sizeof(int8)", unsafe.Sizeof(v_int8))
	fmt.Printf("%-20s = %v \n", "sizeof(int16)", unsafe.Sizeof(v_int16))
	fmt.Printf("%-20s = %v \n", "sizeof(int32)", unsafe.Sizeof(v_int32))

	fmt.Printf("%-20s = %v \n", "sizeof(uint)", unsafe.Sizeof(v_uint))
	fmt.Printf("%-20s = %v \n", "sizeof(uint8)", unsafe.Sizeof(v_uint8))
	fmt.Printf("%-20s = %v \n", "sizeof(uint16)", unsafe.Sizeof(v_uint16))
	fmt.Printf("%-20s = %v \n", "sizeof(uint32)", unsafe.Sizeof(v_uint32))

	fmt.Printf("%-20s = %v \n", "sizeof(float32)", unsafe.Sizeof(v_float32))
	fmt.Printf("%-20s = %v \n", "sizeof(float64)", unsafe.Sizeof(v_float64))

	fmt.Printf("%-20s = %v \n", "sizeof(time.Time)", unsafe.Sizeof(v_time))

	fmt.Printf("%-20s = %v \n", "sizeof(string)", unsafe.Sizeof(v_string))

	fmt.Printf("%-20s = %v \n", "sizeof([]byte)", unsafe.Sizeof(v_byte_slice))
	fmt.Printf("%-20s = %v \n", "sizeof([]uint8)", unsafe.Sizeof(v_uint8_slice))

	return nil
}

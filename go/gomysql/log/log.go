/*
简单日志
*/
package log

import (
	"fmt"
	"log"
	"os"
)

var (
	//日志格式,包含日期,文件名及对应行数，
	logFlags    = log.LstdFlags | log.Lshortfile
	errorLogger = log.New(os.Stdout, "[error]", logFlags)
	infoLogger  = log.New(os.Stdout, "[info ]", logFlags)
)

//对外暴露接口
var (
	Error = errorLogger.Printf
	Info  = infoLogger.Printf
)

//出错时退出
func ExitOnError(format string, args ...interface{}) {
	errorLogger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(2)
}

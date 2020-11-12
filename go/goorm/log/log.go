package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	mu       sync.Mutex
	logFlags = log.LstdFlags | log.Lshortfile
	errorLog = log.New(os.Stdout, "[error]", logFlags)
	infoLog  = log.New(os.Stdout, "[info]", logFlags)
	loggers  = []*log.Logger{errorLog, infoLog}
)

//对外暴露方法
var (
	Errorf = errorLog.Printf
	Error  = errorLog.Println
	Infof  = infoLog.Printf
	Info   = infoLog.Println
)

//定义日志等级
type LogLevel int

const (
	InfoLevel LogLevel = iota
	ErrorLevel
	Disabled
)

//设置日志等级
func SetLogLevel(logLevel LogLevel) {
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if logLevel > ErrorLevel {
		errorLog.SetOutput(ioutil.Discard)
	}

	if logLevel > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}

}

package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

//获取当前命令执行目录
func GetBinPath() (curPath string, err error) {
	curPath, err = os.Getwd()
	return
}

//生成目录
func GenFolder(folder string) (string, error) {
	curPath, err := GetBinPath()
	if err != nil {
		return "", err
	}

	//不存在，则创建目录
	folderPath := filepath.Join(curPath, folder)
	_, err = os.Stat(folderPath)
	if err != nil {
		if !os.IsExist(err) {
			if err = os.Mkdir(folderPath, os.ModePerm); err != nil {
				return "", err
			}
		}
	}

	return folderPath, nil
}

//出错时退出
func ExitOnError(format string, args ...interface{}) {
	logger := log.New(os.Stderr, "[exit_on_error]", log.LstdFlags|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(2)
}

//日志记录
func DoLog(format string, args ...interface{}) {
	logger := log.New(os.Stdout, "[do_log]", log.LstdFlags|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, args...))
}

//浏览器打开url
func OpenBrowser(addr string) {
	//构造 url
	targetUrl := fmt.Sprintf("http://%s", addr)

	//不同平台区分 执行命令
	var cmdName string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmdName = "cmd"
		args = []string{"/c", "start", targetUrl}
	default:
		return
	}

	//构造超时命令
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = exec.CommandContext(ctx, cmdName, args...).Run()
}

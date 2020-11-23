package utils

import (
	"os"
	"path/filepath"
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

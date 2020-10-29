package main

import (
	"fmt"
	"sync"
)

//单例对象结构体
type ConfigData struct {
	Host string
	Port int
}

//单例对象
var configInstance *ConfigData
var once sync.Once

//懒惰加载，调用时候，首次生成对象
func GetConfigInstance() *ConfigData {
	if configInstance == nil {
		once.Do(func() {
			configInstance = &ConfigData{
				Host: "localhost",
				Port: 8080,
			}
		})
	}

	return configInstance
}

func main() {
	target := GetConfigInstance()
	fmt.Printf("%#v \n", target)
}

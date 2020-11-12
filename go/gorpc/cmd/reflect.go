/*
通过反射，可以很容易获取到一个结构体的所有方法
*/
package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	typ := reflect.TypeOf(&wg)
	fmt.Printf("//num of method is %d \n", typ.NumMethod())

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		methodType := method.Type
		argv := make([]string, 0, methodType.NumIn())
		returns := make([]string, 0, methodType.NumOut())

		//j 从 1 开始， 第 0个入参是 wg 自己
		for j := 1; j < methodType.NumIn(); j++ {
			argv = append(argv, methodType.In(j).Name())
		}

		for j := 0; j < methodType.NumOut(); j++ {
			returns = append(returns, methodType.Out(j).Name())
		}

		fmt.Printf("func (w *%s) %s(%s) %s \n",
			typ.Elem().Name(),
			method.Name,
			strings.Join(argv, ","),
			strings.Join(returns, ","),
		)
	}
}

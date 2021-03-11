/*
修饰器
*/

package codepattern

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type SumFunc func(int64, int64) int64

//反射获取函数名
func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TimedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {
		//延迟调用，计算时间
		defer func(begin time.Time) {
			fmt.Printf("%s spend %v\n", getFuncName(f), time.Since(begin))
		}(time.Now())

		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	if start > end {
		start, end = end, start
	}

	for i := start; i <= end; i++ {
		sum += i
	}

	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}

	return (start + end) * (end - start + 1) / 2
}

//////使用反射
func TimedFuncDecorator(decoPtr, fn interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value

	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			defer func(begin time.Time) {
				fmt.Printf("%s spend %v\n", getFuncName(fn), time.Since(begin))
			}(time.Now())

			out = targetFunc.Call(in)

			return
		})

	decoratedFunc.Set(v)
	return
}

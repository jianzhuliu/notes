package main

import (
	"codepattern"
	"fmt"
	"sync"
)

const nProcess = 5

//生成切片
func makeRange(min, max int) []int {
	out := make([]int, max-min+1)
	for i := range out {
		out[i] = min + i
	}

	return out
}

//合并多个channel
func merge(splitChan []<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(splitChan))
	for _, ch := range splitChan {
		go func(in <-chan int) {
			for n := range in {
				out <- n
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	//*
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for n := range codepattern.Sqrt(codepattern.Odd(codepattern.Echo(nums))) {
		fmt.Println(n)
	}

	fmt.Println("==========================")

	for n := range codepattern.Pipeline(nums, codepattern.Echo, codepattern.Odd, codepattern.Sqrt) {
		fmt.Println(n)
	}

	fmt.Println("==========================")
	//*/

	//分段求所有奇数之和
	nums = makeRange(-1000000, 1000001)
	in := codepattern.Echo(nums)

	var splitChan [nProcess]<-chan int

	for i := range splitChan {
		splitChan[i] = codepattern.Sum(codepattern.Odd(in))
	}

	out := <-codepattern.Sum(merge(splitChan[:]))
	fmt.Println(out)
}

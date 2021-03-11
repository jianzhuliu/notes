/*
管道编程
*/

package codepattern

type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func Pipeline(nums []int, echo EchoFunc, pipeFuncs ...PipeFunc) <-chan int {
	out := echo(nums)

	for _, fn := range pipeFuncs {
		out = fn(out)
	}

	return out
}

//将整数切片传递到 channel中
func Echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}

		close(out)
	}()

	return out
}

//平方
func Sqrt(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * n
		}

		close(out)
	}()

	return out
}

//过滤奇数
func Odd(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			if n%2 != 0 {
				out <- n
			}
		}

		close(out)
	}()

	return out
}

//求和
func Sum(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		var sum int
		for n := range in {
			sum += n
		}

		out <- sum
		close(out)
	}()

	return out
}

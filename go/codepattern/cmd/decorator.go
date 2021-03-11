package main

import (
	"codepattern"
	"fmt"
)

func main() {
	sum1 := codepattern.TimedSumFunc(codepattern.Sum1)
	sum2 := codepattern.TimedSumFunc(codepattern.Sum2)

	fmt.Printf("Sum1 = %d\n", sum1(-10000, 10000))
	fmt.Printf("Sum2 = %d\n", sum2(-10000, 10000))

	var d_sum1 codepattern.SumFunc
	_ = codepattern.TimedFuncDecorator(&d_sum1, codepattern.Sum1)
	fmt.Printf("d_Sum1 = %d\n", d_sum1(-10000, 10001))

	d_sum2 := codepattern.Sum2
	_ = codepattern.TimedFuncDecorator(&d_sum2, codepattern.Sum2)
	fmt.Printf("d_Sum2 = %d\n", d_sum2(-10000, 10002))
}

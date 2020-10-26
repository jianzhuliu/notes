package main

import "fmt"

func DeferFunc(i int) (t int) {

	fmt.Println("t = ", t)

	return 2
}

func main() {
	DeferFunc(10)
}

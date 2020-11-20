package main

import (
	"fmt"
	"multiflag"
)

func main() {
	err := multiflag.Run()
	if err != nil {
		fmt.Println("[main] run(),fail,", err)
	}
}

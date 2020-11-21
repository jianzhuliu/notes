package main

import (
	"fmt"
	"gomysql/command"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Println("Run()|fail|", err)
	}
}

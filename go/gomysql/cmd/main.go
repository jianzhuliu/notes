package main

import (
	"fmt"
	"gomysql"
)

func main() {
	err := gomysql.Run()
	if err != nil {
		fmt.Println("Run()|fail|", err)
	}
}

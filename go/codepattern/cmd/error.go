package main

import (
	"codepattern"
	"fmt"
)

func main() {
	a := codepattern.Account{Name: "take", Age: 60, Gender: 1}
	a.CheckName().CheckAge().CheckGender().Print()
	if err := a.Error(); err != nil {
		fmt.Println(err)
	}
}

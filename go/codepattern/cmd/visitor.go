package main

import (
	"codepattern"
	"fmt"
)

func main() {
	info := codepattern.Info{}

	var v codepattern.Visitor = &info
	v = codepattern.LogVisitor{v}
	v = codepattern.NameVisitor{v}
	v = codepattern.OtherThingsVisitor{v}

	loadInfo := func(info *codepattern.Info, err error) error {
		info.Name = "mysql"
		info.Namespace = "db"
		info.OtherThings = "this is a mysql visitor"

		return nil
	}

	v.Visit(loadInfo)

	fmt.Println("\n============================")
	var v2 codepattern.Visitor = &info
	v2 = codepattern.NewDecoratorVisitor(v2, codepattern.NameVisitorFunc, codepattern.OtherThingsVisitorFunc)
	v2.Visit(loadInfo)
}

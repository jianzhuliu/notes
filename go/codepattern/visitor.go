/*
visitor
*/

package codepattern

import "fmt"

type VisitorFunc func(*Info, error) error

type Visitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Name        string
	Namespace   string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

/////////////////////////修饰器

type DecoratorVisitor struct {
	visitor       Visitor
	decoratorFunc []VisitorFunc
}

func NewDecoratorVisitor(visitor Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return visitor
	}
	return DecoratorVisitor{visitor, fn}
}

func (v DecoratorVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		err = fn(info, err)

		for i := range v.decoratorFunc {
			if err = v.decoratorFunc[i](info, err); err != nil {
				return err
			}
		}

		return err
	})
}

func NameVisitorFunc(info *Info, err error) error {
	fmt.Println("NameVisitorFunc before")
	if err == nil {
		fmt.Printf("=====> Name=%s, Namespace=%s\n", info.Name, info.Namespace)
	}

	fmt.Println("NameVisitorFunc after")
	return nil
}

func OtherThingsVisitorFunc(info *Info, err error) error {
	fmt.Println("OtherThingsVisitorFunc before")
	if err == nil {
		fmt.Printf("=====> OtherThings=%s\n", info.OtherThings)
	}

	fmt.Println("OtherThingsVisitorFunc after")
	return nil
}

/////////////////////////具体实现

type NameVisitor struct {
	Visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.Visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor before")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("=====> Name=%s, Namespace=%s\n", info.Name, info.Namespace)
		}

		fmt.Println("NameVisitor after")

		return err
	})
}

type OtherThingsVisitor struct {
	Visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.Visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor before")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("=====> OtherThings=%s\n", info.OtherThings)
		}

		fmt.Println("OtherThingsVisitor after")

		return err
	})
}

type LogVisitor struct {
	Visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.Visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor before")
		err = fn(info, err)
		fmt.Println("LogVisitor after")

		return err
	})
}

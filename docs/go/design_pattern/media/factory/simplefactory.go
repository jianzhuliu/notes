package main

import "fmt"

//定义一个接口
type Factory interface {
	Made()
}

//对外提供一个接口，根据参数来匹配工厂
func NewFactory(t int) Factory {
	if t == 1 {
		return NewAFacatory()
	} else {
		return NewBFacatory()
	}
}

//A 工厂
type AFacatory struct {
}

func NewAFacatory() AFacatory {
	return AFacatory{}
}

func (a AFacatory) Made() {
	fmt.Println("made in A")
}

//B 工厂
type BFacatory struct {
}

func NewBFacatory() BFacatory {
	return BFacatory{}
}

func (b BFacatory) Made() {
	fmt.Println("made in B")
}

func main() {
	var f Factory
	f = NewFactory(1)
	f.Made()

	f = NewFactory(2)
	f.Made()
}

package main

import "fmt"

//造车需要的流程
type MadeCar interface {
	SetEngine() //设置引擎发动机
	SetWheel()  //设置轮子
	Result()    //造成了
}

//造车的工厂
type CarFactory interface {
	Create() MadeCar
}

//A 厂
type AFactory struct {
}

func NewAFactory() AFactory {
	return AFactory{}
}

//A 厂可以提供一个造车服务
func (AFactory) Create() MadeCar {
	return ACarFactory{}
}

// A 厂的造车部门
type ACarFactory struct{}

func (ACarFactory) SetEngine() {
	fmt.Println("A begin to set engine")
}

func (ACarFactory) SetWheel() {
	fmt.Println("A begin to set wheel")
}

//A 厂造车
func (a ACarFactory) Result() {
	a.SetEngine()
	a.SetWheel()
	fmt.Println("made in A")
}

//B 厂
type BFactory struct {
}

func NewBFactory() BFactory {
	return BFactory{}
}

//B 厂可以提供一个造车服务
func (BFactory) Create() MadeCar {
	return BCarFactory{}
}

// B 厂的造车部门
type BCarFactory struct{}

func (BCarFactory) SetEngine() {
	fmt.Println("B begin to set engine")
}

func (BCarFactory) SetWheel() {
	fmt.Println("B begin to set wheel")
}

//B 厂造车
func (b BCarFactory) Result() {
	b.SetEngine()
	b.SetWheel()
	fmt.Println("made in B")
}

func main() {
	var f CarFactory
	f = NewAFactory() //分给 A 厂来造
	madeCar(f)

	f = NewBFactory() //分给 A 厂来造
	madeCar(f)

}

//开始造车
func madeCar(f CarFactory) {
	madeCar := f.Create()
	madeCar.Result()
}

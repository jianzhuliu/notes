package main

import "fmt"

//造引擎
type EngineApi interface {
	SetEngine()
}

//造轮子
type WheelApi interface {
	SetWheel()
}

//能造车的，就要能提供造引擎及造轮子的接口
type CarApi interface {
	CreateEngineApi() EngineApi
	CreateWheelApi() WheelApi
}

//A 厂来了
type AFactory struct {
}

func NewAFactory() AFactory {
	return AFactory{}
}

func (AFactory) CreateEngineApi() EngineApi {
	return AEngineApiImpl{}
}

func (AFactory) CreateWheelApi() WheelApi {
	return AWheelApiImpl{}
}

//A 厂造引擎实现接口
type AEngineApiImpl struct {
}

func (AEngineApiImpl) SetEngine() {
	fmt.Println("A is set engine")
}

//A 厂造轮子实现接口
type AWheelApiImpl struct {
}

func (AWheelApiImpl) SetWheel() {
	fmt.Println("A is set wheel")
}

func main() {
	var f CarApi
	f = NewAFactory()
	f.CreateEngineApi().SetEngine()
	f.CreateWheelApi().SetWheel()
}

package main

import "fmt"

//观察者必须实现统一接口，用于被观察者动态变化后通知
type ObserverInterface interface {
	Update(*Observable)
}

//被观察者
type Observable struct {
	Content   string              //需要更新的内容
	Observers []ObserverInterface //实现了统一接口的观察者列表
}

//创建一个被观察者对象
func NewObservable() *Observable {
	return &Observable{Observers: make([]ObserverInterface, 0)}
}

//添加观察者
func (observable *Observable) Attach(observerInterface ObserverInterface) {
	observable.Observers = append(observable.Observers, observerInterface)
}

//通知观察者
func (observable *Observable) Notify() {
	for _, observerInterface := range observable.Observers {
		observerInterface.Update(observable)
	}
}

//被观察者动态变化，开始广播了
func (observable *Observable) BroadCast(content string) {
	observable.Content = content

	//通知观察者
	observable.Notify()
}

//观察者
type Observer struct {
	Name string
}

//创建观察者
func NewObserver(name string) *Observer {
	return &Observer{Name: name}
}

func (observer *Observer) Update(observable *Observable) {
	fmt.Printf("%s receive %s \n", observer.Name, observable.Content)
}

func main() {
	//创建被观察者
	observable := NewObservable()

	//创建多个观察者
	observer1 := NewObserver("张三")
	observer2 := NewObserver("李四")
	observer3 := NewObserver("王五")

	//添加观察者
	observable.Attach(observer1)
	observable.Attach(observer2)
	observable.Attach(observer3)

	//被观察者动态变化
	observable.BroadCast("观察者模式")

}

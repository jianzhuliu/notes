/*
责任链模式
*/
package godesign

import (
	"log"
)

//责任接口
type IResponsibility interface {
	Do(*PaymentContext) error                //自身业务逻辑
	SetNext(IResponsibility) IResponsibility //设置下一个责任对象
	Run(*PaymentContext) error               //调用执行
}

//责任链对象
type ResponsibilityChain struct {
	next IResponsibility
}

func NewResponsibilityChain() *ResponsibilityChain {
	return &ResponsibilityChain{}
}

//公用部分接口方法实现
func (r *ResponsibilityChain) SetNext(responsibility IResponsibility) IResponsibility {
	r.next = responsibility
	return responsibility
}

func (r *ResponsibilityChain) Run(c *PaymentContext) error {
	if r.next != nil {
		//调用自身业务
		if err := r.next.Do(c); err != nil {
			return err
		}

		//如果有下一个责任，则调用下一个
		return r.next.Run(c)
	}

	return nil
}

//////////////////////////具体实现
//参数验证
type ResponsibilityParamCheck struct {
	ResponsibilityChain
}

func (r *ResponsibilityParamCheck) Do(c *PaymentContext) error {
	log.Println("参数校验")
	return nil
}

//支付
type ResponsibilityPay struct {
	ResponsibilityChain
}

func (r *ResponsibilityPay) Do(c *PaymentContext) error {
	log.Println("支付")
	return nil
}

//发奖
type ResponsibilityGift struct {
	ResponsibilityChain
}

func (r *ResponsibilityGift) Do(c *PaymentContext) error {
	log.Println("发奖")
	return nil
}

//发短信
type ResponsibilityMessage struct {
	ResponsibilityChain
}

func (r *ResponsibilityMessage) Do(c *PaymentContext) error {
	log.Println("发短信")
	return nil
}

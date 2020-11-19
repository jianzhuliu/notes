/*
策略模式+简单工厂模式
调用方，提供策略
*/
package godesign

import (
	"log"
)

type PayType int

//支付方式
const (
	PayTypeAli PayType = iota
	PayTypeWechat
	PayTypeBank
)

//支付上下文信息
type PaymentContext struct {
	PayType PayType
	Money   int
}

func NewPaymentContext(payType PayType, money int) *PaymentContext {
	return &PaymentContext{
		PayType: payType,
		Money:   money,
	}
}

//支付接口
type Istrategy interface {
	Pay(*PaymentContext) error
}

//策略执行者
type StrategyOperator struct {
	strategy Istrategy
}

func NewStrategyOperator() *StrategyOperator {
	return &StrategyOperator{}
}

//设置策略
func (o *StrategyOperator) SetStrategy(strategy Istrategy) {
	o.strategy = strategy
}

//统一由策略执行者调度
func (o *StrategyOperator) Pay(c *PaymentContext) error {
	return o.strategy.Pay(c)
}

//简单工厂模式，根据支付方式，获取支付策略
func GetStrategy(c *PaymentContext) Istrategy {
	switch c.PayType {
	case PayTypeAli:
		return &AliPay{}
	case PayTypeWechat:
		return &WechatPay{}
	case PayTypeBank:
		return &BankPay{}
	default:
		return nil
	}
}

////////////////////////////////具体实现
//微信支付
type WechatPay struct {
}

func (p *WechatPay) Pay(c *PaymentContext) error {
	log.Println("微信支付", c.Money)
	return nil
}

//支付宝支付
type AliPay struct {
}

func (p *AliPay) Pay(c *PaymentContext) error {
	log.Println("支付宝支付", c.Money)
	return nil
}

//银行支付
type BankPay struct {
}

func (p *BankPay) Pay(c *PaymentContext) error {
	log.Println("银行支付", c.Money)
	return nil
}

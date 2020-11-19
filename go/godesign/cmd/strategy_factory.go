package main

import (
	"fmt"
	"godesign"
)

func main() {
	paymentContext := godesign.NewPaymentContext(godesign.PayTypeAli, 30)
	strategy := godesign.GetStrategy(paymentContext)
	if strategy == nil {
		fmt.Println("没有对应的支付方式", paymentContext)
		return
	}

	strategyOperator := godesign.NewStrategyOperator()
	strategyOperator.SetStrategy(strategy)
	strategyOperator.Pay(paymentContext)
}

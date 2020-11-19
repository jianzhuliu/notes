package main

import (
	"fmt"
	"godesign"
)

func main() {
	paymentContext := godesign.NewPaymentContext(godesign.PayTypeAli, 30)

	responsibilityChain := godesign.NewResponsibilityChain()
	responsibilityChain.SetNext(&godesign.ResponsibilityParamCheck{}).
		SetNext(&godesign.ResponsibilityPay{}).
		SetNext(&godesign.ResponsibilityGift{}).
		SetNext(&godesign.ResponsibilityMessage{})

	err := responsibilityChain.Run(paymentContext)
	if err != nil {
		fmt.Println(err)
	}
}

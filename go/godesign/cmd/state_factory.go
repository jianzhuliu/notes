package main

import (
	"context"
	"godesign"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stateContext := godesign.NewStateContext(15998585234, "短信内容")
	stateManager := godesign.NewStateManager(ctx, time.Second)

	for i := 0; i < 5; i++ {
		stateManager.GetCurrentState().Send(stateContext)
		time.Sleep(time.Second)
	}
}

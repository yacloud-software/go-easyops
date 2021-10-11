package main

import (
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/tokens"
	"time"
)

func SimplePing() {
	initClient()
xloop:
	ctx := tokens.ContextWithToken()
	_, err := cl.SimplePing(ctx, &common.Void{})
	if err != nil {
		fmt.Printf("failed to ping: %s\n", err)
	}
	if *loop {
		time.Sleep(time.Duration(300) * time.Millisecond)
		goto xloop
	}

}

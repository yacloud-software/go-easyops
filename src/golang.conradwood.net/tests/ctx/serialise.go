package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	gctx "golang.conradwood.net/go-easyops/ctx"
)

func serialise() error {
	ctx := authremote.Context()
	s, err := gctx.SerialiseContextToString(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Context:\n%s\n", s)
	return nil
}

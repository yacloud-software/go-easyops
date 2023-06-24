package client

import (
	"context"
	"golang.conradwood.net/go-easyops/ctx"
)

func getContext() context.Context {
	cb := ctx.NewContextBuilder()
	return cb.ContextWithAutoCancel()
}

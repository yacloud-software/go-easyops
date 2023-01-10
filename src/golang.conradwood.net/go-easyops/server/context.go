package server

import (
	"context"
	"golang.conradwood.net/go-easyops/ctx"
)

func context_Background() context.Context {
	cb := ctx.NewContextBuilder()
	return cb.ContextWithAutoCancel()

}

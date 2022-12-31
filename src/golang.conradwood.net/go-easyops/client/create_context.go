package client

import (
	"context"
	"golang.conradwood.net/go-easyops/tokens"
)

func getContext() context.Context {
	return tokens.DISContextWithToken()
}

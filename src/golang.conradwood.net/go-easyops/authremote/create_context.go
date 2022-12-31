package authremote

import (
	"context"
	"golang.conradwood.net/go-easyops/tokens"
)

func getContext() context.Context {
	return tokens.DISContextWithToken()
}
func getContextWithTimeout(secs uint64) context.Context {
	return tokens.DISContextWithTokenAndTimeout(secs)
}

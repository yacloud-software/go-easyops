package rpc

import (
	"context"
	apb "golang.conradwood.net/apis/auth"
)

func GetUser(ctx context.Context) *apb.User {
	co := fromContext(ctx)
	return co.user
}

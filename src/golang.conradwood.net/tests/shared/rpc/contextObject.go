package rpc

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/tokens"
	"google.golang.org/grpc/metadata"
)

const (
	current_version = 2
)

type contextObject struct {
	user    *apb.User
	service *apb.User
}

func (co *contextObject) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, CTXKEY, co)
	newmd := metadata.Pairs(tokens.METANAME, co.serialize())
	ctx = metadata.NewOutgoingContext(ctx, newmd)
	return ctx
}

func fromContext(ctx context.Context) *contextObject {
	ifa := (ctx.Value(CTXKEY))
	co := ifa.(*contextObject)
	return co
}

// multiline description
func (co *contextObject) PrettyString() string {
	ud := ""
	if co.user != nil {
		ud = fmt.Sprintf("[email=%s, id=%s]", co.user.Email, co.user.ID)
	}
	sd := ""
	if co.service != nil {
		sd = fmt.Sprintf("[email=%s, id=%s]", co.service.Email, co.service.ID)
	}
	return fmt.Sprintf("User   : %s %s\nService: %s %s\n",
		auth.Description(co.user), ud,
		auth.Description(co.service), sd,
	)
}
func (co *contextObject) serialize() string {
	panic("THIS USED TO USE THE RPCINTERCEPTOR")
}

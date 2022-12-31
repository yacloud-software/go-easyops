package authremote

import (
	"context"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
)

func DIS_contextFromStruct(ctx context.Context, inctx *ge.InContext) context.Context {
	bs, err := utils.Marshal(inctx)
	if err != nil {
		panic("cannot marshal context")
	}
	newmd := metadata.Pairs(tokens.METANAME, bs)
	nctx := metadata.NewOutgoingContext(ctx, newmd)
	return nctx
}

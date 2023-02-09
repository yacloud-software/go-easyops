package ctxv1

import (
	"context"
	"fmt"
	//	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

var (
	ser_prefix = []byte("SER-CTX-V1")
)

func getPrefix() []byte {
	return ser_prefix
}
func Serialise(ctx context.Context) ([]byte, error) {
	ls := GetLocalState(ctx)
	ic := &ge.InContext{
		ImCtx: &ge.ImmutableContext{
			User:           ls.User(),
			CreatorService: ls.CreatorService(),
			RequestID:      ls.RequestID(),
			Session:        ls.Session(),
		},
		MCtx: &ge.MutableContext{
			CallingService: ls.CallingService(),
			Debug:          ls.Debug(),
			Trace:          ls.Trace(),
			Tags:           ls.RoutingTags(),
		},
	}
	var b []byte
	var err error
	b, err = utils.MarshalBytes(ic)
	if err != nil {
		return nil, err
	}

	prefix := ser_prefix
	b = append(prefix, b...)
	return b, nil
}
func DeserialiseWithTimeout(t time.Duration, buf []byte) (context.Context, error) {
	if len(buf) < len(ser_prefix) {
		return nil, fmt.Errorf("v1 context too short to deserialise (len=%d)", len(buf))
	}
	for i, b := range ser_prefix {
		if buf[i] != b {
			show := buf
			if len(show) > 10 {
				show = show[:10]
			}
			return nil, fmt.Errorf("v1 context has invalid prefix at pos %d (first 10 bytes: %s)", i, utils.HexStr(show))
		}
	}
	ud := buf[len(ser_prefix):]
	ctx := context.Background()
	shared.Debugf(ctx, "a v1deserialise: %s\n", utils.HexStr(buf))
	shared.Debugf(ctx, "b v1deserialise: %s\n", utils.HexStr(ud))
	ic := &ge.InContext{}
	err := utils.UnmarshalBytes(ud, ic)
	if err != nil {
		return nil, err
	}
	cb := &v1ContextBuilder{}
	if ic.ImCtx != nil {
		cb.WithUser(ic.ImCtx.User)
	} else {
		panic("no imctx")
	}
	if ic.MCtx != nil {
		cb.WithCallingService(ic.MCtx.CallingService)
	}
	return cb.ContextWithAutoCancel(), nil
}

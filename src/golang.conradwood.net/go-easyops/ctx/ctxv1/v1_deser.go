package ctxv1

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

var (
	ser_prefix = []byte("SER-CTX-V1")
)

func GetPrefix() []byte {
	return ser_prefix
}
func Serialise(ctx context.Context) ([]byte, error) {
	ls := GetLocalState(ctx)
	u := ls.User()
	var b []byte
	var err error
	if u != nil {
		b, err = utils.MarshalBytes(u)
		if err != nil {
			return nil, err
		}
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
	u := &auth.SignedUser{}
	err := utils.UnmarshalBytes(ud, u)
	if err != nil {
		return nil, err
	}
	cb := &v1ContextBuilder{}
	cb.WithUser(u)
	return cb.ContextWithAutoCancel(), nil
}

package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	//	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"os"
	"strings"
	"time"
)

const (
	SERBINPREFIX = "CTXUSER-BIN-"
	SERSTRPREFIX = "CTXUSER-STR-"
)

var (
//	rpci rc.RPCInterceptorServiceClient
)

// return a context with token and/or from environment or so
func Context(t time.Duration) context.Context {
	if tokens.GetUserTokenParameter() != "" || tokens.GetServiceTokenParameter() != "" {
		ctx := tokens.ContextWithTokenAndTimeout(uint64(t.Seconds()))
		return ctx
	}
	sctx := os.Getenv("GE_CTX")
	if sctx == "" {
		fmt.Fprintf(os.Stderr, "[go-easyops] Warning no context with tokens or environment at all\n")
		return nil
	}
	ctx, err := RecreateContextWithTimeout(t, []byte(sctx))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[go-easyops] Warning failed to recreate context from environment: %s\n", err)
		return nil
	}
	return ctx
}

// this will create a context for a userobject. if the userobject is signed, it will "just work"
func ContextForUser(u *auth.User) (context.Context, error) {
	if !common.VerifySignature(u) {
		return nil, fmt.Errorf("[go-easyops] no context (User signature invalid)")
	}
	token := tokens.GetServiceTokenParameter()
	mt := &rc.InMetadata{
		FooBar:       "local",
		ServiceToken: token,
		UserID:       u.ID,
		User:         u,
	}
	return contextForMeta(mt)
}
func contextForMetaWithTimeout(t time.Duration, mt *rc.InMetadata) (context.Context, error) {
	cs := &rpc.CallState{
		Started:  time.Now(),
		Debug:    true,
		Metadata: mt,
	}
	ctx := tokens.ContextWithToken()
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	err := cs.UpdateContextFromResponseWithTimeout(t)
	if err != nil {
		return nil, err
	}
	return cs.Context, nil

}
func contextForMeta(mt *rc.InMetadata) (context.Context, error) {
	cs := &rpc.CallState{
		Started:  time.Now(),
		Debug:    true,
		Metadata: mt,
	}
	ctx := tokens.ContextWithToken()
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	err := cs.UpdateContextFromResponse()
	if err != nil {
		return nil, err
	}
	return cs.Context, nil

}

// this returns a byte sequence, max 256 bytes long which may be used to recreate a users' context at some point in the future
func serialiseContextRaw(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("Cannot serialise 'nil' context")
	}
	md := tryGetMetadata(ctx)

	if md == nil {
		return nil, fmt.Errorf("No metadata in callstate")
	}
	data, err := proto.Marshal(md)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// this recreates a context from a previously stored state (see SerialiseContext())
func RecreateContextWithTimeout(t time.Duration, bs []byte) (context.Context, error) {
	uid := string(bs)
	var userdata []byte
	var err error
	if strings.HasPrefix(uid, SERSTRPREFIX) {
		s := string(bs[len(SERSTRPREFIX):])
		userdata, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(uid, SERBINPREFIX) {
		userdata = bs[len(SERBINPREFIX):]
	} else {
		return nil, fmt.Errorf("invalid serialised context prefix (%s)", uid)
	}

	md := &rc.InMetadata{}
	//	au := &auth.User{}
	err = proto.Unmarshal(userdata, md)
	if err != nil {
		return nil, err
	}
	if md.User != nil && !common.VerifySignature(md.User) {
		return nil, fmt.Errorf("[go-easyops] no context (User signature invalid)")
	}
	if md.Service != nil && !common.VerifySignature(md.Service) {
		return nil, fmt.Errorf("[go-easyops] no context (Service signature invalid)")
	}
	ctx, err := contextForMetaWithTimeout(t, md)
	return ctx, err
}

func SerialiseContextToString(ctx context.Context) (string, error) {
	b, err := serialiseContextRaw(ctx)
	if err != nil {
		return "", err
	}
	if len(b) == 0 {
		return "", fmt.Errorf("context serialised to zero bytes")
	}
	s := base64.StdEncoding.EncodeToString(b)
	if len(s) == 0 {
		return "", fmt.Errorf("context serialised to zero bytes")
	}
	return SERSTRPREFIX + s, nil
}

func SerialiseContext(ctx context.Context) ([]byte, error) {
	b, err := serialiseContextRaw(ctx)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("context serialised to zero bytes")
	}
	res := []byte(SERBINPREFIX)
	res = append(res, b...)
	return res, nil
}
func tryGetMetadata(ctx context.Context) *rc.InMetadata {
	cs := rpc.CallStateFromContext(ctx)
	if cs != nil {
		return cs.Metadata
	}
	u := GetUser(ctx)
	s := GetService(ctx)
	mt := &rc.InMetadata{
		FooBar:       "local",
		ServiceToken: tokens.GetServiceTokenParameter(),
		User:         u,
		Service:      s,
	}
	if mt.User != nil {
		mt.UserID = u.ID
	}
	return mt
}

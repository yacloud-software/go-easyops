package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	//"google.golang.org/protobuf/proto"
	//	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	pctx "golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"time"
)

const (
	SERBINPREFIX = "CTXUSER-BIN-"
	SERSTRPREFIX = "CTXUSER-STR-"
)

var (
// rpci rc.RPCInterceptorServiceClient
)

// return a context with token and/or from environment or so
// this function is obsolete and deprecated. use authremote.Context() instead
func DISContext(t time.Duration) context.Context {
	if tokens.GetUserTokenParameter() != "" || tokens.GetServiceTokenParameter() != "" {
		ctx := tokens.DISContextWithTokenAndTimeout(uint64(t.Seconds()))
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

func ForkContext(ictx context.Context) (context.Context, error) {
	if cmdline.ContextWithBuilder() {
		u := GetSignedUser(ictx)
		cb := pctx.NewContextBuilder()
		cb.WithUser(u)
		nctx := cb.ContextWithAutoCancel()
		return nctx, nil
	}
	u := GetSignedUser(ictx)
	return DISContextForSignedUser(u)
}

// this will create a context for a userobject. if the userobject is signed, it will "just work"
// this function is obsolete and deprecated. use authremote.Context() instead
func DISContextForSignedUser(su *auth.SignedUser) (context.Context, error) {
	if su == nil {
		return nil, fmt.Errorf("[go-easyops] ctxforuser: no user specified")
	}
	u := common.VerifySignedUser(su)
	if u == nil {
		return nil, fmt.Errorf("[go-easyops] ctxforuser: no context (User signature invalid)")
	}
	token := tokens.GetServiceTokenParameter()
	mt := &rc.InMetadata{
		FooBar:       "local",
		ServiceToken: token,
		UserID:       u.ID,
		SignedUser:   su,
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
	ctx := tokens.DISContextWithToken()
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
	ctx := tokens.DISContextWithToken()
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
	if cmdline.ContextWithBuilder() {
		return pctx.SerialiseContext(ctx)
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
	if pctx.IsSerialisedByBuilder(bs) {
		return pctx.DeserialiseContextWithTimeout(t, bs)
	}
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
		us := uid
		if len(us) > 10 {
			us = us[:10]
		}
		return nil, fmt.Errorf("invalid serialised context prefix (%s (%s))", us, utils.HexStr([]byte(us)))
	}

	md := &rc.InMetadata{}
	//	au := &auth.User{}
	err = proto.Unmarshal(userdata, md)
	if err != nil {
		return nil, err
	}
	if md.User != nil && !common.VerifySignature(md.User) {
		//fmt.Printf("user: %#v\n", md.User)
		return nil, fmt.Errorf("[go-easyops] no context (User signature invalid)")
	}
	if md.Service != nil && !common.VerifySignature(md.Service) {
		return nil, fmt.Errorf("[go-easyops] no context (Service signature invalid)")
	}
	ctx, err := contextForMetaWithTimeout(t, md)
	return ctx, err
}

func SerialiseContextToString(ctx context.Context) (string, error) {
	if cmdline.ContextWithBuilder() {
		return pctx.SerialiseContextToString(ctx)
	}
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
	if cmdline.ContextWithBuilder() {
		return pctx.SerialiseContext(ctx)
	}
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

// get signed session from context or nil if none
func GetSignedSession(ctx context.Context) *auth.SignedSession {
	cs := rpc.CallStateFromContext(ctx)
	if cs == nil {
		return nil
	}

	s := cs.SignedSession()
	if s == nil {
		return nil
	}
	if !common.VerifyBytes(s.Session, s.Signature) {
		return nil
	}
	return s
}

// get session token from context or "" if none
func GetSessionToken(ctx context.Context) string {
	s := GetSignedSession(ctx)
	if s == nil {
		return ""
	}
	sess := &auth.Session{}
	err := utils.UnmarshalBytes(s.Session, sess)
	if err != nil {
		fmt.Printf("[go-easyops] invalid session (%s)\n", err)
		return ""
	}
	return sess.Token

}

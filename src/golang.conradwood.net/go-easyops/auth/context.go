package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"golang.conradwood.net/apis/auth"
	//"google.golang.org/protobuf/proto"
	//	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	pctx "golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/tokens"
	"os"
	"time"
)

const (
	SERBINPREFIX = "CTXUSER-BIN-"
	SERSTRPREFIX = "CTXUSER-STR-"
)

// return a context with token and/or from environment or so
// this function is obsolete and deprecated. use authremote.Context() instead
func DISContext(t time.Duration) context.Context {
	if tokens.GetUserTokenParameter() != "" || tokens.GetServiceTokenParameter() != "" {
		ctx := tokens.DISContextWithTokenAndTimeout(uint64(t.Seconds()))
		return ctx
	}
	sctx := cmdline.GetEnvContext()
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
		s := GetSignedService(ictx)
		cb := pctx.NewContextBuilder()
		cb.WithUser(u)
		cb.WithCallingService(s)
		nctx := cb.ContextWithAutoCancel()
		return nctx, nil
	}
	u := GetSignedUser(ictx)
	if u == nil {
		return DISContext(time.Duration(10) * time.Second), nil
	}
	return DISContextForSignedUser(u)
}

// this will create a context for a userobject. if the userobject is signed, it will "just work"
// this function is obsolete and deprecated. use authremote.Context() instead
func DISContextForSignedUser(su *auth.SignedUser) (context.Context, error) {
	if su == nil {
		return nil, fmt.Errorf("[go-easyops] ctxforuser: no user specified")
	}
	panic("obsolete codepath")
}

// this returns a byte sequence, max 256 bytes long which may be used to recreate a users' context at some point in the future
func serialiseContextRaw(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("Cannot serialise 'nil' context")
	}
	if cmdline.ContextWithBuilder() && pctx.IsContextFromBuilder(ctx) {
		return pctx.SerialiseContext(ctx)
	}

	panic("obsolete codepath")
}

// this recreates a context from a previously stored state (see SerialiseContext())
func RecreateContextWithTimeout(t time.Duration, bs []byte) (context.Context, error) {
	if pctx.IsSerialisedByBuilder(bs) {
		return pctx.DeserialiseContextWithTimeout(t, bs)
	}
	panic("obsolete context to deserialise - context serialised by a version of go easyops no longer supported")
}

func SerialiseContextToString(ctx context.Context) (string, error) {
	if cmdline.ContextWithBuilder() && pctx.IsContextFromBuilder(ctx) {
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
	if cmdline.ContextWithBuilder() && pctx.IsContextFromBuilder(ctx) {
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

// get signed session from context or nil if none
func DISGetSignedSession(ctx context.Context) *auth.SignedSession {
	return nil
	/*
		ls := pctx.GetLocalState(ctx)
		res := ls.Session()
		if res != nil {
			return res
		}
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
	*/
}

// get session token from context or "" if none
func GetSessionToken(ctx context.Context) string {
	ls := pctx.GetLocalState(ctx)
	res := ls.Session()
	if res != nil {
		return res.SessionID
	}
	return ""
	/*
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
	*/
}

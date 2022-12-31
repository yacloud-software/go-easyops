package authremote

/****************************************************************

code in this file is not in use.



****************************************************************/

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	"time"
)

// a local context value
type CallStateV2 struct {
	inctx *ge.InContext
}

func DIS_ContextV2WithTimeoutAndTags(t time.Duration, rt *ge.CTXRoutingTags) context.Context {
	user := getLocalUserAccount()
	ctx, cnc := DIS_ContextV2WithTimeoutAndTagsForUser(t, "no_req_id", user, rt)
	go auto_cancel(cnc, t)
	return ctx
}

// automatically cancel after duration
func auto_cancel(cf context.CancelFunc, t time.Duration) {
	time.Sleep(t)
	cf()
}

/*
creates a new context for a user, with a timeout and routing tags and a cancel function
userid may be "" (empty).
*/
func DIS_ContextV2WithTimeoutAndTagsForUser(t time.Duration, reqid string, user *apb.SignedUser, rt *ge.CTXRoutingTags) (context.Context, context.CancelFunc) {
	if cmdline.IsStandalone() {
		f := func() {}
		return standalone_ContextWithTimeoutAndTags(t, rpc.Tags_ge_to_rpc(rt)), f
	}
	ctx, cnc := context.WithTimeout(context.Background(), t)
	inctx := DIS_build_new_ctx_meta_struct(reqid, user, nil)
	inctx.MCtx.Tags = rt
	lm := &CallStateV2{inctx: inctx}
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAMEV2, lm)
	ctx = DIS_contextFromStruct(ctx, inctx)
	return ctx, cnc
}

/*
build the struct we need to add to the context. used to create new contexts (e.g. in h2gproxy or in command line clients)
it will determine the service itself. user and sudo may be nil.
this is intented to be used as outbound context to other services
*/
func DIS_build_new_ctx_meta_struct(requestid string, user, sudo *apb.SignedUser) *ge.InContext {
	fmt.Printf("[go-easyops] Building meta for user %s\n", auth.Description(common.VerifySignedUser(user)))
	lsvc := getLocalServiceAccount()
	res := &ge.InContext{
		ImCtx: &ge.ImmutableContext{
			CreatorService: lsvc,
			RequestID:      requestid,
			User:           user,
			SudoUser:       sudo,
		},
		MCtx: &ge.MutableContext{
			CallingService: lsvc,
		},
	}
	return res
}

package ctx

import (
	"context"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/rpc"
	"time"
)

type v1LocalState struct {
	callstate      *rpc.CallState
	builder        *v1ContextBuilder
	callingservice *auth.SignedUser
	started        time.Time
}

func GetLocalState(ctx context.Context) *v1LocalState {
	v := ctx.Value(shared.LOCALSTATENAME)
	res, ok := v.(*v1LocalState)
	if !ok {
		return nil
	}
	return res
}
func (ls *v1LocalState) CreatorService() *auth.SignedUser {
	return nil
}
func (ls *v1LocalState) CallingService() *auth.SignedUser {
	if ls == nil {
		return nil
	}
	return ls.callingservice
}
func (ls *v1LocalState) Debug() bool {
	return false
}
func (ls *v1LocalState) Trace() bool {
	return false
}
func (ls *v1LocalState) User() *auth.SignedUser {
	if ls == nil || ls.callstate == nil || ls.callstate.Metadata == nil {
		return nil
	}
	return ls.callstate.Metadata.SignedUser
}
func (ls *v1LocalState) Session() *auth.SignedSession {
	return nil
}
func (ls *v1LocalState) RequestID() string {
	return "v1reqid"
}
func (ls *v1LocalState) RoutingTags() *ge.CTXRoutingTags {
	return nil
}

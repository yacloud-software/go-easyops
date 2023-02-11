package ctxv1

import (
	"context"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/rpc"
	"time"
)

type v1LocalState struct {
	this_is_v1_local_state string
	callstate              *rpc.CallState
	builder                *v1ContextBuilder
	callingservice         *auth.SignedUser
	routingtags            *ge.CTXRoutingTags
	started                time.Time
}

func GetLocalState(ctx context.Context) *v1LocalState {
	if ctx == nil {
		panic("cannot get localstate from nil context")
	}
	v := ctx.Value(shared.LOCALSTATENAME)
	if v == nil {
		shared.Debugf(ctx, "[go-easyops] ctxv1 warning, tried to extract localstate from context which is not a v1 context\n")
	}
	res, ok := v.(*v1LocalState)
	if !ok {
		return nil
	}
	return res
}
func (ls *v1LocalState) CreatorService() *auth.SignedUser {
	//v1 does not have a creator service
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
	if ls == nil {
		return nil
	}
	return ls.routingtags
}

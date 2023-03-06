package ctxv1

import (
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.yacloud.eu/apis/session"
	"time"
)

type v1LocalState struct {
	this_is_v1_local_state string
	callstate              *rpc.CallState
	builder                *v1ContextBuilder
	callingservice         *auth.SignedUser
	routingtags            *ge.CTXRoutingTags
	started                time.Time
	session                *session.Session
	requestid              string
}

func assertV1LocalStateImplementsInterface() shared.LocalState {
	return &v1LocalState{}
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
func (ls *v1LocalState) Session() *session.Session {
	if ls == nil {
		return nil
	}
	return ls.session
}
func (ls *v1LocalState) RequestID() string {
	if ls == nil {
		return ""
	}
	return ls.requestid
}
func (ls *v1LocalState) RoutingTags() *ge.CTXRoutingTags {
	if ls == nil {
		return nil
	}
	return ls.routingtags
}
func (ls *v1LocalState) Info() string {
	return "v1localstate"
}

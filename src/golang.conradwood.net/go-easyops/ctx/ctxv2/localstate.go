package ctxv2

import (
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
)

type localState struct {
	builder *contextBuilder
}

// this function only serves to assert that localState satisfies its interface (compile-error otherwise)
func assert_localstate_interface() shared.LocalState {
	return &localState{}
}
func (ls *localState) CreatorService() *auth.SignedUser {
	//v1 does not have a creator service
	return nil
}
func (ls *localState) CallingService() *auth.SignedUser {
	if ls == nil {
		return nil
	}
	return ls.builder.creatorservice
}
func (ls *localState) Debug() bool {
	return false
}
func (ls *localState) Trace() bool {
	return false
}
func (ls *localState) User() *auth.SignedUser {
	return ls.builder.user
}
func (ls *localState) Session() *auth.SignedSession {
	if ls == nil {
		return nil
	}
	return ls.builder.session
}
func (ls *localState) RequestID() string {
	return ls.builder.requestid
}
func (ls *localState) RoutingTags() *ge.CTXRoutingTags {
	if ls == nil {
		return nil
	}
	return ls.builder.routing_tags
}

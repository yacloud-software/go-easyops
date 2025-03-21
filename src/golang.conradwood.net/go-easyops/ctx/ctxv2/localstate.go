package ctxv2

import (
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.yacloud.eu/apis/session"
)

type localState struct {
	this_is_v2_local_state string
	builder                *contextBuilder
	callingservice         *auth.SignedUser // who called us (different from contextbuilder, which contains this service's id)
}

// this function only serves to assert that localState satisfies its interface (compile-error otherwise)
func assert_localstate_interface() shared.LocalState {
	return &localState{this_is_v2_local_state: "v2_localstate"}
}
func (ls *localState) CreatorService() *auth.SignedUser {
	if ls == nil || ls.builder == nil {
		return nil
	}
	return ls.builder.ge_context.ImCtx.CreatorService
}
func (ls *localState) CallingService() *auth.SignedUser {
	if ls == nil {
		return nil
	}
	return ls.callingservice
	//	return ls.builder.ge_context.MCtx.CallingService // NOTE: this is always 'us' (the service we are running at, because the inbound2outbound makes it so)
}
func (ls *localState) Info() string {
	if ls.builder == nil {
		return "nobuilder"
	}
	return "localstate_from_ctxv2_builder"
}
func (ls *localState) Experiments() []*ge.Experiment {
	if ls == nil || ls.builder == nil {
		return nil
	}
	return ls.builder.ge_context.MCtx.Experiments
}
func (ls *localState) Debug() bool {
	if ls == nil || ls.builder == nil {
		return false
	}
	return ls.builder.ge_context.MCtx.Debug
}
func (ls *localState) Trace() bool {
	if ls == nil || ls.builder == nil {
		return false
	}
	return ls.builder.ge_context.MCtx.Trace
}
func (ls *localState) User() *auth.SignedUser {
	return ls.builder.ge_context.ImCtx.User
}
func (ls *localState) SudoUser() *auth.SignedUser {
	return ls.builder.ge_context.ImCtx.SudoUser
}
func (ls *localState) Session() *session.Session {
	if ls == nil {
		return nil
	}
	return ls.builder.ge_context.ImCtx.Session
}
func (ls *localState) RequestID() string {
	return ls.builder.ge_context.ImCtx.RequestID
}
func (ls *localState) RoutingTags() *ge.CTXRoutingTags {
	if ls == nil {
		return nil
	}
	return ls.builder.ge_context.MCtx.Tags
}
func (ls *localState) Services() []*ge.ServiceTrace {
	return ls.builder.ge_context.MCtx.Services
}
func (ls *localState) AuthTags() []string {
	return ls.builder.ge_context.ImCtx.AuthTags
}

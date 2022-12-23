package rpc

import (
	abp "golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/common"
)

type LocalCallState interface {
	User() *abp.User
	CallingService() *abp.User
	SignedSession() *abp.SignedSession
	RoutingTags() *ge.CTXRoutingTags
	RequestID() string
}
type CallStateV2 struct {
	inCtx *ge.InContext
}

func NewCallStateV2(inctx *ge.InContext) LocalCallState {
	return &CallStateV2{inCtx: inctx}
}
func (cs *CallStateV2) RequestID() string {
	return "foorequestid"
}
func (cs *CallStateV2) User() *abp.User {
	return common.VerifySignedUser(cs.inCtx.ImCtx.User)
}
func (cs *CallStateV2) CallingService() *abp.User {
	return common.VerifySignedUser(cs.inCtx.MCtx.CallingService)
}
func (cs *CallStateV2) SignedSession() *abp.SignedSession {
	return nil
}
func (cs *CallStateV2) RoutingTags() *ge.CTXRoutingTags {
	return nil
}

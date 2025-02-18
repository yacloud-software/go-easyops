package shared

import (
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.yacloud.eu/apis/session"
)

type emptyLocalState struct {
	THIS_IS_EMPTY_LOCAL_STATE string // marker
}

func IsEmptyLocalState(ls LocalState) bool {
	_, f := ls.(*emptyLocalState)
	return f
}

func newEmptyLocalState() *emptyLocalState {
	return &emptyLocalState{THIS_IS_EMPTY_LOCAL_STATE: "this is empty local state"}
}
func (e *emptyLocalState) CreatorService() *auth.SignedUser { return nil }
func (e *emptyLocalState) CallingService() *auth.SignedUser { return nil }
func (e *emptyLocalState) Debug() bool                      { return false }
func (e *emptyLocalState) Trace() bool                      { return false }
func (e *emptyLocalState) User() *auth.SignedUser           { return nil }
func (e *emptyLocalState) SudoUser() *auth.SignedUser       { return nil }
func (e *emptyLocalState) Session() *session.Session        { return nil }
func (e *emptyLocalState) RequestID() string                { return "" }
func (e *emptyLocalState) RoutingTags() *ge.CTXRoutingTags  { return nil }
func (e *emptyLocalState) Info() string                     { return "EMPTY_LOCALSTATE" }
func (e *emptyLocalState) Experiments() []*ge.Experiment    { return nil }
func (e *emptyLocalState) Services() []*ge.ServiceTrace     { return nil }
func (e *emptyLocalState) GetAuthTags() []string            { return nil }

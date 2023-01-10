package ctx

import (
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
)

type EmptyLocalState struct {
}

func (e *EmptyLocalState) CreatorService() *auth.SignedUser { return nil }
func (e *EmptyLocalState) CallingService() *auth.SignedUser { return nil }
func (e *EmptyLocalState) Debug() bool                      { return false }
func (e *EmptyLocalState) Trace() bool                      { return false }
func (e *EmptyLocalState) User() *auth.SignedUser           { return nil }
func (e *EmptyLocalState) Session() *auth.SignedSession     { return nil }
func (e *EmptyLocalState) RequestID() string                { return "" }
func (e *EmptyLocalState) RoutingTags() *ge.CTXRoutingTags  { return nil }

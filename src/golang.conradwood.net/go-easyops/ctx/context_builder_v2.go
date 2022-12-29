package ctx

import (
	"context"
	"golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/goeasyops"
	"time"
)

// build V2 Contexts
type V2ContextBuilder struct {
}

/*
	return the context from this builder based on the options and WithXXX functions
*/
func (c *V2ContextBuilder) Context() (context.Context, context.CancelFunc) {
	return nil, nil
}

// automatically cancels context after timeout
func (c *V2ContextBuilder) ContextWithAuthCancel() context.Context {
	return nil
}

/*
add a user to context
*/
func (c *V2ContextBuilder) WithUser(user *auth.User) {
}

/*
add a creator service to context
*/
func (c *V2ContextBuilder) WithCreatorService(user *auth.User) {
}

/*
add a calling service (e.g. "me") to context
*/
func (c *V2ContextBuilder) WithCallingService(user *auth.User) {
}

/*
add a session to the context
*/
func (c *V2ContextBuilder) WithSession(user *auth.SignedSession) {
}

// mark context as with debug
func (c *V2ContextBuilder) WithDebug() {
}

// mark context as with trace
func (c *V2ContextBuilder) WithTrace() {
}
func (c *V2ContextBuilder) WithRoutingTags(pb.CTXRoutingTags) {
}
func (c *V2ContextBuilder) WithRequestID(reqid string) {
}
func (c *V2ContextBuilder) WithParentContext(context context.Context) {
}
func (c *V2ContextBuilder) WithTimeout(time.Duration) {
}

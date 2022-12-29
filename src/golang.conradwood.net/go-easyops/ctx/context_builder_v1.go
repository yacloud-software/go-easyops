package ctx

import (
	"context"
	"golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/goeasyops"
	"time"
)

// build V2 Contexts
type V1ContextBuilder struct {
}

/*
	return the context from this builder based on the options and WithXXX functions
*/
func (c *V1ContextBuilder) Context() (context.Context, context.CancelFunc) {
	return nil, nil
}

// automatically cancels context after timeout
func (c *V1ContextBuilder) ContextWithAuthCancel() context.Context {
	return nil
}

/*
add a user to context
*/
func (c *V1ContextBuilder) WithUser(user *auth.User) {
}

/*
add a creator service to context
*/
func (c *V1ContextBuilder) WithCreatorService(user *auth.User) {
}

/*
add a calling service (e.g. "me") to context
*/
func (c *V1ContextBuilder) WithCallingService(user *auth.User) {
}

/*
add a session to the context
*/
func (c *V1ContextBuilder) WithSession(user *auth.SignedSession) {
}

// mark context as with debug
func (c *V1ContextBuilder) WithDebug() {
}

// mark context as with trace
func (c *V1ContextBuilder) WithTrace() {
}
func (c *V1ContextBuilder) WithRoutingTags(pb.CTXRoutingTags) {
}
func (c *V1ContextBuilder) WithRequestID(reqid string) {
}
func (c *V1ContextBuilder) WithParentContext(context context.Context) {
}
func (c *V1ContextBuilder) WithTimeout(time.Duration) {
}

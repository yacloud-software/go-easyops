package ctx

import (
	"context"
	//	"fmt"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	"time"
)

const (
	METANAME = "goeasyops_meta" // marshaled proto, must match tokens.METANAME (avoiding import cycle)
)

// build V2 Contexts. That is, a context with metadata serialised into an rpc InContext struct
type V1ContextBuilder struct {
	requestid    string
	timeout      time.Duration
	parent       context.Context
	user         *auth.SignedUser
	service      *auth.SignedUser
	session      *auth.SignedSession
	routing_tags *ge.CTXRoutingTags
}

/*
return the context from this builder based on the options and WithXXX functions
*/
func (c *V1ContextBuilder) Context() (context.Context, context.CancelFunc) {
	cs := &rpc.CallState{
		Context: c.parent,
		RPCIResponse: &rc.InterceptRPCResponse{
			RequestID:           c.requestid,
			Source:              "ctxbuilder",
			SignedCallerUser:    c.user,
			SignedCallerService: c.service,
			CallerService:       common.VerifySignedUser(c.service),
			CallerUser:          common.VerifySignedUser(c.user),
		},
		Metadata: &rc.InMetadata{
			RequestID:     c.requestid,
			FooBar:        "foo_builder",
			SignedService: c.service,
			Service:       common.VerifySignedUser(c.service),
			SignedUser:    c.user,
			SignedSession: c.session,
			RoutingTags:   rpc.Tags_ge_to_rpc(c.routing_tags),
		},
	}
	//	fmt.Printf("Build with service: %s\n", describeUser(cs.Metadata.SignedService))
	cs.UpdateContextFromResponseWithTimeout(c.timeout)
	return cs.Context, cnc
}
func cnc() {
}

// automatically cancels context after timeout
func (c *V1ContextBuilder) ContextWithAutoCancel() context.Context {
	res, _ := c.Context()
	return res
}

/*
add a user to context
*/
func (c *V1ContextBuilder) WithUser(user *auth.SignedUser) {
	c.user = user
}

/*
add a creator service to context - v1 does not distinguish between creator and caller
*/
func (c *V1ContextBuilder) WithCreatorService(user *auth.SignedUser) {
	if user != nil {
		c.service = user
	}
}

/*
add a calling service (e.g. "me") to context
*/
func (c *V1ContextBuilder) WithCallingService(user *auth.SignedUser) {
	c.service = user
}

/*
add a session to the context - v1 does not have sessions
*/
func (c *V1ContextBuilder) WithSession(sess *auth.SignedSession) {
	c.session = sess
}

// mark context as with debug
func (c *V1ContextBuilder) WithDebug() {
}

// mark context as with trace
func (c *V1ContextBuilder) WithTrace() {
}
func (c *V1ContextBuilder) WithRoutingTags(tags *ge.CTXRoutingTags) {
	c.routing_tags = tags
}
func (c *V1ContextBuilder) WithRequestID(reqid string) {
	c.requestid = reqid
}
func (c *V1ContextBuilder) WithParentContext(context context.Context) {
	c.parent = context
}
func (c *V1ContextBuilder) WithTimeout(t time.Duration) {
	c.timeout = t
}

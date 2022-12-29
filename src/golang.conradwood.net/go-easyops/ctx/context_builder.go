/*
Package ctx contains methods to build authenticated contexts and retrieve information of them.
Package ctx is a "leaf" package - it is imported from many other goeasyops packages but does not import other goeasyops packages.

This package supports the following usecases:

* Create a new context with an authenticated user from a service (e.g. a web-proxy)

* Create a new context with a user and no service (e.g. a commandline)

* Create a new context from a service without a user (e.g. a service triggering a gRPC periodically)

* Update a context service (unary inbound gRPC interceptor)

Furthermore, go-easyops in general will parse "latest" context version and "latest-1" context versions. That is so that functions such as auth.User(context) return the right thing wether or not called from a service that has been updated or from a service that is not yet on latest. The context version it generates is selected via cmdline switches.

The context returned is ready to be used for outbound calls as-is.
The context also includes a "value" which is only available locally (does not cross gRPC boundaries) but is used to cache stuff.

*/
package ctx

import (
	"context"
	"golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/cmdline"
	"time"
)

type ContextBuilder interface {

	/*
		return the context from this builder based on the options and WithXXX functions
	*/
	Context() (context.Context, context.CancelFunc)

	// like Context(), but automatically call the CancelFunc after timeout
	ContextWithAuthCancel() context.Context

	/*
	   add a user to context
	*/
	WithUser(user *auth.User)

	/*
	   add a creator service to context
	*/
	WithCreatorService(user *auth.User)

	/*
	   add a calling service (e.g. "me") to context
	*/
	WithCallingService(user *auth.User)

	/*
	   add a session to the context
	*/
	WithSession(user *auth.SignedSession)

	// mark context as with debug
	WithDebug()

	// mark context as with trace
	WithTrace()
	// add routing tags
	WithRoutingTags(pb.CTXRoutingTags)
	//set the requestid
	WithRequestID(reqid string)
	// set a timeout for this context
	WithTimeout(time.Duration)
	// set a parent context
	WithParentContext(context context.Context)
}

// get a new contextbuilder
func NewContextBuilder() ContextBuilder {
	if cmdline.ContextV2() {
		return &V2ContextBuilder{}
	}
	return &V1ContextBuilder{}
}

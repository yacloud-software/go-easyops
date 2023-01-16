/*
Package ctx contains methods to build authenticated contexts and retrieve information of them.
Package ctx is a "leaf" package - it is imported from many other goeasyops packages but does not import (many) other goeasyops packages.

This package supports the following usecases:

* Create a new context with an authenticated user from a service (e.g. a web-proxy)

* Create a new context with a user and no service (e.g. a commandline)

* Create a new context from a service without a user (e.g. a service triggering a gRPC periodically)

* Update a context service (unary inbound gRPC interceptor)

Furthermore, go-easyops in general will parse "latest" context version and "latest-1" context versions. That is so that functions such as auth.User(context) return the right thing wether or not called from a service that has been updated or from a service that is not yet on latest. The context version it generates is selected via cmdline switches.

The context returned is ready to be used for outbound calls as-is.
The context also includes a "value" which is only available locally (does not cross gRPC boundaries) but is used to cache stuff.

Definition of CallingService: the LocalValue contains the service who called us. The context metadata contains this service definition (which in then is transmitted to downstream services)
*/
package ctx

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

const (
	LOCALSTATENAME = "goeasysops_localstate"
)

var (
	debug = flag.Bool("ge_debug_context", false, "if true print context debug stuff")
)

// the local state, this is not transmitted across grpc boundaries.
type LocalState interface {
	CreatorService() *auth.SignedUser
	CallingService() *auth.SignedUser
	Debug() bool
	Trace() bool
	User() *auth.SignedUser
	Session() *auth.SignedSession
	RequestID() string
	RoutingTags() *ge.CTXRoutingTags
}

type ContextBuilder interface {
	/*
		This function parses metadata found in an inbound context and, if successful, returns an "outbound" context with localstate.
		the bool return parameter indicates if it was successful(true) or not(false).
		Note that it requires the LOCAL service, because the calling service is modified and passed to the next service
	*/
	Inbound2Outbound(ctx context.Context, svc *auth.SignedUser) (context.Context, bool)
	/*
		return the context from this builder based on the options and WithXXX functions
	*/
	Context() (context.Context, context.CancelFunc)

	// like Context(), but automatically call the CancelFunc after timeout
	ContextWithAutoCancel() context.Context

	/*
	   add a user to context
	*/
	WithUser(user *auth.SignedUser)

	/*
	   add a creator service to context
	*/
	WithCreatorService(user *auth.SignedUser)

	/*
	   add a calling service (e.g. "me") to context
	*/
	WithCallingService(user *auth.SignedUser)

	/*
	   add a session to the context
	*/
	WithSession(user *auth.SignedSession)

	// mark context as with debug
	WithDebug()

	// mark context as with trace
	WithTrace()
	// add routing tags
	WithRoutingTags(*ge.CTXRoutingTags)
	//set the requestid
	WithRequestID(reqid string)
	// set a timeout for this context
	WithTimeout(time.Duration)
	// set a parent context for cancellation propagation (does not transfer metadata to the new context!)
	WithParentContext(context context.Context)
}

// get a new contextbuilder
func NewContextBuilder() ContextBuilder {
	return &v1ContextBuilder{}
}

// return "localstate" from a context. This is never "nil", but it is not guaranteed that the LocalState interface actually resolves details
func GetLocalState(ctx context.Context) LocalState {
	res := v1_getLocalState(ctx)
	if res == nil {
		if cmdline.ContextWithBuilder() {
			if *debug {
				utils.PrintStack("Localstate missing")
			}
			Debugf("could not get localstate from context (caller: %s)\n", utils.CallingFunction())
		}
		return &EmptyLocalState{}
	}
	return res
}

/*
we receive a context from gRPC (e.g. in a unary interceptor). To use this context for outbount calls we need to copy the metadata, we also need to add a local callstate for the fancy_picker/balancer/dialer. This is what this function does.
*/
func Inbound2Outbound(in_ctx context.Context, local_service *auth.SignedUser) context.Context {
	cb := &v1ContextBuilder{}
	octx, found := cb.Inbound2Outbound(in_ctx, local_service)
	if found {
		svc := common.VerifySignedUser(local_service)
		svs := "[none]"
		if svc != nil {
			svs = fmt.Sprintf("%s (%s)", svc.ID, svc.Email)
		}
		Debugf("converted inbound to outbound context (me.service=%s)\n", svs)
		return octx
	}
	fmt.Printf("[go-easyops] could not parse inbound context!\n")
	return in_ctx
}

func add_context_to_builder(cb ContextBuilder, ctx context.Context) {
	ls := GetLocalState(ctx)
	cb.WithCreatorService(ls.CreatorService())
	if ls.Debug() {
		cb.WithDebug()
	}
	cb.WithRequestID(ls.RequestID())
	if ls.Trace() {
		cb.WithTrace()
	}
	cb.WithUser(ls.User())
	cb.WithSession(ls.Session())
}

func Debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	s1 := fmt.Sprintf("[go-easyops] CONTEXT: ")
	s2 := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s", s1, s2)
}

// for debugging purposes we can convert a context to a human readable string
func Context2String(ctx context.Context) string {
	ls := GetLocalState(ctx)
	if ls == nil {
		return "[no localstate]"
	}
	return fmt.Sprintf("Localstate: %#v", ls)
}

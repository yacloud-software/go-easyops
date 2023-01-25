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
	//	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx/shared"
	v1 "golang.conradwood.net/go-easyops/ctx/v1"
	"golang.conradwood.net/go-easyops/utils"
	// "time"
)

var (
	debug = flag.Bool("ge_debug_context", false, "if true print context debug stuff")
)

// get a new contextbuilder
func NewContextBuilder() shared.ContextBuilder {
	return v1.NewContextBuilder()
}

// return "localstate" from a context. This is never "nil", but it is not guaranteed that the LocalState interface actually resolves details
func GetLocalState(ctx context.Context) shared.LocalState {
	res := v1.GetLocalState(ctx)
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
It is intented to convert any (supported) version of context into the current version of this package
*/
func Inbound2Outbound(in_ctx context.Context, local_service *auth.SignedUser) context.Context {
	cb := NewContextBuilder()
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

func add_context_to_builder(cb shared.ContextBuilder, ctx context.Context) {
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

func SerialiseContext(ctx context.Context) ([]byte, error) {
	return nil, fmt.Errorf("cannot serialisecontext builder context")
}

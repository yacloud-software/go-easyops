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
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx/ctxv1"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
	"time"
)

const (
	SER_PREFIX_STR = "CTX_SER"
)

var (
	debug = flag.Bool("ge_debug_context", false, "if true print context debug stuff")
)

// get a new contextbuilder
func NewContextBuilder() shared.ContextBuilder {
	return ctxv1.NewContextBuilder()
}

// return "localstate" from a context. This is never "nil", but it is not guaranteed that the LocalState interface actually resolves details
func GetLocalState(ctx context.Context) shared.LocalState {
	res := ctxv1.GetLocalState(ctx)
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

// check if 'buf' contains a context, serialised by the builder. a 'true' result implies that it can be deserialised from this package
func IsSerialisedByBuilder(buf []byte) bool {
	if len(buf) < 2 {
		return false
	}
	if strings.HasPrefix(string(buf), SER_PREFIX_STR) {
		// it was serialised by context_builder - as a string
		return true
	}
	version := buf[0]
	buf = buf[1:]
	var b []byte
	if version == 1 {
		b = ctxv1.GetPrefix()
	} else {
		return false
	}

	if bytes.HasPrefix(buf, b) {
		return true
	}

	fmt.Printf("[go-easyops] Not a valid context (%s/%s)\n", string(ctxv1.GetPrefix()), string(buf))
	//	fmt.Printf("a: %s\n", utils.HexStr(b))
	//	fmt.Printf("b: %s\n", utils.HexStr(buf[:20]))
	return false
}

// serialise a context to bunch of bytes
func SerialiseContext(ctx context.Context) ([]byte, error) {
	version := byte(1) // to de-serialise later
	b, err := ctxv1.Serialise(ctx)

	if err != nil {
		return nil, err
	}
	b = append([]byte{version}, b...)
	return b, nil
}

// serialise a context to bunch of bytes
func SerialiseContextToString(ctx context.Context) (string, error) {
	b, err := SerialiseContext(ctx)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(b)
	s = SER_PREFIX_STR + s
	return s, nil
}

// this unmarshals a context from a string into a context
func DeserialiseContextFromString(s string) (context.Context, error) {
	if !strings.HasPrefix(s, SER_PREFIX_STR) {
		return nil, fmt.Errorf("not a valid string to deserialise into a context")
	}
	s = strings.TrimPrefix(s, SER_PREFIX_STR)
	userdata, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return DeserialiseContext(userdata)
}

// this unmarshals a context from a binary blob into a context
func DeserialiseContext(buf []byte) (context.Context, error) {
	if len(buf) < 2 {
		return nil, fmt.Errorf("invalid byte array to deserialise into a context")
	}
	version := buf[0]
	buf = buf[1:]
	var err error
	var res context.Context
	if version == 1 {
		res, err = ctxv1.DeserialiseWithTimeout(time.Duration(10)*time.Second, buf)
	} else {
		return nil, fmt.Errorf("attempt to deserialise incompatible version (%d) to context", version)
	}
	return res, err
}

// this unmarshals a context from a binary blob into a context
func DeserialiseContextWithTimeout(t time.Duration, buf []byte) (context.Context, error) {
	if len(buf) < 2 {
		return nil, fmt.Errorf("invalid byte array to deserialise into a context")
	}
	version := buf[0]
	buf = buf[1:]
	var err error
	var res context.Context
	if version == 1 {
		res, err = ctxv1.DeserialiseWithTimeout(t, buf)
	} else {
		return nil, fmt.Errorf("attempt to deserialise incompatible version (%d) to context", version)
	}
	return res, err
}

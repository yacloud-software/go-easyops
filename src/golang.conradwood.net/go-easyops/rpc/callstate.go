/*
deprecated. See ctx package for a replacement implementation.
*/
package rpc

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.org/x/net/context"
	"time"
)

const (
	LOCALCONTEXTNAME   = "GOEASYOPS_LOCALCTX"
	LOCALCONTEXTNAMEV2 = "GOEASYOPS_LOCALCTX_V2"
	// timeout for newly created context in this package
	DEFAULT_TIMEOUT_SECS = 10
)

var (
	moan_about_no_auth = flag.Bool("ge_debug_print_old_signatures", false, "if true print services which are calling this module with old-style signatures")
	userSourceMetric   = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_user_source",
			Help: "V=1 UNIT=ops DESC=RPC requests with users, split by source",
		},
		[]string{"servicename", "method", "source"},
	)
)

func init() {
	prometheus.MustRegister(userSourceMetric)
}

// information about our local call stack.
// this should always have been an interface, really..
type CallState struct {
	// v1
	Debug           bool
	Started         time.Time
	ServiceName     string          // this is our servicename (the one this module exports)
	MethodName      string          // the method that was called (the one we are implementing)
	CallingMethodID uint64          // who called us (maybe 0)
	Context         context.Context // guaranteed to be the most "up-to-date" context.
	MyServiceID     uint64          // that is us (our local serviceid)
	userCounted     bool
	// v2
	v2            LocalCallState // might be nil for old stuff
	isV2          bool
	signedUser    *auth.SignedUser
	signedService *auth.SignedUser
}

// this now is a "V2" style callstate
func (c *CallState) SetV2(lcs LocalCallState) {
	c.v2 = lcs
	c.isV2 = true
}
func (c *CallState) IsV2() bool {
	return c.isV2
}
func EnableDebug(ctx context.Context) {
	cs := CallStateFromContext(ctx)
	if cs == nil {
		return
	}
	cs.Debug = true
}

func (cs *CallState) RequestID() string {
	if cs.IsV2() {
		return cs.v2.RequestID()
	}
	panic("obsolete codepath")
}

func (cs *CallState) SignedUser() *auth.SignedUser {
	return cs.signedUser

}
func (cs *CallState) SignedService() *auth.SignedUser {
	return cs.signedService
}

// return the authenticated user
func (cs *CallState) User() *auth.User {
	if cs == nil {
		return nil
	}
	if cs.IsV2() {
		return cs.v2.User()
	}
	panic("obsolete codepath")
}

func (cs *CallState) userbysource(src string) {
	if cs.userCounted {
		return
	}
	if *moan_about_no_auth && src != "signedv2" && src != "none" {
		sn := "undef"
		s := cs.CallerService()
		if s != nil {
			sn = fmt.Sprintf("%s/%s", s.ID, s.Email)
		}
		fmt.Printf("[go-easyops] service %s called us with old style signature \"%s\"\n", sn, src)
	}
	l := prometheus.Labels{"servicename": cs.ServiceName, "method": cs.MethodName, "source": src}
	userSourceMetric.With(l).Inc()
}

// return the authenticated service
func (cs *CallState) CallerService() *auth.User {
	if cs == nil {
		return nil
	}
	if cs.IsV2() {
		return cs.v2.CallingService()
	}
	panic("obsolete codepath")
}
func (cs *CallState) TargetString() string {
	if cs == nil {
		return ""
	}
	return fmt.Sprintf("%s.%s", cs.ServiceName, cs.MethodName)
}
func (cs *CallState) CallerString() string {
	if cs == nil {
		return ""
	}
	u := common.VerifySignedUser(cs.SignedUser())
	if u == nil {
		return ""
	}
	return u.Abbrev

}

// print context (if debug enabled)
func (cs *CallState) DebugPrintContext() {
	if cs == nil || !cs.Debug {
		return
	}
	cs.PrintContext()
}
func (cs *CallState) RoutingTags() *ge.CTXRoutingTags {
	if cs == nil {
		return nil
	}
	if cs.IsV2() {
		return cs.v2.RoutingTags()
	}
	panic("obsolete codepath")
}
func (cs *CallState) PrintContext() {
	if cs == nil {
		fmt.Printf("[go-easyops] Context has no Callstate\n")
		return
	}
	fmt.Printf("[go-easyops] printing old style context (%v)", cs.v2)
	/*
		lcv := CallStateFromContext(cs.Context)
		ls := "missing"
		if lcv != nil {
			if lcv.Metadata == nil {
				ls = "metadata missing"
			} else {
				ls = fmt.Sprintf("present (requestid: \"%s\")", lcv.Metadata.RequestID)
			}
		}
		fmt.Printf("[go-easyops] Local Context value: %s\n", ls)
		md, ex := metadata.FromIncomingContext(cs.Context)
		if ex {
			fmt.Printf("[go-easyops] InboundMeta: %v\n", metaToString(md))
		} else {
			fmt.Printf("[go-easyops] InboundMeta: NONE\n")
		}
		md, ex = metadata.FromOutgoingContext(cs.Context)
		if ex {
			fmt.Printf("[go-easyops] OutboundMeta: %v\n", metaToString(md))
		} else {
			fmt.Printf("[go-easyops] OutboundMeta: NONE\n")
		}
	*/
}
func desc(u *auth.User) string {
	if u == nil {
		return "NONE"
	}
	return fmt.Sprintf("%s[%s]", u.ID, u.Email)
}
func (cs *CallState) DISMetadataValue() string {
	return ""
}
func CallStateFromContext(ctx context.Context) *CallState {
	if ctx == nil {
		return nil
	}
	lcv := ctx.Value(LOCALCONTEXTNAME)
	if lcv == nil {
		return nil
	}
	return lcv.(*CallState)
}

// add the information from the interceptor response into the context
// it replaces the context data with information held in callstate
// it adds it in two distinct ways:
// 1. Accessible as a local context Variable
// 2. Ready to be transmitted via metadata (in case the context is use for outbound calls)
func (cs *CallState) UpdateContextFromResponse() error {
	return cs.UpdateContextFromResponseWithTimeout(time.Duration(*tokens.Deadline) * time.Second)
}
func (cs *CallState) UpdateContextFromResponseWithTimeout(t time.Duration) error {
	panic("obsolete codepath")
}

func ContextWithCallState(ctx context.Context) (context.Context, *CallState) {
	cs := &CallState{}
	nc := context.WithValue(ctx, LOCALCONTEXTNAME, cs)
	return nc, cs
}
func (cs *CallState) SignedSession() *auth.SignedSession {
	if cs.IsV2() {
		return cs.v2.SignedSession()
	}
	panic("obsolete codepath")
}

package rpc

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	LOCALCONTEXTNAME = "GOEASYOPS_LOCALCTX"
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

// information about our local call stack
type CallState struct {
	Debug           bool
	Started         time.Time
	ServiceName     string // this is our servicename (the one this module exports)
	MethodName      string // the method that was called (the one we are implementing)
	CallingMethodID uint64 // who called us (maybe 0)
	RPCIResponse    *rc.InterceptRPCResponse
	Metadata        *rc.InMetadata
	Context         context.Context // guaranteed to be the most "up-to-date" context.
	MyServiceID     uint64          // that is us (our local serviceid)
	userCounted     bool
}

func EnableDebug(ctx context.Context) {
	cs := CallStateFromContext(ctx)
	if cs == nil {
		return
	}
	cs.Debug = true
}

func (cs *CallState) RequestID() string {
	im := cs.Metadata
	if im == nil {
		return ""
	}
	return im.RequestID
}

// return the authenticated user
func (cs *CallState) User() *auth.User {
	if cs == nil {
		return nil
	}
	// signedv2 user available? if so return
	if cs.Metadata != nil && cs.Metadata.SignedUser != nil {
		cs.userbysource("signedv2")
		return common.VerifySignedUser(cs.Metadata.SignedUser)
	}
	// signed user available? if so return
	if cs.Metadata != nil && common.VerifySignature(cs.Metadata.User) {
		cs.userbysource("signed")
		return cs.Metadata.User
	}

	if cs.RPCIResponse != nil && cs.RPCIResponse.CallerUser != nil {
		cs.userbysource("rpcinterceptor")
		return cs.RPCIResponse.CallerUser
	}
	if cs.Metadata != nil && cs.Metadata.UserID != "" {
		cs.userbysource("bug")
		fmt.Printf("[go-easyops] cs.RPCIResponse=%#v, cs.Metadata=#%v\n", cs.RPCIResponse, cs.Metadata)
		// don't return UserID (we need user object)
		fmt.Printf("[go-easyops] BUG BUG This should never happen (found a userid in metadata but no RPCResponse.CallerUser)\n")
		return nil
	}
	cs.userbysource("none")
	return nil
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
	if cs.RPCIResponse == nil {
		return nil
	}
	if cs.RPCIResponse.SignedCallerService != nil {
		return common.VerifySignedUser(cs.RPCIResponse.SignedCallerService)
	}
	return cs.RPCIResponse.CallerService
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
	if cs.RPCIResponse == nil {
		return "{caller: no response from rpcinterceptor yet}"
	}
	au := cs.RPCIResponse.CallerUser
	if au == nil {
		return "{ no caller identified }"
	}
	return au.Abbrev
}

// print context (if debug enabled)
func (cs *CallState) DebugPrintContext() {
	if cs == nil || !cs.Debug {
		return
	}
	cs.PrintContext()
}
func (cs *CallState) PrintContext() {
	if cs == nil {
		fmt.Printf("[go-easyops] Context has no Callstate\n")
		return
	}

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
}
func metaToString(md metadata.MD) string {
	mdsa := md[tokens.METANAME]
	if len(mdsa) > 1 {
		return fmt.Sprintf("[manymeta(%d)]", len(mdsa))
	}
	if len(mdsa) == 0 {
		return "[emptymeta]"
	}
	mds := mdsa[0]
	if mds == "" {
		return "[nometa]"
	}
	res := &rc.InMetadata{}
	err := utils.Unmarshal(mds, res)
	if err != nil {
		return fmt.Sprintf("META: error %s\n", err)
	}
	return fmt.Sprintf("%#v", res)
}
func (cs *CallState) MetadataValue() string {
	if cs.Metadata == nil {
		return ""
	}
	/*
		if cs.RPCIResponse == nil && cs.User == nil {
			return ""
		}
	*/
	s, err := utils.Marshal(cs.Metadata)
	if err != nil {
		fmt.Printf("[go-easyops] Warning, unable to marshal metadata: %s\n", err)
	}
	return s

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
	if cs.Context == nil {
		cs.Context, _ = context.WithTimeout(context.Background(), t)
	}
	v := cs.Context.Value(LOCALCONTEXTNAME)
	if v == nil {
		cs.Context = context.WithValue(cs.Context, LOCALCONTEXTNAME, cs)
	}
	newmd := metadata.Pairs(tokens.METANAME, cs.MetadataValue())
	cs.Context = metadata.NewOutgoingContext(cs.Context, newmd)
	/*
		md, exists := metadata.FromOutgoingContext(cs.Context)
		if !exists {
			cs.Context = metadata.AppendToOutgoingContext(cs.Context, tokens.METANAME, cs.MetadataValue())
		} else {
			md.Set(tokens.METANAME, cs.MetadataValue())
		}
	*/
	return nil
}

func ContextWithCallState(ctx context.Context) (context.Context, *CallState) {
	cs := &CallState{}
	nc := context.WithValue(ctx, LOCALCONTEXTNAME, cs)
	return nc, cs
}

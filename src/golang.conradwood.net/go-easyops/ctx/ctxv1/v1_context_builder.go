package ctxv1

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	METANAME = "goeasyops_meta" // marshaled proto, must match tokens.METANAME (avoiding import cycle)
)

var (
	debug = flag.Bool("debug_context_v1", false, "if true debug v1 context builder in more detail")
)

// build V2 Contexts. That is, a context with metadata serialised into an rpc InContext struct
type v1ContextBuilder struct {
	requestid    string
	timeout      time.Duration
	parent       context.Context
	got_parent   bool
	user         *auth.SignedUser
	service      *auth.SignedUser
	session      *auth.SignedSession
	routing_tags *ge.CTXRoutingTags
}

/*
return the context from this builder based on the options and WithXXX functions
*/
func (c *v1ContextBuilder) Context() (context.Context, context.CancelFunc) {
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
			FooBar:        "ctxv1_builder",
			SignedService: c.service,
			Service:       common.VerifySignedUser(c.service),
			SignedUser:    c.user,
			SignedSession: c.session,
			RoutingTags:   rpc.Tags_ge_to_rpc(c.routing_tags),
			User:          common.VerifySignedUser(c.user),
		},
	}
	//	fmt.Printf("Build with service: %s\n", describeUser(cs.Metadata.SignedService))
	//	cs.UpdateContextFromResponseWithTimeout(c.timeout)

	var ctx context.Context
	var cnc context.CancelFunc
	octx := c.parent
	if !c.got_parent {
		octx = context.Background()
	}
	if c.timeout != 0 {
		ctx, cnc = context.WithTimeout(octx, c.timeout)
	} else {
		ctx, cnc = context.WithCancel(octx)
	}
	ls := c.newLocalState(cs)
	ctx = context.WithValue(ctx, shared.LOCALSTATENAME, ls)
	newmd := metadata.Pairs(METANAME, cs.MetadataValue())
	ctx = metadata.NewOutgoingContext(ctx, newmd)

	return ctx, cnc
}

// automatically cancels context after timeout
func (c *v1ContextBuilder) ContextWithAutoCancel() context.Context {
	res, cnc := c.Context()
	if c.timeout != 0 {
		go autocanceler(c.timeout, cnc)
	}
	return res
}
func autocanceler(t time.Duration, cf context.CancelFunc) {
	time.Sleep(t)
	cf()
}

/*
add a user to context
*/
func (c *v1ContextBuilder) WithUser(user *auth.SignedUser) {
	c.user = user
}

/*
add a creator service to context - v1 does not distinguish between creator and caller
*/
func (c *v1ContextBuilder) WithCreatorService(user *auth.SignedUser) {
	if user != nil {
		c.service = user
	}
}

/*
add a calling service (e.g. "me") to context
*/
func (c *v1ContextBuilder) WithCallingService(user *auth.SignedUser) {
	c.service = user
}

/*
add a session to the context - v1 does not have sessions
*/
func (c *v1ContextBuilder) WithSession(sess *auth.SignedSession) {
	c.session = sess
}

// mark context as with debug
func (c *v1ContextBuilder) WithDebug() {
}

// mark context as with trace
func (c *v1ContextBuilder) WithTrace() {
}
func (c *v1ContextBuilder) WithRoutingTags(tags *ge.CTXRoutingTags) {
	c.routing_tags = tags
}
func (c *v1ContextBuilder) WithRequestID(reqid string) {
	c.requestid = reqid
}
func (c *v1ContextBuilder) WithParentContext(context context.Context) {
	c.parent = context
	c.got_parent = true
}
func (c *v1ContextBuilder) WithTimeout(t time.Duration) {
	c.timeout = t
}
func (c *v1ContextBuilder) newLocalState(cs *rpc.CallState) *v1LocalState {
	return &v1LocalState{this_is_v1_local_state: "v1localstate",
		callstate:      cs,
		builder:        c,
		callingservice: cs.Metadata.SignedService,
	}
}
func (c *v1ContextBuilder) Inbound2Outbound(ctx context.Context, svc *auth.SignedUser) (context.Context, bool) {
	/*
		if svc == nil {
			fmt.Printf("[go-easyops] WARNING, creating context from inbound without service\n")
		}
	*/
	// get the proto from metadata:
	md, ex := metadata.FromIncomingContext(ctx)
	if !ex {
		// no metadata at all
		return nil, false
	}
	mdas, fd := md[METANAME]
	if !fd || mdas == nil || len(mdas) != 1 {
		// got metadata, but not our key
		return nil, false
	}
	mds := mdas[0]
	res := &rc.InMetadata{}
	err := utils.Unmarshal(mds, res)
	if err != nil {
		fmt.Printf("[go-easyops] invalid metadata: %s\n", err)
		return nil, false
	}
	if *debug {
		fmt.Printf("[go-easyops] CONTEXT: inbound metadata: user=%s\n", shared.PrettyUser(res.SignedUser))
		fmt.Printf("[go-easyops] CONTEXT: inbound metadata: user=%v\n", res.User)
		fmt.Printf("[go-easyops] CONTEXT: inbound metadata: %v\n", res)
	}
	// calling service is overwritten - store
	cservice := res.SignedService
	// now create new "outbound" context
	c.requestid = res.RequestID
	c.service = svc
	c.WithUser(res.SignedUser)
	c.WithSession(res.SignedSession)
	c.WithParentContext(ctx)
	c.routing_tags = rpc.Tags_rpc_to_ge(res.RoutingTags)
	out_ctx, _ := c.Context()
	GetLocalState(out_ctx).callingservice = cservice

	return out_ctx, true
}
func NewContextBuilder() *v1ContextBuilder {
	return &v1ContextBuilder{}
}

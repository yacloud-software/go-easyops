/*
this context uses a go-easyops proto to store information.
*/

package ctxv2

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/utils"
	"golang.yacloud.eu/apis/session"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	METANAME = "goeasyops_meta_v2" // in this case a serialised ge.InContext proto
)

var (
	ser_prefix = []byte("SER-CTX-V2")
	debug      = flag.Bool("ge_debug_context_v2", false, "if true debug v2 context builder in more detail")
	do_panic   = flag.Bool("ge_panic_v2_on_error", false, "if true panic very often")
)

// build V2 Contexts. That is, a context with metadata serialised into an rpc InContext struct
type contextBuilder struct {
	requestid      string
	timeout        time.Duration
	parent         context.Context
	got_parent     bool
	user           *auth.SignedUser
	service        *auth.SignedUser
	creatorservice *auth.SignedUser
	session        *session.Session
	routing_tags   *ge.CTXRoutingTags
	debug          bool
	trace          bool
}

/*
return the context from this builder based on the options and WithXXX functions
*/
func (c *contextBuilder) Context() (context.Context, context.CancelFunc) {
	ctx, cf, _ := c.contextWithLocalState()
	return ctx, cf
}

/*
return the context from this builder based on the options and WithXXX functions
*/
func (c *contextBuilder) contextWithLocalState() (context.Context, context.CancelFunc, *localState) {
	inctx := &ge.InContext{
		ImCtx: &ge.ImmutableContext{
			RequestID:      c.requestid,
			CreatorService: c.creatorservice,
			User:           c.user,
			Session:        c.session,
		},
		MCtx: &ge.MutableContext{
			CallingService: c.service,
			Debug:          c.debug,
			Trace:          c.trace,
			Tags:           c.routing_tags,
		},
	}
	b, err := utils.Marshal(inctx)
	if err != nil {
		panic(fmt.Sprintf("[go-easyops] unable to marshal context: %s", err))
	}
	ctx, cf := c.buildInitialContext()
	ls := c.newLocalState()
	ctx = context.WithValue(ctx, shared.LOCALSTATENAME, ls)
	newmd := metadata.Pairs(METANAME, b)
	ctx = metadata.NewOutgoingContext(ctx, newmd)
	ls.callingservice = c.service
	panic_if_service_account(common.VerifySignedUser(inctx.ImCtx.User))
	return ctx, cf, ls
}

// build a context from background, parent or so
func (c *contextBuilder) buildInitialContext() (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cnc context.CancelFunc
	octx := c.parent
	if !c.got_parent {
		octx = context.Background()
	}
	if c.timeout != 0 {
		ctx, cnc = context.WithTimeout(context.Background(), c.timeout)
	} else {
		ctx, cnc = context.WithCancel(octx)
	}
	return ctx, cnc
}

// automatically cancels context after timeout
func (c *contextBuilder) ContextWithAutoCancel() context.Context {
	res, cnc := c.Context()
	if c.timeout != 0 && cnc != nil {
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
func (c *contextBuilder) WithUser(user *auth.SignedUser) {
	panic_if_service_account(common.VerifySignedUser(user))
	c.user = user
}

/*
add a creator service to context - v1 does not distinguish between creator and caller
*/
func (c *contextBuilder) WithCreatorService(user *auth.SignedUser) {
	if user != nil {
		c.service = user
	}
}

/*
add a calling service (e.g. "me") to context
*/
func (c *contextBuilder) WithCallingService(user *auth.SignedUser) {
	c.service = user
}

/*
add a session to the context - v1 does not have sessions
*/
func (c *contextBuilder) WithSession(sess *session.Session) {
	c.session = sess
}

// mark context as with debug
func (c *contextBuilder) WithDebug() {
	c.debug = true
}

// mark context as with trace
func (c *contextBuilder) WithTrace() {
	c.trace = true
}
func (c *contextBuilder) WithRoutingTags(tags *ge.CTXRoutingTags) {
	c.routing_tags = tags
}
func (c *contextBuilder) WithRequestID(reqid string) {
	c.requestid = reqid
}
func (c *contextBuilder) WithParentContext(context context.Context) {
	c.parent = context
	c.got_parent = true
}
func (c *contextBuilder) WithTimeout(t time.Duration) {
	c.timeout = t
}
func (c *contextBuilder) newLocalState() *localState {
	ls := &localState{builder: c}
	return ls
}
func (c *contextBuilder) Inbound2Outbound(ctx context.Context, svc *auth.SignedUser) (context.Context, bool) {
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
	res := &ge.InContext{}
	err := utils.Unmarshal(mds, res)
	if err != nil {
		fmt.Printf("[go-easyops] warning invalid inbound v2 context (%s)\n", err)
		return nil, false
	}
	cb := &contextBuilder{}
	cb.WithUser(res.ImCtx.User)
	cb.WithCreatorService(res.ImCtx.CreatorService)
	cb.WithCallingService(svc)
	cb.WithSession(res.ImCtx.Session)
	if res.MCtx.Debug {
		cb.WithDebug()
	}
	if res.MCtx.Trace {
		cb.WithTrace()
	}
	cb.WithRoutingTags(res.MCtx.Tags)
	cb.WithRequestID(res.ImCtx.RequestID)
	cb.WithParentContext(ctx)
	ctx, _, ls := cb.contextWithLocalState() // always has a parent context, which means it needs no auto-cancel, uses parent cancelfunc
	// localstate has a different calling service (the original one)
	ls.callingservice = res.MCtx.CallingService
	panic_if_service_account(common.VerifySignedUser(res.ImCtx.User))
	return ctx, true
}
func NewContextBuilder() *contextBuilder {
	return &contextBuilder{}
}

func metadata_to_ctx(md metadata.MD, found bool) (*ge.InContext, error) {
	if !found {
		return nil, nil
	}
	mdas, fd := md[METANAME]
	if !fd || mdas == nil || len(mdas) != 1 {
		// got metadata, but not our key
		return nil, nil
	}
	mds := mdas[0]
	res := &ge.InContext{}
	err := utils.Unmarshal(mds, res)
	if err != nil {
		//		fmt.Printf("[go-easyops] warning invalid inbound v2 context (%s)\n", err)
		return nil, err
	}
	panic_if_service_account(common.VerifySignedUser(res.ImCtx.User))
	return res, nil

}
func get_metadata(ctx context.Context) (*ge.InContext, error) {
	ic, err := metadata_to_ctx(metadata.FromIncomingContext(ctx))
	if err == nil && ic != nil {
		return ic, nil
	}
	ic, err = metadata_to_ctx(metadata.FromOutgoingContext(ctx))
	return ic, err
}
func Serialise(ctx context.Context) ([]byte, error) {
	ls := shared.GetLocalState(ctx)
	ic := &ge.InContext{
		ImCtx: &ge.ImmutableContext{
			User:           ls.User(),
			CreatorService: ls.CreatorService(),
			RequestID:      ls.RequestID(),
			Session:        ls.Session(),
		},
		MCtx: &ge.MutableContext{
			CallingService: ls.CallingService(),
			Debug:          ls.Debug(),
			Trace:          ls.Trace(),
			Tags:           ls.RoutingTags(),
		},
	}
	panic_if_service_account(common.VerifySignedUser(ic.ImCtx.User))
	var b []byte
	var err error
	b, err = utils.MarshalBytes(ic)
	if err != nil {
		return nil, err
	}

	prefix := ser_prefix
	b = append(prefix, b...)
	return b, nil
}

/*
		ge, err := get_metadata(ctx)
		if err != nil {
			return nil, err
		}
		if ge == nil {
			return nil, fmt.Errorf("[go-easyops] no metadata in context to serialise")
		}
		b, err := utils.MarshalBytes(ge)
		if err != nil {
			return nil, err
		}
		panic("cannot serialise v2 contexts yet")
	}
*/
func DeserialiseContextWithTimeout(t time.Duration, buf []byte) (context.Context, error) {
	if len(buf) < len(ser_prefix) {
		return nil, fmt.Errorf("v1 context too short to deserialise (len=%d)", len(buf))
	}
	for i, b := range ser_prefix {
		if buf[i] != b {
			show := buf
			if len(show) > 18 {
				show = show[:18]
			}
			fmt.Printf("\nEXPECTED: %s\n", utils.HexStr(ser_prefix))
			fmt.Printf("GOT     : %s\n", utils.HexStr(buf))
			return nil, fmt.Errorf("v2 context has invalid prefix at pos %d (first 10 bytes: %s)", i, utils.HexStr(show))
		}
	}
	ud := buf[len(ser_prefix):]
	ctx := context.Background()
	shared.Debugf(ctx, "a v2deserialise: %s", utils.HexStr(buf))
	shared.Debugf(ctx, "b v2deserialise: %s", utils.HexStr(ud))
	ic := &ge.InContext{}
	err := utils.UnmarshalBytes(ud, ic)
	if err != nil {
		return nil, err
	}
	cb := &contextBuilder{}
	if ic.ImCtx != nil {
		panic_if_service_account(common.VerifySignedUser(ic.ImCtx.User))
		cb.WithUser(ic.ImCtx.User)
	} else {
		panic("no imctx")
	}
	if ic.MCtx != nil {
		cb.WithCallingService(ic.MCtx.CallingService)
	}
	cb.WithTimeout(t)
	return cb.ContextWithAutoCancel(), nil
}

func panic_if_service_account(u *auth.User) {
	if u == nil {
		return
	}
	if u.ServiceAccount {
		if *do_panic {
			panic(fmt.Sprintf("attempt to build context with serviceaccount as user %s (%s)", u.ID, u.Email))
		}
		fmt.Printf("[go-easyops] WARNING -- creating context with user as serviceaccount (%s) (%s)\n", u.ID, u.Email)
	}
}

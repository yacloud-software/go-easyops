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
The context also includes a "value" which is only available locally (does not cross gRPC boundaries) but is used to provide information from the context.

Definition of CallingService: the LocalValue contains the service who called us. The context metadata contains this service definition (which in then is transmitted to downstream services)

Contexts are transformed on each RPC. typically this goes like this:

  - Context is created ( via ContextBuilder() or authremote.Context() ). This context has a localstate and OUTBOUND metadata.

  - Client calls an RPC

  - Server transforms inbound context into a new context and adds itself as callingservice ( ctx.inbound2outbound() ). The inbound context has no localstate and INBOUND metadata. The new context has a localstate and OUTBOUND metadata.

  - Server becomes client, calls another rpc..
*/
package ctx

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx/ctxv2"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/utils"

	//	"golang.yacloud.eu/apis/session"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	SER_PREFIX_STR = "CTX_SER_STRING"
)

var (
	SER_PREFIX_BYT = []byte("CTX_SER_BYTE")
)

// get a new contextbuilder
func NewContextBuilder() shared.ContextBuilder {
	i := cmdline.GetContextBuilderVersion()
	if i == 1 {
		panic("obsolete codepath")
	} else if i == 2 {
		return ctxv2.NewContextBuilder()
	} else {
		// hm....
		return ctxv2.NewContextBuilder()
	}
}

// return "localstate" from a context. This is never "nil", but it is not guaranteed that the LocalState interface actually resolves details
func GetLocalState(ctx context.Context) shared.LocalState {
	return shared.GetLocalState(ctx)
}

// returns all "known" contextbuilders. we use this for received contexts to figure out which version it is
func getAllContextBuilders() map[int]shared.ContextBuilder {
	return map[int]shared.ContextBuilder{
		2: ctxv2.NewContextBuilder(),
	}
}

/*
we receive a context from gRPC (e.g. in a unary interceptor). To use this context for outbound calls we need to copy the metadata, we also need to add a local callstate for the fancy_picker/balancer/dialer. This is what this function does.
It is intented to convert any (supported) version of context into the current version of this package
*/
func Inbound2Outbound(in_ctx context.Context, local_service *auth.SignedUser) context.Context {
	for version, cb := range getAllContextBuilders() {
		octx, found := cb.Inbound2Outbound(in_ctx, local_service)
		if found {
			svc := common.VerifySignedUser(local_service)
			svs := "[none]"
			if svc != nil {
				svs = fmt.Sprintf("%s (%s)", svc.ID, svc.Email)
			}
			cmdline.DebugfContext("converted inbound (version=%d) to outbound context (me.service=%s)", version, svs)
			cmdline.DebugfContext("New Context: %s", Context2String(octx))
			ls := GetLocalState(octx)
			if ls == nil || shared.IsEmptyLocalState(ls) {
				utils.PrintStack("[go-easyops] no localstate for newly created context")
				return nil
			}
			cmdline.DebugfContext("Localstate %s: %#v\n", ls.Info(), ls)
			cmdline.DebugfContext("Localstate Detail:\n%#s\n", shared.LocalState2string(ls))

			return octx
		}
	}
	cmdline.DebugfContext("[go-easyops] could not parse inbound context!")
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
	for _, e := range ls.Experiments() {
		cb.EnableExperiment(e.Name)
	}
	cb.WithUser(ls.User())
	cb.WithSession(ls.Session())
}

// md,source,version -> metadata source: 0=none,1=inbound,2=outbound
func getMetadataFromContext(ctx context.Context) (string, int, int) {
	source := 1
	md, ex := metadata.FromIncomingContext(ctx)
	if !ex {
		source = 2
		md, ex = metadata.FromOutgoingContext(ctx)
		if !ex {
			// no metadata at all
			return "", 0, 0
		}
	}

	mdas, fd := md[ctxv2.METANAME]
	if fd {
		if len(mdas) != 1 {
			return "", source, 2
		}
		return mdas[0], source, 2
	}

	return "", 0, 0
}
func shortSessionText(ls shared.LocalState, maxlen int) string {
	s := ls.Session()
	if s == nil {
		return "nosession"
	}
	sl := s.SessionID
	if len(sl) > maxlen {
		sl = sl[:maxlen]
	}
	return sl
}

func Context2DetailString(ctx context.Context) string {
	return "TODO context_builder.go"
}

// for debugging purposes we can convert a context to a human readable string
func Context2String(ctx context.Context) string {
	md, src, version := getMetadataFromContext(ctx)

	ls := GetLocalState(ctx)
	if ls == nil || shared.IsEmptyLocalState(ls) {
		return fmt.Sprintf("[no localstate] md[src=%d,version=%d]", src, version)
	}
	if ls.User() != nil || ls.CallingService() != nil {
		sesstxt := shortSessionText(ls, 20)
		return fmt.Sprintf("Localstate[userid=%s,callingservice=%s,session=%s] md[src=%d,version=%d]", shared.PrettyUser(ls.User()), shared.PrettyUser(ls.CallingService()), sesstxt, src, version)
	}
	if src == 0 {
		return fmt.Sprintf("no localstate, no metadata (%v)\n", ctx)
	}
	if version == 2 {
		res := &ge.InContext{}
		err := utils.Unmarshal(md, res)
		if err != nil {
			return fmt.Sprintf("v2 %d metadata invalid (%s)", src, err)
		}
		return fmt.Sprintf("v2 (%d) metadata: %#v %#v\n,ls=[%s]", src, res.ImCtx, res.MCtx, shared.LocalState2string(ls))
	} else if version == 1 {
		panic("unsupported context version")
	}
	return fmt.Sprintf("Unsupported metadata version %d\n", version)

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
	if bytes.HasPrefix(buf, SER_PREFIX_BYT) {
		return true
	}
	/*
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
	*/
	cmdline.DebugfContext("[go-easyops] Not a ctxbuilder context (%s)", utils.HexStr(buf))
	//	cmdline.DebugfContext(ctx,"a: %s", utils.HexStr(b))
	//	cmdline.DebugfContext(ctx,"b: %s", utils.HexStr(buf[:20]))
	return false
}

// serialise a context to bunch of bytes
func SerialiseContext(ctx context.Context) ([]byte, error) {
	if !IsContextFromBuilder(ctx) {
		utils.PrintStack("incompatible context")
		return nil, fmt.Errorf("cannot serialise a context which was not built by builder")
	}
	version := byte(2) // to de-serialise later
	b, err := ctxv2.Serialise(ctx)

	if err != nil {
		return nil, err
	}
	chk := shared.Checksum(b)
	b = append([]byte{version, chk}, b...)
	b = append(SER_PREFIX_BYT, b...)
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

// this unmarshals a context from a string into a context. Short for DeserialiseContextFromStringWithTimeout()
func DeserialiseContextFromString(s string) (context.Context, error) {
	return DeserialiseContextFromStringWithTimeout(time.Duration(10)*time.Second, s)
}

// this unmarshals a context from a binary blob into a context. Short for DeserialiseContextWithTimeout()
func DeserialiseContext(buf []byte) (context.Context, error) {
	return DeserialiseContextWithTimeout(time.Duration(10)*time.Second, buf)
}

// this unmarshals a context from a string into a context
func DeserialiseContextFromStringWithTimeout(t time.Duration, s string) (context.Context, error) {
	if !strings.HasPrefix(s, SER_PREFIX_STR) {
		return nil, fmt.Errorf("not a valid string to deserialise into a context")
	}
	s = strings.TrimPrefix(s, SER_PREFIX_STR)
	userdata, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return DeserialiseContextWithTimeout(t, userdata)
}

// this unmarshals a context from a binary blob into a context
func DeserialiseContextWithTimeout(t time.Duration, buf []byte) (context.Context, error) {
	if !IsSerialisedByBuilder(buf) {
		panic("context not serialised by builder")
	}
	if len(buf) < 2 {
		return nil, fmt.Errorf("invalid byte array to deserialise into a context")
	}
	cmdline.DebugfContext("Deserialising %s", utils.HexStr(buf))
	tbuf := buf[len(SER_PREFIX_BYT):]
	s := string(buf)
	if strings.HasPrefix(s, SER_PREFIX_STR) {
		// it's a string...
		return DeserialiseContextFromStringWithTimeout(t, s)
	}
	if !bytes.HasPrefix(buf, SER_PREFIX_BYT) {
		// it's not a byte
		return nil, fmt.Errorf("context does not have ser_prefix_byt (%s)", utils.HexStr(buf))
	}

	version := tbuf[0]
	chk := tbuf[1]
	tbuf = tbuf[2:]
	c := shared.Checksum(tbuf)
	if c != chk {
		cmdline.DebugfContext("ERROR IN CHECKSUM (%d vs %d)", c, chk)
	}
	cmdline.DebugfContext("deserialising from version %d\n", version)
	var err error
	var res context.Context
	if version == 1 {
		// trying to deser v1 as v2
		res, err = ctxv2.DeserialiseContextWithTimeout(t, tbuf)
	} else if version == 2 {
		res, err = ctxv2.DeserialiseContextWithTimeout(t, tbuf)
	} else {
		cmdline.DebugfContext("a: %s", utils.HexStr(buf))
		utils.PrintStack("incompatible version %d", version)
		return nil, fmt.Errorf("(2) attempt to deserialise incompatible version (%d) to context", version)
	}
	if err != nil {
		cmdline.DebugfContext("unable to create context (%s)\n", err)
		return nil, err
	}
	cerr := res.Err()
	if cerr != nil {
		if cerr != nil {
			fmt.Printf("[go-easyops] created faulty context\n")
		}
	}
	cmdline.DebugfContext("Deserialised context: %s\n", Context2String(res))
	return res, err
}

// returns true if this context was build by the builder
func IsContextFromBuilder(ctx context.Context) bool {
	if ctx.Value(shared.LOCALSTATENAME) != nil {
		return true
	}
	return false
}

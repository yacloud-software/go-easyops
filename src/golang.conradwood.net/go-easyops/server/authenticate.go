package server

import (
	"context"
	"fmt"

	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/auth"

	//	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	//	"golang.conradwood.net/go-easyops/ctx"
	"sync"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/peer"
)

var (
	debuglock  sync.Mutex
	gettingrpc = false
	rpclock    sync.Mutex
	//	disable_interceptor = flag.Bool("ge_disable_interceptor", false, "if true, will not use rpc interceptor for access checks (very experimental!)")
	//verify_interceptor  = flag.Bool("ge_verify_noninterceptor", true, "if true, will compare the non-interceptor with interceptor by doing the actual intercept call and comparing results")
)

/*
*********************************************************************
newest method of authentication...
*********************************************************************
*/
// return error if not allowed to access
func (sd *serverDef) checkAccess(octx context.Context, rc *rpccall) error {
	if sd.noAuth || cmdline.IsStandalone() {
		return nil
	}
	if auth.GetUser(octx) == nil && auth.GetService(octx) == nil {
		fmt.Printf("[go-easyops] access denied to %s/%s for no-user and no-service to service with auth requirement (caller:%s)\n", rc.ServiceName, rc.MethodName, utils.CallingFunction())
		return errors.Unauthenticated(octx, "denied for access with no user and no service to rpc with auth requirement")
	}
	return nil
}

// authenticate a user (and authorise access to this method/service)
func Authenticate(ictx context.Context, cs *rpc.CallState) error {
	panic("obsolete codepath")
}

/*
	func MetaFromContext(ctx context.Context) *rc.InMetadata {
		headers, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			fmt.Printf("[go-easyops] Warning - cannot extract metadata from context (peer=%s)\n", peerFromContext(ctx))
			return nil
		}
		ims := headers[tokens.METANAME]
		if ims == nil || len(ims) == 0 {
			fmt.Printf("[go-easyops] Warning - metadata in context is nil or 0 (peer=%s)\n", peerFromContext(ctx))
			return nil
		}
		res := &rc.InMetadata{}
		err := utils.Unmarshal(ims[0], res)
		if err != nil {
			fmt.Printf("[go-easyops] Warning - unable to unmarshal metadata (%s)\n", err)
			return nil
		}
		return res
	}
*/
func peerFromContext(ctx context.Context) string {
	s := ""
	t, ok := peer.FromContext(ctx)
	if ok && t != nil && t.Addr != nil {
		s = t.Addr.String()
	}
	return s
}

func username(user *apb.User) string {
	if user == nil {
		return "[nouser]"
	}
	return fmt.Sprintf("[#%s %s]", user.ID, user.Email)
}

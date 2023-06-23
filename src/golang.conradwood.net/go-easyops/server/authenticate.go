package server

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/auth"
	//	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	//	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"sync"
)

var (
	debuglock           sync.Mutex
	rpcclient           rc.RPCInterceptorServiceClient
	gettingrpc          = false
	rpclock             sync.Mutex
	disable_interceptor = flag.Bool("ge_disable_interceptor", false, "if true, will not use rpc interceptor for access checks (very experimental!)")
	verify_interceptor  = flag.Bool("ge_verify_noninterceptor", true, "if true, will compare the non-interceptor with interceptor by doing the actual intercept call and comparing results")
)

/*
*********************************************************************
newest method of authentication...
*********************************************************************
*/
// return error if not allowed to access
func (sd *serverDef) checkAccess(octx context.Context, rc *rpccall) error {
	if sd.NoAuth || cmdline.IsStandalone() {
		return nil
	}
	if auth.GetUser(octx) == nil && auth.GetService(octx) == nil {
		fmt.Printf("[go-easyops] access denied to %s/%s for no-user and no-service to service with auth requirement (caller:%s)\n", rc.ServiceName, rc.MethodName, utils.CallingFunction())
		return errors.Unauthenticated(octx, "denied for access with no user and no service to rpc with auth requirement")
	}
	return nil
}

/*
*********************************************************************
older (obsolete) methods of authentication...
*********************************************************************
*/
func initrpc() error {
	/*
		if gettingrpc {
			return fmt.Errorf("[go-easyops] (auth) RPCInterceptor unavailable")
		}
	*/
	if rpcclient != nil {
		return nil
	}
	rpclock.Lock()
	defer rpclock.Unlock()
	gettingrpc = true
	if rpcclient != nil {
		gettingrpc = false
		return nil
	}
	if rpcclient == nil {
		rpcclient = rc.NewRPCInterceptorServiceClient(client.Connect("rpcinterceptor.RPCInterceptorService"))
	}
	gettingrpc = false
	return nil
}

// authenticate a user (and authorise access to this method/service)
func Authenticate(ictx context.Context, cs *rpc.CallState) error {
	panic("obsolete codepath")
}

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

func peerFromContext(ctx context.Context) string {
	s := ""
	t, ok := peer.FromContext(ctx)
	if ok && t != nil && t.Addr != nil {
		s = t.Addr.String()
	}
	return s
}

func build_access_details(ctx context.Context, irr *rc.InterceptRPCRequest) (*rc.InterceptRPCResponse, error) {
	panic("obsolete codepath")
}

func compare_intercept_responses(res1, res2 *rc.InterceptRPCResponse) error {
	compare_intercept_user("callerservice", res1.CallerService, res2.CallerService)
	compare_intercept_user("calleruser", res1.CallerUser, res2.CallerUser)
	if res1.CalleeServiceID != res2.CalleeServiceID {
		return fmt.Errorf("[go-easyops] intercept check failed: Calleeservice1=%d vs CalleeServiceID2=%d", res1.CalleeServiceID, res2.CalleeServiceID)
	}
	compare_intercept_suser("signedcalleruser", res1.SignedCallerUser, res2.SignedCallerUser)
	compare_intercept_suser("signedcallerservice", res1.SignedCallerService, res2.SignedCallerService)
	return nil
}

func compare_intercept_suser(name string, su1, su2 *apb.SignedUser) {
	u1 := common.VerifySignedUser(su1)
	u2 := common.VerifySignedUser(su2)
	compare_intercept_user(name, u1, u2)

}
func compare_intercept_user(name string, u1, u2 *apb.User) {
	if u1 == nil && u2 == nil {
		return
	}
	u1s := "nil"
	u2s := "nil"
	if u1 != nil {
		u1s = fmt.Sprintf("U1: %s %s\n", u1.ID, u1.Email)
	}
	if u2 != nil {
		u2s = fmt.Sprintf("U2: %s %s\n", u2.ID, u2.Email)
	}
	if u1 == nil {
		fmt.Printf("%s %s\n", u1s, u2s)
		panic(fmt.Sprintf("[go-easyops] intercept user check \"%s\" mismatch (u1 is nil, u2 == %s)", name, u2s))
	}
	if u2 == nil {
		fmt.Printf("%s %s\n", u1s, u2s)
		panic(fmt.Sprintf("[go-easyops] intercept user check \"%s\" mismatch (u2 is nil, u1 == %s)", name, u1s))
	}
	if u1.ID != u2.ID {
		fmt.Printf("%s %s\n", u1s, u2s)
		panic(fmt.Sprintf("[go-easyops] intercept user check \"%s\" mismatch (%v vs %v)", name, u1, u2))
	}
}

func username(user *apb.User) string {
	if user == nil {
		return "[nouser]"
	}
	return fmt.Sprintf("[#%s %s]", user.ID, user.Email)
}

package server

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/common"
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
	debug_auth          = flag.Bool("ge_debug_auth", false, "debug grpc authentication")
	disable_interceptor = flag.Bool("ge_disable_interceptor", false, "if true, will not use rpc interceptor for access checks (very experimental!)")
	verify_interceptor  = flag.Bool("ge_verify_noninterceptor", true, "if true, will compare the non-interceptor with interceptor by doing the actual intercept call and comparing results")
)

func initrpc() error {
	if gettingrpc {
		return fmt.Errorf("[go-easyops] (auth) RPCInterceptor unavailable")
	}
	if rpcclient != nil {
		return nil
	}
	rpclock.Lock()
	defer rpclock.Unlock()
	gettingrpc = true
	if rpcclient != nil {
		return nil
	}
	if rpcclient == nil {
		rpcclient = rc.NewRPCInterceptorServiceClient(client.Connect("rpcinterceptor.RPCInterceptorService"))
	}
	gettingrpc = false
	return nil
}

// authenticate a user (and authorise access to this method/service)
func Authenticate(cs *rpc.CallState) error {
	if *debug_auth {
		cs.Debug = true
	}
	err := initrpc()
	if err != nil {
		return err
	}
	if cs.Debug {
		fmt.Printf("[go-easyops] Calling RPC Interceptor...\n")
	}

	cs.Metadata = MetaFromContext(cs.Context)
	if cs.Debug {
		fmt.Printf("[go-easyops] Inbound metadata: %#v\n", cs.Metadata)
	}

	// call the interceptor
	irr := &rc.InterceptRPCRequest{
		InMetadata: cs.Metadata,
		Service:    cs.ServiceName,
		Method:     cs.MethodName,
	}

	// preserve some of the inbound metadata information (before we overwrite it withour outbound data)
	if cs.Metadata != nil {
		verifySignatures(cs)
		cs.CallingMethodID = cs.Metadata.CallerMethodID
		if cs.Metadata.UserToken == "" && cs.Metadata.ServiceToken == "" && cs.Metadata.User == nil && cs.Metadata.Service == nil && cs.Metadata.SignedUser == nil {
			t, ok := peer.FromContext(cs.Context)
			if ok && t != nil && t.Addr != nil {
				irr.Source = t.Addr.String()
			}
			fmt.Printf("[go-easyops] no identification by caller whatsoever (from %s)\n", irr.Source)
		}
	} else {
	}

	var res *rc.InterceptRPCResponse
	if *disable_interceptor {
		res, err = build_access_details(cs.Context, irr)
	} else {
		res, err = rpcclient.InterceptRPC(cs.Context, irr)
	}
	if err != nil {
		if *debug_auth {
			fmt.Printf("[go-easyops] RPCInterceptor.InterceptRPC() failed: %s\n", utils.ErrorString(err))
		}
		return err
	}
	cs.RPCIResponse = res

	// copy /some/ responses to inmeta
	cs.Metadata.RequestID = cs.RPCIResponse.RequestID
	cs.Metadata.CallerMethodID = cs.RPCIResponse.CallerMethodID
	verifySignatures(cs)
	cs.Metadata.CallerServiceID = cs.MyServiceID
	// all subsequent rpcs propagate OUR servicetoken
	cs.Metadata.ServiceToken = tokens.GetServiceTokenParameter()
	cs.Metadata.FooBar = "authmoo"
	if *debug_auth {
		fmt.Printf("[go-easyops] metadata after rpc interceptor: %#v\n", cs.Metadata)
		fmt.Printf("[go-easyops] RPC Interceptor (reject=%t) said: %v\n", cs.RPCIResponse.Reject, cs.RPCIResponse)
	}
	return nil
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

/*
 signatures from response. verify and copy response to metadata
goes through all user and service accounts and invalid ones are removed
*/
func verifySignatures(cs *rpc.CallState) {
	// we got a response, so copy stuff across
	r := cs.RPCIResponse
	var sigu, sigs *apb.SignedUser
	var u, s *apb.User
	if r != nil {
		if r.SignedCallerUser != nil {
			sigu = r.SignedCallerUser
			u = common.VerifySignedUser(sigu)
			if u != nil {
				r.CallerUser = u
			} else {
				r.CallerUser = nil
				r.SignedCallerUser = nil
			}

		}
		if r.SignedCallerService != nil {
			sigs = r.SignedCallerService
			s = common.VerifySignedUser(sigs)
			if s != nil {
				r.CallerService = s
			} else {
				r.CallerService = nil
				r.SignedCallerService = nil
			}

		}
		if u == nil {
			u = r.CallerUser
		}
		if s == nil {
			s = r.CallerService
		}
	}
	m := cs.Metadata
	if m != nil {
		if u == nil {
			if m.SignedUser != nil {
				sigu = m.SignedUser
				u = common.VerifySignedUser(sigu)
				if u == nil {
					sigu = nil
				}
			} else if m.User != nil {
				u = m.User
			}
		}

		if s == nil {
			if m.SignedUser != nil {
				// this isn't right? it sets the user as the service
				/*
					sigs = m.SignedUser
					s = common.VerifySignedUser(sigs)
					if s == nil {
						sigs = nil
					}
				*/
			} else if m.Service != nil {
				s = m.Service
			}
		}
		cs.Metadata.User = u
		cs.Metadata.SignedUser = sigu
		cs.Metadata.Service = s
		cs.Metadata.SignedService = sigs
		if u != nil {
			cs.Metadata.UserID = u.ID
			if !common.VerifySignature(u) {
				fmt.Printf("[go-easyops] invalid user signature\n")
				cs.Metadata.User = nil
				cs.Metadata.SignedUser = nil
			}
		}
		if s != nil {
			if !common.VerifySignature(s) {
				fmt.Printf("[go-easyops] invalid service signature\n")
				cs.Metadata.Service = nil
				cs.Metadata.SignedService = nil
			}

		}
	}
}

func build_access_details(ctx context.Context, irr *rc.InterceptRPCRequest) (*rc.InterceptRPCResponse, error) {
	debuglock.Lock()
	defer debuglock.Unlock()
	md := irr.InMetadata
	if md == nil {
		panic("no metadata for request!!")
	}

	fmt.Printf("ServiceToken        : %s\n", md.ServiceToken)
	fmt.Printf("UserToken           : %s\n", md.UserToken)
	fmt.Printf("CallerService       : %s\n", username(md.Service))
	fmt.Printf("Signed-CallerService: %s\n", username(common.VerifySignedUser(md.SignedService)))

	fmt.Printf("Metadata: %#v\n", md)
	reqid := md.RequestID
	if reqid == "" {
		reqid = utils.RandomString(32)
	}
	res := &rc.InterceptRPCResponse{
		RequestID:           reqid,
		Reject:              false,
		RejectReason:        rc.RejectReason_NonSpecific,
		CallerService:       md.Service,
		CallerUser:          md.User,
		CalleeServiceID:     md.CallerServiceID,
		SignedCallerService: md.SignedService,
		SignedCallerUser:    md.SignedUser,
		Source:              "go-easyops_local",
	}
	if md.SignedUser == nil && md.User != nil {
		st := fmt.Sprintf("[go-easyops] Warning: received context with user but no signed user from service %v", md.Service)
		fmt.Println(st)
		if *verify_interceptor {
			panic(st)
		} else {
			return nil, fmt.Errorf("got user, but not signed user")
		}
	}

	if md.ServiceToken != "" {
		res.SignedCallerService = authremote.SignedGetByToken(ctx, md.ServiceToken)
		res.CallerService = common.VerifySignedUser(res.SignedCallerService)
	}
	if res.CalleeServiceID == 0 {
		if res.CallerService != nil {
			sr, err := get_service_id(ctx, res.CallerService.ID)
			if err != nil {
				return nil, err
			}
			res.CalleeServiceID = sr.ID
		}
	}
	if res.SignedCallerUser == nil {
		res.SignedCallerUser = authremote.SignedGetByToken(ctx, md.UserToken)
		res.CallerUser = common.VerifySignedUser(res.SignedCallerUser)
	}

	if res.CallerUser == nil {
		res.CallerUser = authremote.GetByToken(ctx, md.UserToken)
	}
	if res.CallerService == nil {
		res.CallerService = authremote.GetByToken(ctx, md.ServiceToken)
	}

	if !*verify_interceptor {
		return res, nil
	}
	// debugging, compare stuff
	r_res, err := rpcclient.InterceptRPC(ctx, irr)
	if err != nil {
		return nil, err
	}
	err = compare_intercept_responses(res, r_res)
	if err != nil {
		panic(fmt.Sprintf("Error on intercept response: %s\n", err))
		//		return nil, err
	}

	return res, nil
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

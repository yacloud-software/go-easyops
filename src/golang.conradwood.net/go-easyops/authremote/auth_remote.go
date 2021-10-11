package authremote

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
	"os"
	"sync"
	"time"
)

var (
	userbyidcache    = cache.NewResolvingCache("userbyid", time.Duration(60)*time.Second, 9999)
	userbyemailcache = cache.NewResolvingCache("userbyemail", time.Duration(60)*time.Second, 9999)
	userbytokencache = cache.NewResolvingCache("userbytoken", time.Duration(60)*time.Second, 9999)
	rpci             rc.RPCInterceptorServiceClient
	authServer       apb.AuthenticationServiceClient
	authServerLock   sync.Mutex
	authManager      apb.AuthManagerServiceClient
	authManagerLock  sync.Mutex
	contextRetrieved = false
	lastUser         *apb.SignedUser
	lastService      *apb.SignedUser
)

func Context() context.Context {
	client.GetSignatureFromAuth()
	return ContextWithTimeout(time.Duration(10) * time.Second)
}

/*
create a new context with routing tags (routing criteria to route to specific instances of a service)
*/
func NewContextWithRouting(kv map[string]string) context.Context {
	return DerivedContextWithRouting(Context(), kv)
}

/*
derive  a context with routing tags (routing criteria to route to specific instances of a service)
*/
func DerivedContextWithRouting(cv context.Context, kv map[string]string) context.Context {
	cv = context.WithValue(cv, "routingtags", kv)
	return cv
}

/* this context gives a context with a full userobject
todo so it _has_ to call external servers to get a signed userobject.
if started_by_autodeployer will use tokens.ContextWithToken()
else if environment variable with context, will use auth.Context() (with variable)
else create context by asking auth service for a signed user object
*/
func ContextWithTimeout(t time.Duration) context.Context {
	if cmdline.Datacenter() {
		return tokens.ContextWithTokenAndTimeout(uint64(t.Seconds()))
	}
	sctx := os.Getenv("GE_CTX")
	if sctx != "" {
		return auth.Context(t)
	}

	if !contextRetrieved {
		lastUser = SignedGetByToken(context.Background(), tokens.GetUserTokenParameter())
		lastService = SignedGetByToken(context.Background(), tokens.GetServiceTokenParameter())
		contextRetrieved = true
	}
	luid := ""
	if lastUser != nil {
		luid = common.VerifySignedUser(lastUser).ID
	}
	cs := &rpc.CallState{
		Metadata: &rc.InMetadata{
			UserID:        luid,
			Service:       common.VerifySignedUser(lastService),
			User:          common.VerifySignedUser(lastUser),
			SignedService: lastService,
			SignedUser:    lastUser,
		},
		RPCIResponse: &rc.InterceptRPCResponse{
			CallerUser:          common.VerifySignedUser(lastUser),
			CallerService:       common.VerifySignedUser(lastService),
			SignedCallerUser:    lastUser,
			SignedCallerService: lastService,
		},
	}
	err := cs.UpdateContextFromResponseWithTimeout(t)
	if err != nil {
		panic(fmt.Sprintf("bad context: %s", err))
	}
	return cs.Context
}

func GetAuthManagerClient() apb.AuthManagerServiceClient {
	managerClient()
	return authManager
}

// compat with 'create', synonym for GetAuthClient()
func GetAuthenticationServiceClient() apb.AuthenticationServiceClient {
	return GetAuthClient()
}

// compat with 'create', synonym for GetAuthClient()
func GetAuthenticationService() apb.AuthenticationServiceClient {
	return GetAuthClient()
}

func GetAuthClient() apb.AuthenticationServiceClient {
	authClient()
	return authServer
}

// create an outbound context for a given user. user must be valid and signed
// this is an expensive call ! (also calls rpcinterceptor)
// this is not privileged (user must be signed)
func ContextForUser(user *apb.User) (context.Context, error) {
	return ContextForUserWithTimeout(user, 0) //default timeout
}
func ContextForUserWithTimeout(user *apb.User, secs uint64) (context.Context, error) {
	if user == nil {
		return nil, fmt.Errorf("Missing user")
	}
	if rpci == nil {
		rpci = rc.NewRPCInterceptorServiceClient(client.Connect("rpcinterceptor.RPCInterceptorService"))
	}
	token := tokens.GetServiceTokenParameter()
	mt := &rc.InMetadata{
		FooBar:       "local",
		ServiceToken: token,
		UserID:       user.ID,
		User:         user,
	}
	mts, err := utils.Marshal(mt)
	if err != nil {
		return nil, err
	}
	cs := &rpc.CallState{
		Started:  time.Now(),
		Debug:    true,
		Metadata: mt,
	}
	newmd := metadata.Pairs(tokens.METANAME, mts)
	ctx := tokens.ContextWithToken()
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	res := metadata.NewOutgoingContext(ctx, newmd)
	cs.Context = ctx
	resp, err := rpci.InterceptRPC(ctx, &rc.InterceptRPCRequest{InMetadata: mt, Service: "local", Method: "local", Source: "go-easyops/auth/context.go"})
	if err != nil {
		return nil, err
	}
	cs.RPCIResponse = resp
	//	cs.PrintContext()
	if resp != nil {
		if mt != nil && mt.User == nil {
			mt.User = resp.CallerUser
			mt.SignedUser = resp.SignedCallerUser
		}
		if mt != nil && mt.Service == nil {
			mt.Service = resp.CallerService
			mt.SignedService = resp.SignedCallerService
		}
		if mt != nil && mt.SignedUser == nil {
			mt.SignedUser = resp.SignedCallerUser
			if mt.User == nil {
				mt.User = common.VerifySignedUser(mt.SignedUser)
			}
		}
		if resp.CallerUser != nil && resp.SignedCallerUser == nil {
			fmt.Printf("[go-easyops] WARNING: authremote.ContextForUser created incomplete context with user, but no signeduser\n")
		}
	}
	// now rebuild metadata again to add to outbound context
	// this must be the final step
	mts, err = utils.Marshal(mt)
	if err != nil {
		return nil, err
	}
	newmd = metadata.Pairs(tokens.METANAME, mts)
	if secs == 0 {
		ctx = tokens.ContextWithToken()
	} else {
		ctx = tokens.ContextWithTokenAndTimeout(secs)
	}
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	res = metadata.NewOutgoingContext(ctx, newmd)
	cs.Context = ctx

	return res, nil

}

// create an outbound context for a given user by id (with current service token)
// this is an expensive call ! (also calls rpcinterceptor)
// it is also privileged
func ContextForUserID(userid string) (context.Context, error) {
	return ContextForUserIDWithTimeout(userid, 0)
}
func ContextForUserIDWithTimeout(userid string, to time.Duration) (context.Context, error) {
	if userid == "" || userid == "0" {
		return nil, fmt.Errorf("Missing userid")
	}
	if rpci == nil {
		rpci = rc.NewRPCInterceptorServiceClient(client.Connect("rpcinterceptor.RPCInterceptorService"))
	}
	token := tokens.GetServiceTokenParameter()
	if token == "" {
		return nil, fmt.Errorf("no service token parameter to generate contextforuser by id")
	}
	mt := &rc.InMetadata{
		FooBar:       "local",
		ServiceToken: token,
		UserID:       userid,
	}
	mts, err := utils.Marshal(mt)
	if err != nil {
		return nil, err
	}
	cs := &rpc.CallState{
		Started:  time.Now(),
		Debug:    true,
		Metadata: mt,
	}
	newmd := metadata.Pairs(tokens.METANAME, mts)
	ctx := tokens.ContextWithToken()
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	res := metadata.NewOutgoingContext(ctx, newmd)
	cs.Context = ctx
	resp, err := rpci.InterceptRPC(ctx, &rc.InterceptRPCRequest{InMetadata: mt, Service: "local", Method: "local", Source: "go-easyops/auth/context.go"})
	if err != nil {
		return nil, err
	}
	cs.RPCIResponse = resp
	//	cs.PrintContext()
	if resp != nil && mt != nil && mt.User == nil {
		mt.User = resp.CallerUser
		mt.SignedUser = resp.SignedCallerUser
	}
	if resp != nil && mt != nil && mt.Service == nil {
		mt.Service = resp.CallerService
		mt.SignedService = resp.SignedCallerService
	}

	// now rebuild metadata again to add to outbound context
	// this must be the final step
	mts, err = utils.Marshal(mt)
	if err != nil {
		return nil, err
	}
	newmd = metadata.Pairs(tokens.METANAME, mts)
	if to == 0 {
		ctx = tokens.ContextWithToken()
	} else {
		ctx = tokens.ContextWithTokenAndTimeout(uint64(to.Seconds()))
	}
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, cs)
	res = metadata.NewOutgoingContext(ctx, newmd)
	cs.Context = ctx

	return res, nil

}
func GetUserByID(ctx context.Context, userid string) (*apb.User, error) {
	if userid == "" {
		return nil, fmt.Errorf("[go-easyops] No userid provided")
	}
	o, err := userbyidcache.Retrieve(userid, func(k string) (interface{}, error) {
		managerClient()
		res, err := authManager.GetUserByID(ctx, &apb.ByIDRequest{UserID: k})
		return res, err
	})
	if err != nil {
		return nil, err
	}
	return o.(*apb.User), nil
}
func GetUserByEmail(ctx context.Context, email string) (*apb.User, error) {
	if email == "" {
		return nil, fmt.Errorf("[go-easyops] No email provided")
	}
	o, err := userbyemailcache.Retrieve(email, func(k string) (interface{}, error) {
		managerClient()
		res, err := authManager.GetUserByEmail(ctx, &apb.ByEmailRequest{Email: k})
		return res, err
	})
	if err != nil {
		return nil, err
	}
	return o.(*apb.User), nil
}
func WhoAmI() *apb.User {
	tok := tokens.GetUserTokenParameter()
	return GetByToken(context.Background(), tok)
}
func GetByToken(ctx context.Context, token string) *apb.User {
	if token == "" {
		return nil
	}
	authClient()
	ar, err := authServer.GetByToken(ctx, &apb.AuthenticateTokenRequest{Token: token})
	if err != nil {
		return nil
	}
	if !ar.Valid {
		return nil
	}
	if !ar.User.Active {
		return nil
	}
	return ar.User
}
func SignedGetByToken(ctx context.Context, token string) *apb.SignedUser {
	if token == "" {
		return nil
	}
	su, err := userbytokencache.Retrieve(token, func(k string) (interface{}, error) {
		authClient()
		ar, err := authServer.SignedGetByToken(ctx, &apb.AuthenticateTokenRequest{Token: k})
		if err != nil {
			return nil, err
		}
		if !ar.Valid {
			return nil, fmt.Errorf("user not valid")
		}
		u := common.VerifySignedUser(ar.User)
		if !u.Active {
			return nil, fmt.Errorf("user not active")
		}
		return ar.User, nil
	})
	if err != nil {
		return nil
	}
	if su == nil {
		return nil
	}
	return su.(*apb.SignedUser)
}

func authClient() {
	if authServer == nil {
		authServerLock.Lock()
		defer authServerLock.Unlock()
		if authServer != nil {
			return
		}
		authServer = apb.NewAuthenticationServiceClient(client.Connect("auth.AuthenticationService"))
	}
}
func managerClient() {
	if authManager == nil {
		authManagerLock.Lock()
		defer authManagerLock.Unlock()
		if authManager != nil {
			return
		}
		authManager = apb.NewAuthManagerServiceClient(client.Connect("auth.AuthManagerService"))
	}
}

/*
This package provides access to user information which require network I/O, for example lookup of users by email.

It also provides some wrappers to create a new context. That is for historic reasons. Developers should use and port code to use the ctx package instead.
*/
package authremote

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/tokens"
	"sync"
	"time"
)

var (
	userbyemailcache = cache.NewResolvingCache("userbyemail", time.Duration(60)*time.Second, 9999)
	userbytokencache = cache.NewResolvingCache("userbytoken", time.Duration(60)*time.Second, 9999)
	authServer       apb.AuthenticationServiceClient
	authServerLock   sync.Mutex
	authManager      apb.AuthManagerServiceClient
	authManagerLock  sync.Mutex
	contextRetrieved = false
	lastUser         *apb.SignedUser
	lastService      *apb.SignedUser
)

func init() {
	common.AddRegistryChangeReceiver(registry_changed)
}
func registry_changed() {
	authManager = nil
	authServer = nil
}
func Context() context.Context {
	client.GetSignatureFromAuth()
	return ContextWithTimeout(time.Duration(10) * time.Second)
}

/*
create a new context with routing tags (routing criteria to route to specific instances of a service)
if fallback is true, fallback to any service without tags if none is found (default was false)
*/
func NewContextWithRouting(kv map[string]string, fallback bool) context.Context {
	return DerivedContextWithRouting(Context(), kv, fallback)
}

/*
derive  a context with routing tags (routing criteria to route to specific instances of a service)
if fallback is true, fallback to any service without tags if none is found (default was false)
*/
func DerivedContextWithRouting(cv context.Context, kv map[string]string, fallback bool) context.Context {
	if cv == nil {
		panic("cannot derive context from nil context")
	}
	cri := &ge.CTXRoutingTags{Tags: kv, FallbackToPlain: fallback}

	if cmdline.ContextWithBuilder() {
		_, s := GetLocalUsers()
		if s == nil {
			s = auth.GetSignedService(cv)
		}
		cb := ctx.NewContextBuilder()
		cb.WithUser(auth.GetSignedUser(cv))
		cb.WithCreatorService(s)
		cb.WithCallingService(s)
		cb.WithRoutingTags(cri)
		//	cb.WithTimeout(t)
		cb.WithParentContext(cv)
		nctx := cb.ContextWithAutoCancel()
		if auth.GetSignedService(nctx) == nil && s != nil {
			fmt.Printf("[go-easyops] context: %s\n", ctx.Context2String(nctx))
			fmt.Printf("[go-easyops] Localstate: %#v\n", ctx.GetLocalState(nctx))
			fmt.Printf("[go-easyops] WARNING derived context (v=%d) includes no service, but should\n", cmdline.GetContextBuilderVersion())
			//return nil
		}
		if nctx == nil {
			panic("no context")
		}
		return nctx

	}
	panic("deprecated codepath")
}

/*
get a context with routing tags, specified by proto
*/
func NewContextWithRoutingTags(rt *ge.CTXRoutingTags) context.Context {
	return ContextWithTimeoutAndTags(time.Duration(10)*time.Second, rt)
}

/*
	this context gives a context with a full userobject

todo so it _has_ to call external servers to get a signed userobject.
if started_by_autodeployer will use getContext()
else if environment variable with context, will use auth.Context() (with variable)
else create context by asking auth service for a signed user object
*/
func ContextWithTimeout(t time.Duration) context.Context {
	return ContextWithTimeoutAndTags(t, nil)
}

// get the user and service we are running as. Do not cache this result! (on boot the result may change once auth comes available)
func GetLocalUsers() (*apb.SignedUser, *apb.SignedUser) {
	client.GetSignatureFromAuth()
	if cmdline.DebugAuth() {
		fmt.Printf("[go-easyops] debugauth, contextretrieved=%v, localuser=%s,gotsig=%v\n", contextRetrieved, auth.SignedDescription(lastUser), client.GotSig())
	}
	if !client.GotSig() {
		if cmdline.DebugAuth() {
			fmt.Printf("[go-easyops] debugauth no local users, we do not yet have a signature\n")
		}
		return nil, nil
	}
	if !contextRetrieved {
		utok := tokens.GetUserTokenParameter()
		//		fmt.Printf("utok: \"%s\"\n", utok)
		lastUser = SignedGetByToken(context_background(), utok)
		lastService = SignedGetByToken(context_background(), tokens.GetServiceTokenParameter())
		lu := common.VerifySignedUser(lastUser)
		if lastUser != nil && lu == nil {
			fmt.Printf("[go-easyops] Warning - local user signature invalid\n")
			return nil, nil
		}
		if lu != nil {
			if lu.ServiceAccount {
				fmt.Printf("[go-easyops] Error - local user resolved to a service account\n")
				panic("invalid user configuration")
			}
		}
		if lastService != nil && common.VerifySignedUser(lastService) == nil {
			fmt.Printf("[go-easyops] Warning - local service signature invalid\n")
			return nil, nil
		}
		contextRetrieved = true
	}
	return lastUser, lastService
}

/*
create a new context with routing tags. This is an EXPERIMENTAL API and very likely to change in future
*/
func ContextWithTimeoutAndTags(t time.Duration, rt *ge.CTXRoutingTags) context.Context {
	if cmdline.IsStandalone() {
		return standalone_ContextWithTimeoutAndTags(t, rt)
	}
	sctx := cmdline.GetEnvContext()
	if sctx != "" {
		if ctx.IsSerialisedByBuilder([]byte(sctx)) {
			ctx, err := ctx.DeserialiseContextWithTimeout(t, []byte(sctx))
			if err != nil {
				fmt.Printf("[go-easyops] weird context GE_CTX (%s)\n", err)
			} else {
				return ctx
			}

		}
		//				fmt.Printf("Recreating context from environment variable GE_CTX\n")
		res, err := auth.RecreateContextWithTimeout(t, []byte(sctx))
		if err == nil {
			return res
		} else {
			fmt.Printf("[go-easyops] invalid context in environment variable GE_CTX\n")
		}
	}
	if cmdline.ContextWithBuilder() {
		u, s := GetLocalUsers()
		cb := ctx.NewContextBuilder()
		cb.WithUser(u)
		cb.WithCreatorService(s)
		cb.WithCallingService(s)
		cb.WithRoutingTags(rt) //rpc.Tags_rpc_to_ge(rt))
		cb.WithTimeout(t)
		return cb.ContextWithAutoCancel()
	}

	panic("[go-easyops] DEPRECATED CONTEXT creation!\n")
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
// this is an expensive call
// this is not privileged (user must be signed)
func ContextForUser(user *apb.User) (context.Context, error) {
	return ContextForUserWithTimeout(user, 0) //default timeout
}
func ContextForUserWithTimeout(user *apb.User, secs uint64) (context.Context, error) {
	if user == nil {
		return nil, fmt.Errorf("Missing user")
	}

	if cmdline.ContextWithBuilder() {
		su, err := GetSignedUserByID(Context(), user.ID)
		if err != nil {
			return nil, err
		}
		cb := ctx.NewContextBuilder()
		cb.WithTimeout(time.Duration(secs) * time.Second)
		cb.WithUser(su)
		_, svc := GetLocalUsers()
		cb.WithCreatorService(svc)
		cb.WithCallingService(svc)
		return cb.ContextWithAutoCancel(), nil
	}
	panic("obsolete codepath")
}

// create an outbound context for a given user by id (with current service token)
// this is an expensive call
// it is also privileged
func ContextForUserID(userid string) (context.Context, error) {
	return ContextForUserIDWithTimeout(userid, 0)
}
func ContextForUserIDWithTimeout(userid string, to time.Duration) (context.Context, error) {
	if userid == "" || userid == "0" {
		return nil, fmt.Errorf("Missing userid")
	}
	if cmdline.ContextWithBuilder() {
		su, err := GetSignedUserByID(Context(), userid)
		if err != nil {
			return nil, err
		}
		cb := ctx.NewContextBuilder()
		cb.WithTimeout(to)
		cb.WithUser(su)
		_, svc := GetLocalUsers()
		cb.WithCreatorService(svc)
		cb.WithCallingService(svc)
		return cb.ContextWithAutoCancel(), nil
	}
	panic("obsolete codepath")

}
func GetUserByID(ctx context.Context, userid string) (*apb.User, error) {
	if userid == "" {
		return nil, fmt.Errorf("[go-easyops] No userid provided")
	}
	return usercache_GetUserByID(ctx, userid)
}

func GetSignedUserByID(ctx context.Context, userid string) (*apb.SignedUser, error) {
	if userid == "" {
		return nil, fmt.Errorf("[go-easyops] No userid provided")
	}
	return usercache_GetSignedUserByID(ctx, userid)
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
	return GetByToken(context_background(), tok)
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
		//		utils.PrintStack("[go-easyops] attempt to get user for empty token\n")
		return nil
	}
	su, err := userbytokencache.Retrieve(token, func(k string) (interface{}, error) {
		authClient()
		if cmdline.DebugAuth() {
			fmt.Printf("[go-easyops] getting user for token \"%s\"...\n", k[:5])
		}
		ar, err := authServer.SignedGetByToken(ctx, &apb.AuthenticateTokenRequest{Token: k})
		if err != nil {
			if cmdline.DebugAuth() {
				fmt.Printf("[go-easyops] getting user for token \"%s\" failed: %s\n", k, err)
			}
			return nil, err
		}
		if !ar.Valid {
			if cmdline.DebugAuth() {
				fmt.Printf("[go-easyops] getting user for token \"%s\" invalid", k[:5])
			}
			return nil, fmt.Errorf("user not valid")
		}
		u := common.VerifySignedUser(ar.User)
		if !u.Active {
			if cmdline.DebugAuth() {
				fmt.Printf("[go-easyops] getting user for token \"%s\" inactive", k[:5])
			}
			return nil, fmt.Errorf("user not active")
		}
		if cmdline.DebugAuth() {
			fmt.Printf("[go-easyops] getting user for token \"%s\" resulted in \"%s\"", k[:5], ar.User)
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

func context_background() context.Context {
	cb := ctx.NewContextBuilder()
	return cb.ContextWithAutoCancel()
}

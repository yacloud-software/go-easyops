package server

import (
	"context"
	"fmt"
	au "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/registry"
)

type ServerDef interface {

	//if something needs to be done to errors before they propagate up the stack, then this hook can be used to do so
	SetErrorHandler(e func(ctx context.Context, fn string, err error))
	/*
	 set to true if this server does NOT require authentication (default: it does need authentication).
	 This should normally not be necessary. Normally, a service needs to be called with EITHER a service account OR a user account OR both. There are very special circumstances where this is not possible, for example, the registry and the auth service cannot be called with a service or user account, because in order to get one, the service needs to lookup and call the auth service. Thus registry and auth both expose their RPCs as "NoAuth". In normal circumstances this is never necessary.
	*/
	SetNoAuth()
	// the tcp port to listen on
	SetPort(port int)
	// register the implementation of the gRPC service
	SetRegister(r Register)
	DontRegister() // if this service should not register with the registry initially
	// assume the service is directly accessible on a public ip. this disables functionalitity normally filtered out by proxies, such as /internal/ helpers and reflection. Normally not needed. Typically h2gproxy proxies requests.
	SetPublic()
	/*
	   set a callback that is called AFTER grpc server started successfully
	*/
	SetOnStartupCallback(f func())
	AddTag(key, value string) // add a routing tag to a serverdef

}

// no longer exported - please use NewServerDef instead
type serverDef struct {
	callback    func()   // called if/when server started up successfully
	port        int      // the port the GRPC server should listen on
	Certificate []byte   // do not override the default. Exposed due to an implementation limitation
	Key         []byte   // do not override the default. Exposed due to an implementation limitation
	CA          []byte   // do not override the default. Exposed due to an implementation limitation
	register    Register // do not override the default. Exposed due to an implementation limitation
	/*
	 set to true if this server does NOT require authentication (default: it does need authentication).
	 This should normally not be necessary. Normally, a service needs to be called with EITHER a service account OR a user account OR both. There are very special circumstances where this is not possible, for example, the registry and the auth service cannot be called with a service or user account, because in order to get one, the service needs to lookup and call the auth service. Thus registry and auth both expose their RPCs as "NoAuth". In normal circumstances this is never necessary.
	*/
	noAuth bool
	// set to false if this service should not register with the registry initially
	registerService bool
	name            string
	types           []pb.Apitype
	registered_id   string
	deployPath      string // do not override the default. Exposed due to an implementation limitation
	serviceID       uint64
	asUser          *au.SignedUser // if we're running as a user rather than a server this is the account
	tags            map[string]string
	/*
	   if something needs to be done to errors before they propagate up the stack, then this hook can be used to do so
	*/
	errorHandler    func(ctx context.Context, function_name string, err error)
	local_service   *au.SignedUser // the local service account
	service_user_id string         // the serviceaccount userid
	public          bool
	port_set        bool
}

func (s *serverDef) SetErrorHandler(e func(ctx context.Context, fn string, err error)) {
	s.errorHandler = e
}
func (s *serverDef) SetNoAuth() {
	s.noAuth = true
}
func (s *serverDef) SetPort(port int) {
	s.port = port
	s.port_set = true
}
func (s *serverDef) SetRegister(r Register) {
	s.register = r
}
func (s *serverDef) DontRegister() {
	s.registerService = false
}
func (s *serverDef) SetPublic() {
	s.public = true
}

/*
set a callback that is called AFTER grpc server started successfully
*/
func (s *serverDef) SetOnStartupCallback(f func()) {
	s.callback = f
}

// add a routing tag to a serverdef
func (s *serverDef) AddTag(key, value string) {
	s.tags[key] = value
}
func (s *serverDef) toString() string {
	return fmt.Sprintf("Port #%d: %s (%v)", s.port, s.name, s.types)
}

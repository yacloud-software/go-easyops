package server

import (
	"context"
	au "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/registry"
)

// no longer exported - please use NewServerDef instead
type serverDef struct {
	callback    func()   // called if/when server started up successfully
	Port        int      // the port the GRPC server should listen on
	Certificate []byte   // do not override the default. Exposed due to an implementation limitation
	Key         []byte   // do not override the default. Exposed due to an implementation limitation
	CA          []byte   // do not override the default. Exposed due to an implementation limitation
	Register    Register // do not override the default. Exposed due to an implementation limitation
	/*
	 set to true if this server does NOT require authentication (default: it does need authentication).
	 This should normally not be necessary. Normally, a service needs to be called with EITHER a service account OR a user account OR both. There are very special circumstances where this is not possible, for example, the registry and the auth service cannot be called with a service or user account, because in order to get one, the service needs to lookup and call the auth service. Thus registry and auth both expose their RPCs as "NoAuth". In normal circumstances this is never necessary.
	*/
	NoAuth bool
	// set to false if this service should not register with the registry initially
	RegisterService bool
	name            string
	types           []pb.Apitype
	registered_id   string
	DeployPath      string // do not override the default. Exposed due to an implementation limitation
	serviceID       uint64
	asUser          *au.SignedUser // if we're running as a user rather than a server this is the account
	tags            map[string]string
	/*
	   if something needs to be done to errors before they propagate up the stack, then this hook can be used to do so
	*/
	ErrorHandler    func(ctx context.Context, function_name string, err error)
	local_service   *au.SignedUser // the local service account
	service_user_id string         // the serviceaccount userid
	public          bool
}

func (s *serverDef) SetPort(port int) {
	s.Port = port
}
func (s *serverDef) SetRegister(r Register) {
	s.Register = r
}
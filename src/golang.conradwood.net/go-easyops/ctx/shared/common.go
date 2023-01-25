package shared

import (
	"context"
	"golang.conradwood.net/apis/auth"
	ge "golang.conradwood.net/apis/goeasyops"
	"time"
)

const (
	LOCALSTATENAME = "goeasysops_localstate"
)

// the local state, this is not transmitted across grpc boundaries.
type LocalState interface {
	CreatorService() *auth.SignedUser
	CallingService() *auth.SignedUser
	Debug() bool
	Trace() bool
	User() *auth.SignedUser
	Session() *auth.SignedSession
	RequestID() string
	RoutingTags() *ge.CTXRoutingTags
}

type ContextBuilder interface {
	/*
		This function parses metadata found in an inbound context and, if successful, returns an "outbound" context with localstate.
		the bool return parameter indicates if it was successful(true) or not(false).
		Note that it requires the LOCAL service, because the calling service is modified and passed to the next service
	*/
	Inbound2Outbound(ctx context.Context, svc *auth.SignedUser) (context.Context, bool)
	/*
		return the context from this builder based on the options and WithXXX functions
	*/
	Context() (context.Context, context.CancelFunc)

	// like Context(), but automatically call the CancelFunc after timeout
	ContextWithAutoCancel() context.Context

	/*
	   add a user to context
	*/
	WithUser(user *auth.SignedUser)

	/*
	   add a creator service to context
	*/
	WithCreatorService(user *auth.SignedUser)

	/*
	   add a calling service (e.g. "me") to context
	*/
	WithCallingService(user *auth.SignedUser)

	/*
	   add a session to the context
	*/
	WithSession(user *auth.SignedSession)

	// mark context as with debug
	WithDebug()

	// mark context as with trace
	WithTrace()
	// add routing tags
	WithRoutingTags(*ge.CTXRoutingTags)
	//set the requestid
	WithRequestID(reqid string)
	// set a timeout for this context
	WithTimeout(time.Duration)
	// set a parent context for cancellation propagation (does not transfer metadata to the new context!)
	WithParentContext(context context.Context)
}

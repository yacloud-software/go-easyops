package errors

// package errors
// grpc Servers should *only* return errors created by this package.
// so instead of fmt.Errorf() or status.Error use
// errors.Error() (this package)
import (
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// mapping as per https://cloud.google.com/apis/design/errors
	grpcToHTTPMap = map[codes.Code]*HTTPError{
		codes.OK:                 {200, "ok", "", ""},
		codes.Unknown:            {500, "unknown method", "", ""},
		codes.InvalidArgument:    {400, "invalid argument", "", ""},
		codes.DeadlineExceeded:   {504, "deadline exceeded", "", ""},
		codes.NotFound:           {404, "not found", "", ""},
		codes.AlreadyExists:      {409, "resource already exists", "", ""},
		codes.PermissionDenied:   {403, "insufficient permission", "", ""},
		codes.ResourceExhausted:  {429, "out of resource quota", "", ""},
		codes.FailedPrecondition: {400, "not possible in current system state", "", ""},
		codes.Aborted:            {409, "concurrency conflict", "", ""},
		codes.OutOfRange:         {400, "invalid range specified", "", ""},
		codes.Unimplemented:      {501, "method not implemented", "", ""},
		codes.Internal:           {500, "internal server error", "", ""},
		codes.Unavailable:        {503, "service unavailable", "", ""},
		codes.DataLoss:           {500, "internal server error", "", ""},
		codes.Unauthenticated:    {401, "missing, invalid, or expired authentication", "", ""},
	}
)

type HTTPError struct {
	ErrorCode           int
	ErrorString         string
	ExtendedErrorString string
	ErrorMessage        string
}

// error if context is not root user or one of the services listed
func NeedServiceOrRoot(ctx context.Context, serviceids []string) error {
	err := NeedsRoot(ctx)
	if err == nil {
		return nil
	}
	u := auth.GetUser(ctx)
	svc := auth.GetService(ctx)
	if svc == nil {
		if u == nil {
			return Unauthenticated(ctx, "please log in")
		} else {
			return AccessDenied(ctx, "not allowed")
		}
	}
	for _, svid := range serviceids {
		if svid == svc.ID {
			return nil
		}
	}
	if u == nil {
		return Unauthenticated(ctx, "please log in")
	} else {
		return AccessDenied(ctx, "not allowed")
	}

}

// function call requires "root" privileges. returns error if user is non-root
func NeedsRoot(ctx context.Context) error {
	u := auth.CurrentUserString(ctx)
	if auth.IsRootUser(auth.GetUser(ctx)) {
		return nil
	}
	return Error(ctx, codes.PermissionDenied, "access denied", "this function requires root privileges (which %s does not have)", u)
}

func NotImplemented(ctx context.Context, method string) error {
	return Error(ctx, codes.Unimplemented, "functionality is not implemented", "function %s not implemented", method)
}
func Unavailable(ctx context.Context, method string) error {
	return Error(ctx, codes.Unavailable, "currently unavailable", "this RPC or data is currently unavailable (%s)", method)
}
func FailedPrecondition(ctx context.Context, logmessage string, a ...interface{}) error {
	return Error(ctx, codes.FailedPrecondition, "state mismatch", logmessage, a...)
}
func AccessDenied(ctx context.Context, logmessage string, a ...interface{}) error {
	return Error(ctx, codes.PermissionDenied, "access denied", logmessage, a...)
}
func NotFound(ctx context.Context, logmessage string, a ...interface{}) error {
	return Error(ctx, codes.NotFound, "not found", logmessage, a...)
}
func Unauthenticated(ctx context.Context, logmessage string, a ...interface{}) error {
	return Error(ctx, codes.Unauthenticated, "access denied", logmessage, a...)
}

// shortcut: we write this so often: user submitted args that aren't valid
func InvalidArgs(ctx context.Context, publicmessage string, logmessage string, a ...interface{}) error {
	return Error(ctx, codes.InvalidArgument, publicmessage, logmessage, a...)
	//	return Error(ctx, codes.FailedPrecondition, publicmessage, logmessage, a...)
}

// include caller/callee information in logmessage
func stdText(ctx context.Context) string {
	user := auth.CurrentUserString(ctx)
	svc := auth.GetService(ctx)
	cs := rpc.CallStateFromContext(ctx)
	caller := "nil"
	callee := "nil"
	if cs != nil {
		callee = cs.ServiceName + "." + cs.MethodName
		if cs.RPCIResponse != nil && svc != nil {
			caller = fmt.Sprintf("%s(%s).method-%d", svc.ID, svc.Email, cs.RPCIResponse.CallerMethodID)
		}
	}

	res := fmt.Sprintf("[%s called %s as user=%s", caller, callee, user)
	if svc == nil {
		res = res + ", noservice"
	} else {
		res = res + ", service=" + svc.ID + " (" + svc.Email + ")"
	}
	res = res + "]"
	return res
}

// really returns a status.Status
func Error(ctx context.Context, code codes.Code, publicmessage string, logmessage string, a ...interface{}) error {
	var err error
	logmessage = stdText(ctx) + logmessage
	log := fmt.Sprintf(logmessage, a...)
	st := status.New(code, publicmessage)
	// encapsulate "status" with logmessage
	add := &common.Status{ErrorCode: int32(code), ErrorDescription: log}
	st, err = st.WithDetails(add)
	if err != nil {
		// this is bad. we can't create an error to reflect the error
		// in case of double-faults there isn't really any other option than to log and exit
		panic(fmt.Sprintf("Double fault, error in error handler whilst creating error for code=%d, publicmessage=%s, logmessage=%s: %s", code, publicmessage, log, err))
	}
	return st.Err()
}
func ToHTTPCode(err error) *HTTPError {
	st := status.Convert(err)
	code := st.Code()
	he, f := grpcToHTTPMap[code]
	if !f {
		he = &HTTPError{ErrorCode: 500,
			ErrorString:         "Unspecified error",
			ExtendedErrorString: fmt.Sprintf("GRPC Error %d", code),
			ErrorMessage:        "Unspecified error",
		}
	}
	return he

}

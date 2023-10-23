package server

import (
	"context"
	"fmt"
	fw "golang.conradwood.net/apis/framework"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/errors"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/metadata"
	"flag"
	"google.golang.org/grpc/status"
	"time"
)

var (
	print_errs = flag.Bool("ge_grpc_print_errors", false, "if true print grpc errors before they propagate to the caller")
)

/*******************************************************************************************
* gRPC calls this interceptor for each call. Be fast and reliable
*******************************************************************************************/

// we authenticate a client here
func (sd *serverDef) UnaryAuthInterceptor(in_ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	pp.ServerRpcEntered()
	defer pp.ServerRpcDone()
	cs := &rpccall{
		ServiceName: ServiceNameFromUnaryInfo(info),
		MethodName:  MethodNameFromUnaryInfo(info),
		Started:     time.Now(),
	}
	if *debug_rpc_serve {
		fmt.Printf("[go-easyops] Debug-rpc called unary rpc \"%s\"\n", info.FullMethod)
	}

	started := time.Now()

	var outbound_ctx context.Context
	var err error

	ctx_build_by := 0
	// we try both types of context parsing (since we can be called by either, old or new service)
	if cmdline.ContextWithBuilder() {
		ctx_build_by = 1
		outbound_ctx, _, err = sd.V1inbound2outbound(in_ctx, cs)
		if err != nil {
			return nil, err
		}
	} else {
		panic("obsolete codepath")
	}
	if err != nil {
		return nil, err
	}
	if *debug_rpc_serve {
		fmt.Printf("[go-easyops] context created through path %d\n", ctx_build_by)
	}
	//fmt.Printf("LS: %#v\n", ls)
	//fmt.Printf("Method: \"%s\"\n", method)
	stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      cs.MethodName,
		"servicename": cs.ServiceName,
	}).Inc()
	defer stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      cs.MethodName,
		"servicename": cs.ServiceName,
	}).Dec()

	grpc_server_requests.With(prometheus.Labels{
		"method":      cs.MethodName,
		"servicename": cs.ServiceName,
	}).Inc()

	print_inbound_debug(cs, outbound_ctx)
	/*************** now call the rpc implementation *****************/
	i, err := handler(outbound_ctx, req)
	if i == nil && err == nil {
		fmt.Printf("[go-easyops] BUG: \"%s.%s\" returned no proto and no error\n", cs.ServiceName, cs.MethodName)
	}
	if *debug_rpc_serve {
		//		fmt.Printf("[go-easyops: result: %v %v\n", i, err)
		fmt.Printf("[go-easyops] Debug-rpc Request: \"%s.%s\" timing: %0.2fs\n", cs.ServiceName, cs.MethodName, time.Since(started).Seconds())
	}
	if err == nil {
		grpc_server_req_durations.WithLabelValues(cs.ServiceName, cs.MethodName).Observe(time.Since(cs.Started).Seconds())
		return i, nil
	}
	// it failed!
	dur := time.Since(cs.Started).Seconds()
	if dur > 5 { // >5 seconds processing time? warn
		fmt.Printf("[go-easyops] Debug-rpc Request: \"%s.%s\" took rather long: %0.2fs (and failed: %s)\n", cs.ServiceName, cs.MethodName, dur, err)
	}
	if *debug_rpc_serve || *print_errs {
		fmt.Printf("[go-easyops] Debug-rpc Request: \"%s.%s\" (called from %s) failed: %s\n", cs.ServiceName, cs.MethodName, auth.UserIDString(auth.GetService(outbound_ctx)), err)
	}
	incFailure(cs.ServiceName, cs.MethodName, err)
	//stdMetrics.grpc_failed_requests.With(prometheus.Labels{"method": method, "servicename": def.name}).Inc()

	// get status from error
	st := status.Convert(err)
	fm := &fw.CallTrace{
		Message: fmt.Sprintf("[go-easyops] GRPC error in method %s.%s()", cs.ServiceName, cs.MethodName),
		Method:  cs.MethodName,
		Service: cs.ServiceName,
	}
	st = AddStatusDetail(st, fm)
	gerr := &ge.GRPCError{
		LogMessage:  fmt.Sprintf("%s", err),
		MethodName:  cs.MethodName,
		ServiceName: cs.ServiceName,
	}
	calling_svc := auth.GetService(outbound_ctx)
	if calling_svc != nil {
		gerr.CallingServiceID = calling_svc.ID
		gerr.CallingServiceEmail = calling_svc.Email
	}
	st = AddErrorDetail(st, gerr)
	re := st.Err()
	sd.logError(outbound_ctx, cs, re)
	eh := sd.errorHandler
	if eh != nil {
		eh(outbound_ctx, cs.MethodName, err)
	}
	return i, st.Err()
}

func (sd *serverDef) V1inbound2outbound(in_ctx context.Context, rc *rpccall) (context.Context, shared.LocalState, error) {
	if sd.local_service == nil {
		if *debug_rpc_serve {
			fmt.Printf("[go-easyops] WARNING, in server.unary_interceptor, we are converting inbound2outbound without a local service account\n")
		}
	}
	u := auth.GetUser(in_ctx)
	if u != nil && u.ServiceAccount {
		return nil, nil, errors.Unauthenticated(in_ctx, "user %s is a serviceaccount", u.ID)
	}
	octx := ctx.Inbound2Outbound(in_ctx, sd.local_service)
	ls := ctx.GetLocalState(octx)
	err := sd.checkAccess(octx, rc)
	if err != nil {
		if *debug_rpc_serve {
			fmt.Printf("[go-easyops] checkaccess error: %s (peer=%s)\n", err, peerFromContext(octx))
			fmt.Printf("[go-easyops] Context: %#v\n", ctx.Context2String(octx))
		}
		return nil, nil, err
	}
	if ls == nil {
		panic("no localstate in newly converted inbound context")
	}
	return octx, ls, nil
}

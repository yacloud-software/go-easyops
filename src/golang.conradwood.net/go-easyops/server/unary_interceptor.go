package server

import (
	"context"
	"fmt"
	fw "golang.conradwood.net/apis/framework"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	//	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

/*******************************************************************************************
* gRPC calls this interceptor for each call. Be fast and reliable
*******************************************************************************************/
// we authenticate a client here
func (sd *serverDef) UnaryAuthInterceptor(in_ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	pp.ServerRpcEntered()
	defer pp.ServerRpcDone()
	started := time.Now()
	cs := &rpc.CallState{
		Started:     time.Now(),
		ServiceName: ServiceNameFromUnaryInfo(info),
		MethodName:  MethodNameFromUnaryInfo(info),
		Context:     in_ctx,
		MyServiceID: sd.serviceID,
	}
	ctx := context.WithValue(in_ctx, rpc.LOCALCONTEXTNAME, cs)
	cs.Context = ctx
	def := getServerDefByName(cs.ServiceName)

	// this can happen if we're looking for a different service than who we are.
	// it's really bad actually (probably a bug);)

	if def == nil {
		s := fmt.Sprintf("[go-easyops] Service not registered! %s", cs.ServiceName)
		fmt.Println(s)
		return nil, errors.Error(cs.Context, codes.Unimplemented, "service unavailable", "service %s is not known here", cs.ServiceName)
	}

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

	// if we're a "noauth" service we MUST NOT call rpcinterceptor (due to the risk of loops)
	if !def.NoAuth {
		err := Authenticate(cs)
		if err != nil {
			fancyPrintf("Debug-rpc Request: Authenticate() failed for \"%s/%s\" => rejected: %s\n", cs.MethodName, cs.ServiceName, err)
			return nil, err
		}
		if cs.RPCIResponse.Reject {
			return nil, errors.Error(cs.Context, codes.PermissionDenied, "access denied", "Access denied to %s for user %s", cs.TargetString(), cs.CallerString())
		}
	}
	if cs.Metadata != nil {
		cs.Metadata.FooBar = "go-easyops"
		if cs.Metadata.RoutingTags != nil && cs.Metadata.RoutingTags.Propagate == false {
			cs.Metadata.RoutingTags = nil
		}
	}

	// this rebuilds the metadata string from cs.Metadata
	cs.UpdateContextFromResponse()
	cs.DebugPrintContext()
	if *debug_rpc_serve {
		user := auth.GetUser(cs.Context)
		svc := auth.GetService(cs.Context)
		fancyPrintf("Debug-rpc Request: \"%s/%s\" invoked by user \"%s\" from service \"%s\"\n", cs.ServiceName, cs.MethodName, auth.UserIDString(user), auth.UserIDString(svc))
	}

	/*************** now call the rpc implementation *****************/
	i, err := handler(cs.Context, req)
	if i == nil && err == nil {
		fmt.Printf("[go-easyops] BUG: \"%s.%s\" returned no proto and no error\n", cs.ServiceName, cs.MethodName)
	}
	if *debug_rpc_serve {
		//		fmt.Printf("[go-easyops: result: %v %v\n", i, err)
		fmt.Printf("[go-easyops] Debug-rpc Request: \"%s.%s\" timing: %v\n", cs.ServiceName, cs.MethodName, time.Since(started))
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
	if *debug_rpc_serve {
		fmt.Printf("[go-easyops] Debug-rpc Request: \"%s.%s\" failed: %s\n", cs.ServiceName, cs.MethodName, err)
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
	re := st.Err()
	sd.logError(cs, re)
	eh := sd.ErrorHandler
	if eh != nil {
		eh(cs.Context, cs.MethodName, err)
	}
	return i, st.Err()
}

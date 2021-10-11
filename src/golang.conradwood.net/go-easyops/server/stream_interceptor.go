package server

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	fw "golang.conradwood.net/apis/framework"
	"golang.conradwood.net/go-easyops/errors"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func (sd *serverDef) StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var err error
	pp.ServerRpcEntered()
	defer pp.ServerRpcDone()

	name := ServiceNameFromStreamInfo(info)
	if *debug_rpc_serve {
		fmt.Printf("[go-easyops] debug: called streaming service %s\n", name)
	}
	method := MethodNameFromStreamInfo(info)
	stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": name,
	}).Inc()
	defer stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": name,
	}).Dec()

	//fmt.Printf("Method: \"%s\"\n", method)
	if isInternalService(name) {
		if *debug_rpc_serve {
			fmt.Printf("Invoking internal service stream handler\n")
		}
		res := handler(srv, stream)
		if *debug_rpc_serve {
			fmt.Printf("internal service stream handler returned: %s\n", res)
		}
		return res
	}

	cs := &rpc.CallState{
		Started:     time.Now(),
		ServiceName: ServiceNameFromStreamInfo(info),
		MethodName:  MethodNameFromStreamInfo(info),
		Context:     stream.Context(),
		MyServiceID: sd.serviceID,
	}
	ctx := context.WithValue(stream.Context(), rpc.LOCALCONTEXTNAME, cs)
	cs.Context = ctx

	def := getServerDefByName(name)
	if def == nil {
		s := fmt.Sprintf("[go-easyops] Service not registered! %s", name)
		fmt.Println(s)
		return errors.Error(cs.Context, codes.Unimplemented, "service unavailable", "service %s is not known here", cs.ServiceName)
	}

	grpc_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": def.name,
	}).Inc()
	// if we're a "noauth" service we MUST NOT call rpcinterceptor (due to the risk of loops)
	if !def.NoAuth {
		err := Authenticate(cs)
		if err != nil {
			return err
		}
		if cs.RPCIResponse.Reject {
			return errors.AccessDenied(cs.Context, "Access denied to %s for user %s", cs.TargetString(), cs.CallerString())
		}
	}
	if cs.Metadata != nil {
		cs.Metadata.FooBar = "nonmoo"
	}
	cs.UpdateContextFromResponse()
	cs.DebugPrintContext()
	nstream := newServerStream(stream, cs.Context)
	err = handler(srv, nstream)
	if err == nil {
		return nil
	}
	if *debug_rpc_serve {
		fmt.Printf("[go-easyops] Call %s.%s failed: %s\n", def.name, method, err)
	}
	incFailure(def.name, method, err)

	// get status from error
	st := status.Convert(err)
	fm := fw.CallTrace{
		Message: fmt.Sprintf("[go-easyops] GRPC error in method %s.%s()", def.name, method),
		Method:  method,
		Service: def.name,
	}

	// add details
	st, errx := st.WithDetails(&fm)

	// if adding details failed, just return the undecorated error message
	if errx != nil {
		sd.logError(cs, err)
		return err
	}

	re := st.Err()
	sd.logError(cs, re)
	return re
}
func MethodNameFromStreamInfo(info *grpc.StreamServerInfo) string {
	full := info.FullMethod
	if full[0] == '/' {
		full = full[1:]
	}
	ns := strings.SplitN(full, "/", 2)
	if len(ns) < 2 {
		return ""
	}
	res := ns[1]
	if res[0] == '/' {
		res = res[1:]
	}
	return ns[1]
}
func ServiceNameFromStreamInfo(info *grpc.StreamServerInfo) string {
	full := info.FullMethod
	if full[0] == '/' {
		full = full[1:]
	}
	ns := strings.SplitN(full, "/", 2)
	return ns[0]
}

type customServerStream struct {
	stream grpc.ServerStream
	ctx    context.Context
}

func newServerStream(in grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	res := &customServerStream{stream: in, ctx: ctx}
	return res
}

func (c *customServerStream) SetHeader(m metadata.MD) error {
	return c.stream.SetHeader(m)
}
func (c *customServerStream) SendHeader(m metadata.MD) error {
	return c.stream.SendHeader(m)
}
func (c *customServerStream) SetTrailer(m metadata.MD) {
	c.stream.SetTrailer(m)
}
func (c *customServerStream) Context() context.Context {
	//	return c.stream.Context()
	return c.ctx
}
func (c *customServerStream) SendMsg(m interface{}) error {
	return c.stream.SendMsg(m)
}
func (c *customServerStream) RecvMsg(m interface{}) error {
	return c.stream.RecvMsg(m)
}

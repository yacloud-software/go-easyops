package server

import (
	"context"
	"fmt"

	fw "golang.conradwood.net/apis/framework"

	//	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"

	//	"golang.conradwood.net/go-easyops/ctx"
	"strings"
	"time"

	"golang.conradwood.net/go-easyops/errors"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (sd *serverDef) StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var err error
	pp.ServerRpcEntered()
	defer pp.ServerRpcDone()
	startRPC()
	defer stopRPC()

	name := ServiceNameFromStreamInfo(info)
	method := MethodNameFromStreamInfo(info)
	rc := &rpccall{ServiceName: name, MethodName: method, Started: time.Now()}
	stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": name,
	}).Inc()
	defer stdMetrics.concurrent_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": name,
	}).Dec()

	if cmdline.IsDebugRPCServer() {
		fmt.Printf("[go-easyops] Debug-rpc: called streaming rpc %s/%s\n", name, method)
	}
	//fmt.Printf("Method: \"%s\"\n", method)
	if isInternalService(name) {
		if cmdline.IsDebugRPCServer() {
			fmt.Printf("Invoking internal service stream handler\n")
		}
		res := handler(srv, stream)
		if cmdline.IsDebugRPCServer() {
			fmt.Printf("internal service stream handler returned: %s\n", res)
		}
		return res
	}

	def := getServerDefByName(name)
	if def == nil {
		s := fmt.Sprintf("[go-easyops] Service not registered! %s", name)
		fmt.Println(s)
		return errors.Error(stream.Context(), codes.Unimplemented, "service unavailable", "service %s is not known here", rc.ServiceName)
	}

	grpc_server_requests.With(prometheus.Labels{
		"method":      method,
		"servicename": def.name,
	}).Inc()

	var out_ctx context.Context

	if cmdline.ContextWithBuilder() {
		out_ctx, _, err = sd.V1inbound2outbound(stream.Context(), rc)
		if err != nil {
			return err
		}
	} else {
		panic("obsolete code path")
	}
	track_inbound_call(name, method, auth.GetService(out_ctx))
	print_inbound_debug(rc, out_ctx)

	nstream := newServerStream(stream, out_ctx)
	err = handler(srv, nstream)
	if err == nil {
		return nil
	}
	if cmdline.IsDebugRPCServer() || *print_errs {
		fmt.Printf("[go-easyops] Call %s.%s failed: %s\n", def.name, method, errors.ErrorStringWithStackTrace(err))
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
		sd.logError(out_ctx, rc, err)
		return err
	}

	re := st.Err()
	sd.logError(out_ctx, rc, re)
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

package client

import (
	_ "context"
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

//ClientMetricsUnaryInterceptor intercepts unary grpc calls made by clients and adds a timestamp to context metadata that can be used by server handlers/interceptors to determine grpc durations from the time a client sent a request to the time the server is ready with a response. This avoids the need for a push gateway for client metrics.
func ClientMetricsUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if ctx == nil {
		panic("missing context")
	}
	pp.ClientRpcEntered()
	defer pp.ClientRpcDone()
	start := time.Now()
	s, m, err := splitMethodAndService(method)
	// keep prometheus happy: (not to panic)
	if s == "" {
		s = "unknown"
	}
	if m == "" {
		m = "unknown"
	}
	// bad programmer! - log it prominently (often)
	if err != nil {
		fmt.Printf("[go-easyops] invalid fqdn method: \"%s\": %s\n", method, err)
	}
	cs := rpc.CallStateFromContext(ctx)
	if cs == nil {
		if *dialer_debug || *debug_rpc_client {
			fmt.Printf("[go-easyops] WARNING - calling external method %s.%s() without callstate\n", s, m)
		}
	} else {
		if !isKnownNotAuthRPCs(s, m) {
			_, ex := metadata.FromOutgoingContext(ctx)
			if !ex {
				fmt.Printf("[go-easyops] WARNING - calling external method %s.%s without metadata (authentication)\n", s, m)
			}
		}
	}
	if *dialer_debug || *debug_rpc_client {
		cs.PrintContext()
		us := "none"
		if cs.User() != nil {
			us = fmt.Sprintf("UserID=%s, Email=%s", cs.User().ID, cs.User().Email)
		}
		fmt.Printf("Invoking method %s.%s as %s...\n", s, m, us)
	}
	grpc_client_sent.With(prometheus.Labels{"method": m, "servicename": s}).Inc()
	err = invoker(ctx, method, req, reply, cc, opts...)
	observeRPC(start, s, m)

	ts := fmt.Sprintf("%fms", time.Since(start).Seconds()*1000)
	if err != nil {
		grpc_client_failed.With(prometheus.Labels{"method": m, "servicename": s}).Inc()
		if *dialer_debug || *debug_rpc_client {
			fmt.Printf("Invoke remote method=%s duration=%v error=%v (Method: \"%s\" in Service: \"%s\")\n", method, time.Since(start)/1e6, err, m, s)
		}
	} else if *debug_rpc_client {
		fmt.Printf("Invoked method %s.%s (%s)...\n", s, m, ts)
	}

	return err
}

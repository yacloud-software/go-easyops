package client

import (
	"context"
	"flag"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"google.golang.org/grpc"
)

const (
	// if we have a long rpc call (e.g. 3 seconds)
	// then we can run into timeouts, during the
	// interceptor auth phase
	// imho - the auth should be handled by "normal"
	// loadbalancer function
	// (it seems) cnw 19/5/2018
	CONST_CALL_TIMEOUT = 4
)

var (
	normal_sleep_time = flag.Int("ge_dialer_sleep_time", 20, "interval in seconds before querying the registry for changes (should be lower than ge_max_block)")

	blockCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_loadbalancer_no_connection",
			Help: "counter incremented each time the loadbalancer has no instances",
		},
		[]string{"servicename"},
	)

	failedQueryCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_loadbalancer_registry_failures",
			Help: "counter incremented each time the loadbalancer fails to query the registry",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(blockCtr, failedQueryCtr)
}
func unaryStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	pp.ClientRpcEntered()
	cs, err := streamer(ctx, desc, cc, method, opts...)
	pp.ClientRpcDone()
	return cs, err
}

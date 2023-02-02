package client

import (
	"flag"
	"golang.conradwood.net/go-easyops/prometheus"
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

package profiling

import (
	"flag"
	"golang.conradwood.net/go-easyops/prometheus"
	"time"
)

var (
	rpcCtr       = 0
	sqlCtr       = 0
	serverRPCCtr = 0
	ms           = flag.Int("ge_profiling_interval", 300, "interval in `milliseconds` to collect codepath stuff")
	cm_total     = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ge_profiling_total_samples",
			Help: "total number of codepath samples (whilst at least one server grpc is being executed)",
		},
		[]string{"rpcactive"},
	)
	cm_ctr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ge_profiling_samples_wait",
			Help: "total number of codepath samples where at least one thread was blocked",
		},
		[]string{"codepath"},
	)
	gm_ctr = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ge_profiling_current_blocks",
			Help: "threads within codepath",
		},
		[]string{"codepath"},
	)
)

func init() {
	prometheus.MustRegister(cm_total, cm_ctr, gm_ctr)
	go statsCheck()
}

func statsCheck() {
	for {
		gm_ctr.With(prometheus.Labels{"codepath": "grpccall"}).Set(float64(rpcCtr))
		gm_ctr.With(prometheus.Labels{"codepath": "sqlquery"}).Set(float64(sqlCtr))
		if serverRPCCtr == 0 {
			cm_total.With(prometheus.Labels{"rpcactive": "false"}).Inc()
			time.Sleep(time.Duration(*ms) * time.Millisecond)
			continue
		}
		cm_total.With(prometheus.Labels{"rpcactive": "true"}).Inc()

		// below we copy each variable into a local copy
		// so to avoid race conditions between "if" statement and Add()
		a := serverRPCCtr
		if a > 0 {
			cm_ctr.With(prometheus.Labels{"codepath": "serving"}).Add(float64(a))
		}
		a = rpcCtr
		if a > 0 {
			cm_ctr.With(prometheus.Labels{"codepath": "grpccall"}).Add(float64(a))
		}
		a = sqlCtr
		if a > 0 {
			cm_ctr.With(prometheus.Labels{"codepath": "sqlquery"}).Add(float64(a))
		}
		time.Sleep(time.Duration(*ms) * time.Millisecond)
	}
}

func ClientRpcEntered() {
	rpcCtr++
}
func ClientRpcDone() {
	rpcCtr--
}
func ServerRpcEntered() {
	serverRPCCtr++
}
func ServerRpcDone() {
	serverRPCCtr--
}
func SqlEntered() {
	sqlCtr++
}
func SqlDone() {
	sqlCtr--
}

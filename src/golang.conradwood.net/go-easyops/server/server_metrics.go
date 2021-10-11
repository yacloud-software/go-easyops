package server

import (
	//	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/prometheus"
	"time"
)

var (
	grpc_server_req_durations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "grpc_server_request_durations",
			Help:       "V=1 UNIT=s DESC=RPC latency distributions, measured on server",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			MaxAge:     time.Hour,
		},
		[]string{"servicename", "method"},
	)
	grpc_failed_requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_failed",
			Help: "number of grpc requests failed by this server",
		},
		[]string{"servicename", "method", "grpccode"},
	)
	grpc_server_requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_received",
			Help: "total number of grpc requests received by this server",
		},
		[]string{"servicename", "method"},
	)
)

func init() {
	prometheus.MustRegister(grpc_server_req_durations, grpc_failed_requests, grpc_server_requests)
}

// deprecated - to be moved to global vars
type ServerMetrics struct {
	concurrent_server_requests *prometheus.GaugeVec
	inv_auth                   *prometheus.CounterVec
}

// deprecated - to be moved to global vars
func NewServerMetrics() *ServerMetrics {
	res := &ServerMetrics{
		concurrent_server_requests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "grpc_concurrent_requests",
				Help: "how many rpcs are served concurrently at any point in time",
			},
			[]string{"servicename", "method"},
		),
		inv_auth: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_invalid_auth",
				Help: "V=1 UNIT=ops DESC=RPC requests with invalid or obsolete authentication",
			},
			[]string{"servicename", "method", "reason"},
		),
	}
	prometheus.MustRegister(res.inv_auth, res.concurrent_server_requests)
	return res
}

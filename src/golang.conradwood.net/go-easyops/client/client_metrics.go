package client

import (
	"golang.conradwood.net/go-easyops/prometheus"
	"time"
)

var (
	grpc_client_responsetime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "grpc_client_responsetime",
			Help: "V=1 unit=s DESC=responsetimes of services called",
		},
		[]string{"servicename", "method"},
	)
)

func init() {
	prometheus.MustRegister(grpc_client_responsetime)
}
func observeRPC(started time.Time, servicename, method string) {
	diff := time.Since(started).Seconds()
	l := prometheus.Labels{"method": method, "servicename": servicename}
	grpc_client_responsetime.With(l).Observe(diff)
}

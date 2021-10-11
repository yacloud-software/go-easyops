package client

import (
	_ "context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"google.golang.org/grpc"
	"os"
	"strings"
)

const (
	use_fancy_balancer = true
)

var (
	known_not_auth_rpcs = []string{
		"rpcinterceptor.RPCInterceptorService.InterceptRPC",
	}
	// I think part of a refactoring, the metrics below
	// should move into a metrics package, together with
	// the server metrics.
	// then we should have a single metric:
	// "grpc_requests_total{direction="sent|received"}
	// cnw
	grpc_client_sent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_sent",
			Help: "V=1 unit=ops total number of grpc requests sent by this instance",
		},
		[]string{"servicename", "method"},
	)
	grpc_client_failed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_sent_failed",
			Help: "V=1 unit=ops total number of grpc requests sent by this instance and failed",
		},
		[]string{"servicename", "method"},
	)
	debug_rpc_client = flag.Bool("ge_debug_rpc_client", false, "set to true to debug remote invokations")
	dialer_debug     = flag.Bool("ge_debug_dialer", false, "set to true to debug the grpc dialer")
)

func init() {
	prometheus.MustRegister(grpc_client_sent, grpc_client_failed)
}

// opens a tcp connection to an ip (no loadbalancing obviously)
func ConnectWithIPNoBlock(ip string) (*grpc.ClientConn, error) {
	return connectWithIPOptions(ip, false)
}
func ConnectWithIP(ip string) (*grpc.ClientConn, error) {
	return connectWithIPOptions(ip, true)
}
func connectWithIPOptions(ip string, block bool) (*grpc.ClientConn, error) {
	if *dialer_debug {
		fmt.Println("DialService: Dialling " + ip + " and blocking until successful connection...")
	}
	var err error
	var conn *grpc.ClientConn
	if block {
		conn, err = grpc.Dial(
			ip,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(GetClientCreds()),
			grpc.WithUnaryInterceptor(ClientMetricsUnaryInterceptor),
			grpc.WithStreamInterceptor(unaryStreamInterceptor),
		)
	} else {
		conn, err = grpc.Dial(
			ip,
			grpc.WithTransportCredentials(GetClientCreds()),
			grpc.WithUnaryInterceptor(ClientMetricsUnaryInterceptor),
			grpc.WithStreamInterceptor(unaryStreamInterceptor),
		)
	}
	if err != nil {
		return nil, err
	}

	if *dialer_debug {
		fmt.Printf("Connected to %s\n", ip)
	}

	return conn, nil

}

//
func Connect(serviceNameOrPath string) *grpc.ClientConn {
	return ConnectAt(cmdline.GetClientRegistryAddress(), serviceNameOrPath)
}

// convenience method to get a loadbalanced connection to a service
// use path or servicename (path prefered, it contains the version)
// unless it successfullly connects it will NOT return
// (it will either terminate the process or loop)
func ConnectAt(registryadr string, serviceNameOrPath string) *grpc.ClientConn {
	common.AddBlockedServiceName(serviceNameOrPath)
	conn, err := dialService(registryadr, serviceNameOrPath)
	// an error in this case reflects a LOCAL error, such as
	// no route to host or out-of-memory.
	// if a service is not available at the time of the call
	// it will block until one becomes available.
	// since it is a local error it is appropriate to exit.
	// a system administrator has to repair the machine before
	// the software can continue.
	if err != nil {
		fmt.Printf("Failed to dial %s: %s\n", serviceNameOrPath, err)
		os.Exit(10)
	}
	if *dialer_debug {
		fmt.Printf("[go-easyops]Connected to %s\n", serviceNameOrPath)
	}
	common.RemoveBlockedServiceName(serviceNameOrPath)
	return conn
}

// opens a tcp connection to a path.
func dialService(registry string, serviceName string) (*grpc.ClientConn, error) {
	GetSignatureFromAuth() // this is triggered here, because we _must_ have a valid signature later. if it has been called earlire it is a noop
	if *dialer_debug {
		fmt.Println("DialService: Dialling " + serviceName + " and blocking until successful connection...")
	}
	var err error
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(
		"go-easyops://"+serviceName+"/"+serviceName+"@"+registry, // "go-easyops://" url scheme registered in fancy_resolver.go
		grpc.WithContextDialer(CustomDialer),                     // custom dialer to distinguish between direct and proxy tcp connections
		grpc.WithBlock(),                                         // do not return until at least one connection is up
		grpc.WithBalancerName("fancybalancer"),                   // "fancybalancer" registered in fancy_balancer.go
		grpc.WithTransportCredentials(GetClientCreds()),          // transport credentials: default hardcoded certificates
		grpc.WithUnaryInterceptor(ClientMetricsUnaryInterceptor), // this is called for every unary RPC
		grpc.WithStreamInterceptor(unaryStreamInterceptor),       // this is called for every stream RPC
	)

	if err != nil {
		return nil, err
	}

	if *dialer_debug {
		fmt.Printf("Connected to %s\n", serviceName)
	}

	return conn, nil
}

// given a fqdn like so:
// "/auth.AuthenticationService/GetByToken"
// it'll return service and method as strings
func splitMethodAndService(fqdn string) (string, string, error) {
	ms := strings.Split(fqdn, "/")
	if len(ms) != 3 {
		return "", "", fmt.Errorf("%s is not a valid service name (contains %d parts instead of 3)", fqdn, len(ms))
	}
	return ms[1], ms[2], nil
}
func isKnownNotAuthRPCs(s, m string) bool {
	sn := fmt.Sprintf("%s.%s", s, m)
	for _, k := range known_not_auth_rpcs {
		if k == sn {
			return true
		}
	}
	return false
}

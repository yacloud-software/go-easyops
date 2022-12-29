/*
Package standalone handles applications running in "standalone" mode (that is without failover/datacenter etc)
*/
package standalone

// a standalone "registry replacement"

import (
	"context"
	"flag"
	"fmt"
	reg "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"net"
	"strings"
)

var (
	debug = flag.Bool("ge_debug_standalone", false, "if true debug the standalone dialler")
	// if we receive requests to resolve any of those (and they don't happen to be available), we panic
	UNSUPPORTED_SERVICES = []string{"rpcinterceptor.RPCInterceptorService"}
)

// register a service with the local "registry"
func RegisterService(rsr *reg.RegisterServiceRequest) (string, error) {
	dir := cmdline.LocalRegistrationDir()
	if !utils.FileExists(dir) {
		utils.Bail("failed to create register dir", utils.RecreateSafely(dir))
	}
	b, err := utils.MarshalBytes(rsr)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal proto (BUG!): %s", err))
	}
	fname := fmt.Sprintf("%s/%s.proto", dir, rsr.ServiceName)
	t := utils.FileExists(fname)
	err = utils.WriteFile(fname, b)
	if err != nil {
		panic(fmt.Sprintf("failed to write standalone registration details to %s: %s", fname, err))
	}
	if !t {
		fmt.Printf("[go-easyops] Registered service \"%s\" in %s\n", rsr.ServiceName, fname)
	}
	return "local", nil
}

// dial a service (using local "registry" to find the local port)
// serviceurl is something like "direct://[host]:port"
func DialService(ctx context.Context, serviceurl string) (net.Conn, error) {
	if !strings.HasPrefix(serviceurl, "direct://") {
		panic(fmt.Sprintf("[go-easyops] standalone unsupported serviceurl \"%s\"", serviceurl))
	}
	dialstring := strings.TrimPrefix(serviceurl, "direct://")
	if *debug {
		fmt.Printf("[go-easyops] standalone dialler dialing \"%s\"\n", dialstring)
	}
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", dialstring)

	if err != nil {
		fmt.Printf("[go-easyops] standalone dialling error: %s\n", err)
	}
	if *debug {
		fmt.Printf("[go-easyops] standalone dialler connection established to \"%s\"\n", dialstring)
	}
	return conn, err
}

// find all registrations for a given service
func Registry_V2GetTarget(ctx context.Context, req *reg.V2GetTargetRequest) (*reg.V2GetTargetResponse, error) {
	if *debug {
		fmt.Printf("[go-easyops] standalone resolving service \"%s\"...\n", req.ServiceName)
	}
	res := &reg.V2GetTargetResponse{}
	dir := cmdline.LocalRegistrationDir()
	for _, serviceName := range req.ServiceName {
		fname := fmt.Sprintf("%s/%s.proto", dir, serviceName)
		if !utils.FileExists(fname) {
			fmt.Printf("[go-easyops] standalone resolver: file %s does not exist\n", fname)
			is_well_known_unsupported(serviceName)
			return res, nil
		}
		b, err := utils.ReadFile(fname)
		if err != nil {
			fmt.Printf("[go-easyops] standalone resolver: cannot read file %s: %s\n", fname, err)
			return nil, err
		}
		rsr := &reg.RegisterServiceRequest{}
		err = utils.UnmarshalBytes(b, rsr)
		if err != nil {
			fmt.Printf("[go-easyops] standalone resolver: cannot parse file %s: %s\n", fname, err)
			return nil, err
		}
		t, err := registration_request_to_target(rsr)
		if err != nil {
			fmt.Printf("[go-easyops] standalone resolver: failed to convert file %s: %s\n", fname, err)
			return nil, err
		}
		res.Targets = append(res.Targets, t)
	}
	if len(res.Targets) == 0 {
		for _, s := range req.ServiceName {
			is_well_known_unsupported(s)
		}
	}

	return res, nil
}
func is_well_known_unsupported(s string) {
	for _, u := range UNSUPPORTED_SERVICES {
		if u == s {
			panic(fmt.Sprintf("[go-easyops] standalone does not support well-known service \"%s\".", s))
		}
	}
}
func registration_request_to_target(rsr *reg.RegisterServiceRequest) (*reg.Target, error) {
	res := &reg.Target{
		IP:          "localhost",
		ServiceName: rsr.ServiceName,
		Port:        rsr.Port,
		ApiType:     []reg.Apitype{reg.Apitype_grpc},
		RoutingInfo: &reg.RoutingInfo{},
	}
	return res, nil
}

package client

import (
	"flag"
	"fmt"
	"strings"
)

var (
	nodebugServices = []string{
		"registry.Registry",
		"auth.AuthenticationService",
		"auth.AuthManagerService",
		"rpcinterceptor.RPCInterceptorService",
		"errorlogger.ErrorLogger",
	}
	debug_fancy = flag.Bool("ge_debug_fancy_dialer", false, "debug the fancy resolver and balancer")
)

type serviceNamer interface {
	ServiceName() string
}

func fancyPrintf(sn serviceNamer, msg string, args ...interface{}) {
	if !*debug_fancy {
		return
	}
	s := sn.ServiceName()
	for _, f := range nodebugServices {
		if strings.HasPrefix(s, f) {
			return
		}
	}
	l := "[go-easyops] [" + s + "] "
	fmt.Printf(l+msg, args...)
}

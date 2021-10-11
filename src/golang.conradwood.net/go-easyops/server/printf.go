package server

import (
	"flag"
	"fmt"
)

var (
	fancyprintf = flag.Bool("ge_rpc_server_quiet", true, "if false, server will be very quiet and not print normal output")
	setquiet    = false
)

func SetQuietMode() {
	setquiet = true
}

func fancyPrintf(msg string, args ...interface{}) {
	if !*fancyprintf {
		return
	}
	if setquiet {
		return
	}

	l := "[go-easyops] "
	fmt.Printf(l+msg, args...)
}

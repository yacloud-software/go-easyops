package server

import (
	"context"
	"fmt"

	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/ctx"
)

func print_inbound_debug(rc *rpccall, myctx context.Context) {
	if !cmdline.IsDebugRPCServer() {
		return
	}
	s := rc.FullMethod()
	fmt.Printf("[go-easyops] Debug-RPC[%s]: (builder=%d) Invoked by user %s and service %s\n", s, cmdline.GetContextBuilderVersion(), auth.UserIDString(auth.GetUser(myctx)), auth.UserIDString(auth.GetService(myctx)))
	if auth.GetUser(myctx) == nil {
		fmt.Printf("[go-easyops] Context: %#v\n", ctx.Context2String(myctx))
	}

}

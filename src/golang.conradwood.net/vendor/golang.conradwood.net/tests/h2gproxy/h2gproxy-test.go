package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/create"
	"golang.conradwood.net/go-easyops/auth"
	_ "golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/tests/shared/rpc"
	"os"
)

func main() {
	flag.Parse()
	useremail := "easyops-test-user"
	ctx, err := rpc.ContextWithLogin(useremail, "easyops-test-password")
	utils.Bail("no login context", err)
	rpc.PrintContext(ctx)
	user := rpc.GetUser(ctx)
	fmt.Printf("Got user: %s\n", auth.Description(user))
	if user.Email != useremail {
		fail("User email mismatch: %s != %s\n", useremail, user.Email)
	}
	_, err = create.GetEasyOpsTestClient().SimplePing(ctx, &common.Void{})
	utils.Bail("Failed simpleping()", err)
}

func fail(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(10)
}

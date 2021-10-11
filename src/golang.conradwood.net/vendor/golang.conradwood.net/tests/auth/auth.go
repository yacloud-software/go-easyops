package main

import (
	//	au "golang.conradwood.net/apis/auth"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/create"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()

	ctx := auth.Context()
	_, err := auth.SerialiseContext(ctx)
	utils.Bail("failed to serialise context I just created", err)

	am := create.GetAuthManagerService()
	ctx = tokens.ContextWithToken()
	me, err := am.WhoAmI(ctx, &common.Void{})
	utils.Bail("failed to get me", err)
	fmt.Printf("Result of whoami():\n")
	auth.PrintUser(me)

	ctx, err = auth.ContextForUser(me)
	utils.Bail("failed to get context for user", err)
	me = auth.GetUser(ctx)
	fmt.Printf("Result of context for user with whoami():\n")
	auth.PrintUser(me)

	sctx, err := auth.SerialiseContext(ctx)
	fmt.Printf("Serialised: %s\n", sctx[:24])
	utils.Bail("failed to serialise context", err)
	ctx, err = auth.RecreateContext([]byte(sctx))
	utils.Bail("failed to deserialise context", err)
	me = auth.GetUser(ctx)
	fmt.Printf("Result of context serialise/deserialise:\n")
	auth.PrintUser(me)

	me, err = create.GetAuthManagerService().WhoAmI(ctx, &common.Void{})
	utils.Bail("failed to call authmanager.WhoAmI()", err)
	fmt.Printf("Result of calling AuthManager.WhoAmi():\n")
	auth.PrintUser(me)
	fmt.Printf("Done\n")

}

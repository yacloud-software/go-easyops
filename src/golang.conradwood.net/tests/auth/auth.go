package main

import (
	"flag"
	"fmt"
	au "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	cm "golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	flag.Parse()

	ctx := auth.Context(time.Duration(5) * time.Second)
	_, err := auth.SerialiseContext(ctx)
	utils.Bail("failed to serialise context I just created", err)

	am := authremote.GetAuthManagerClient()
	ctx = tokens.ContextWithToken()
	me, err := am.WhoAmI(ctx, &common.Void{})
	utils.Bail("failed to get me", err)
	fmt.Printf("Result of whoami():\n")
	auth.PrintUser(me)
	b := cm.VerifySignature(me)
	if b {
		fmt.Printf("Signature Valid\n")
	} else {
		fmt.Printf("Signature inalid!!\n")
	}

	ctx, err = auth.ContextForUser(me)
	utils.Bail("failed to get context for user", err)
	me = auth.GetUser(ctx)
	fmt.Printf("Result of context for user with whoami():\n")
	auth.PrintUser(me)

	sctx, err := auth.SerialiseContext(ctx)
	fmt.Printf("Serialised: %s\n", sctx[:24])
	utils.Bail("failed to serialise context", err)
	ctx, err = auth.RecreateContextWithTimeout(time.Duration(5)*time.Second, []byte(sctx))
	utils.Bail("failed to deserialise context", err)
	me = auth.GetUser(ctx)
	fmt.Printf("Result of context serialise/deserialise:\n")
	auth.PrintUser(me)

	me, err = authremote.GetAuthManagerClient().WhoAmI(ctx, &common.Void{})
	utils.Bail("failed to call authmanager.WhoAmI()", err)
	fmt.Printf("Result of calling AuthManager.WhoAmi():\n")
	auth.PrintUser(me)

	su, err := authremote.GetAuthManagerClient().SignedGetUserByID(ctx, &au.ByIDRequest{UserID: me.ID})
	utils.Bail("failed to call authmanager.WhoAmI()", err)
	fmt.Printf("Result of calling AuthManager.SignedGetUserByID():\n")
	auth.PrintSignedUser(su)

	fmt.Printf("Done\n")

}

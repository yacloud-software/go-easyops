package main

import (
	"context"
	"fmt"
	au "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

var (
	timeout    = time.Duration(10) * time.Second
	me_user    *au.SignedUser
	me_service *au.SignedUser
)

func client() {
	fmt.Printf("testing context...\n")

	// first check that new and old context carry the same user/service information once created.
	ctx_def := authremote.Context()
	fmt.Printf("context with version %d:\n", cmdline.GetContextBuilderVersion())
	printContext(ctx_def)
	me_user = auth.GetSignedUser(ctx_def)
	me_service = auth.GetSignedService(ctx_def)
	if me_user == nil && me_service == nil {
		fmt.Printf("failed no service and no user in context. cannot proceed\n")
		os.Exit(10)
	}

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	nctx := authremote.Context()
	fmt.Printf("context with version %d:\n", cmdline.GetContextBuilderVersion())
	printContext(nctx)
	utils.Bail("(1) ctx assertion", AssertEqualContexts(ctx_def, nctx))

	m := map[string]string{"foo": "bar"}
	nctx_derived := authremote.DerivedContextWithRouting(nctx, m, true)
	//	printContext(nctx)
	utils.Bail("(2) ctx assertion", AssertEqualContexts(nctx, nctx_derived))

	// check simple, new functions.
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	ctx1 := authremote.Context()
	s, err := ctx.SerialiseContextToString(ctx1)
	utils.Bail("(1) failed to serialise context to string", err)
	fmt.Printf("Serialised to: %s\n", utils.HexStr([]byte(s)))
	_, err = ctx.DeserialiseContextFromString(s)
	utils.Bail("(1) failed to deserialise context to string", err)
	_, err = ctx.DeserialiseContext([]byte(s))
	utils.Bail("(2) failed to deserialise context to string", err)
	_, err = auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, []byte(s))
	utils.Bail("(3) failed to deserialise context to string", err)

	checkFile("/tmp/context.env")
	b, err := auth.SerialiseContext(ctx1)
	utils.Bail("failed to serialise", err)
	//fmt.Println(utils.Hexdump("Context: ", b))
	ctx2, err := auth.RecreateContextWithTimeout(timeout, b)
	utils.Bail("failed to deserialise", err)
	mustBeSame(ctx1, ctx2)

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	fmt.Printf("Checking 'derived' context with routing...\n")
	ctx1 = authremote.Context()
	ctx2 = authremote.DerivedContextWithRouting(ctx1, make(map[string]string), true)
	mustBeSame(ctx1, ctx2)
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	fmt.Printf("Checking 'derived' context with routing...\n")
	ctx1 = authremote.Context()
	ctx2 = authremote.DerivedContextWithRouting(ctx1, make(map[string]string), true)
	mustBeSame(ctx1, ctx2)

	// now check all the various methods of deserialising
	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	ctx1 = authremote.Context()
	checkSer(ctx1)

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	ctx1 = authremote.Context()
	checkSer(ctx1)

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	ctx1 = authremote.Context()
	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	checkSer(ctx1)

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	ctx1 = authremote.Context()
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	checkSer(ctx1)

}

// check if serialisation and deserialisation arrives at the same context as the one passed in
func checkSer(ctx1 context.Context) {
	cmdline.SetEnvContext("")
	b, err := auth.SerialiseContext(ctx1)
	utils.Bail("failed to serialise", err)
	fmt.Printf("Serialised to: %s\n", utils.HexStr(b))
	ctx2, err := auth.RecreateContextWithTimeout(timeout, b)
	utils.Bail("failed to deserialise", err)
	mustBeSame(ctx1, ctx2)

	cmdline.SetEnvContext(string(b))
	ctx2 = authremote.Context()
	mustBeSame(ctx1, ctx2)
	cmdline.SetEnvContext("")

	s, err := auth.SerialiseContextToString(ctx1)
	utils.Bail("failed to serialise", err)
	cmdline.SetEnvContext(s)
	ctx2 = authremote.Context()
	mustBeSame(ctx1, ctx2)

	_, err = auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, []byte(s))
	if err != nil {
		utils.PrintStack("deser string fail")
	}
	utils.Bail("failed to deserialise string", err)
	cmdline.SetEnvContext("")

}

func printContext(ictx context.Context) {
	u := auth.GetUser(ictx)
	s := auth.GetService(ictx)
	fmt.Printf("User   : %s\n", auth.UserIDString(u))
	fmt.Printf("Service: %s\n", auth.UserIDString(s))
	fmt.Printf("ctx2str: %s\n", ctx.Context2String(ictx))
}
func mustBeSame(ctx1, ctx2 context.Context) {
	u1 := auth.GetUser(ctx1)
	u2 := auth.GetUser(ctx2)
	if auth.UserIDString(u1) != auth.UserIDString(u2) {
		panic(fmt.Sprintf("users do not match (%s vs %s)", auth.UserIDString(u1), auth.UserIDString(u2)))
	}
	s1 := auth.GetService(ctx1)
	s2 := auth.GetService(ctx2)
	if auth.UserIDString(s1) != auth.UserIDString(s2) {
		panic("services do not match")
	}
	//fmt.Printf("Contexts match\n")
}

func checkFile(filename string) {
	b, err := utils.ReadFile(filename)
	if err != nil {
		fmt.Printf("ignoring file \"%s\", could not read it\n", filename)
		return
	}
	ctx1, err := ctx.DeserialiseContext(b)
	if err != nil {
		fmt.Printf("Failed to deserialise context in file \"%s\": %s\n", filename, err)
		os.Exit(10)
	}
	fmt.Printf("Serialised Context from file %s:\n", filename)
	printContext(ctx1)
}

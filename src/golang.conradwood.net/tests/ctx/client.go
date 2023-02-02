package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

var (
	timeout = time.Duration(10) * time.Second
)

func client() {
	fmt.Printf("testing context...\n")
	checkFile("/tmp/context.env")
	ctx1 := authremote.Context()
	fmt.Printf("Default context:\n")
	printContext(ctx1)
	b, err := auth.SerialiseContext(ctx1)
	utils.Bail("failed to serialise", err)
	//fmt.Println(utils.Hexdump("Context: ", b))
	ctx2, err := auth.RecreateContextWithTimeout(timeout, b)
	utils.Bail("failed to deserialise", err)
	mustBeSame(ctx1, ctx2)

	cmdline.SetContextWithBuilder(false)
	fmt.Printf("Checking 'derived' context with routing...\n")
	ctx1 = authremote.Context()
	ctx2 = authremote.DerivedContextWithRouting(ctx1, make(map[string]string), true)
	mustBeSame(ctx1, ctx2)
	cmdline.SetContextWithBuilder(true)
	fmt.Printf("Checking 'derived' context with routing...\n")
	ctx1 = authremote.Context()
	ctx2 = authremote.DerivedContextWithRouting(ctx1, make(map[string]string), true)
	mustBeSame(ctx1, ctx2)

	// now check all the various methods of deserialising
	cmdline.SetContextWithBuilder(false)
	ctx1 = authremote.Context()
	checkSer(ctx1)

	cmdline.SetContextWithBuilder(true)
	ctx1 = authremote.Context()
	checkSer(ctx1)

	cmdline.SetContextWithBuilder(true)
	ctx1 = authremote.Context()
	cmdline.SetContextWithBuilder(false)
	checkSer(ctx1)

	cmdline.SetContextWithBuilder(false)
	ctx1 = authremote.Context()
	cmdline.SetContextWithBuilder(true)
	checkSer(ctx1)

}

// check if serialisation and deserialisation arrives at the same context as the one passed in
func checkSer(ctx1 context.Context) {
	cmdline.SetEnvContext("")
	b, err := auth.SerialiseContext(ctx1)
	utils.Bail("failed to serialise", err)
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

	cmdline.SetEnvContext("")

}

func printContext(ictx context.Context) {
	u := auth.GetUser(ictx)
	s := auth.GetService(ictx)
	fmt.Printf("User   : %s\n", auth.UserIDString(u))
	fmt.Printf("Service: %s\n", auth.UserIDString(s))
}
func mustBeSame(ctx1, ctx2 context.Context) {
	u1 := auth.GetUser(ctx1)
	u2 := auth.GetUser(ctx2)
	if auth.UserIDString(u1) != auth.UserIDString(u2) {
		panic("users do not match")
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

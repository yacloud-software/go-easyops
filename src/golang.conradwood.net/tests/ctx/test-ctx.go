package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	flag.Parse()
	fmt.Printf("testing context...\n")
	t := time.Duration(10) * time.Second
	ctx1 := authremote.Context()
	printContext(ctx1)
	b, err := auth.SerialiseContext(ctx1)
	utils.Bail("failed to serialise", err)
	fmt.Println(utils.Hexdump("Context: ", b))
	ctx2, err := auth.RecreateContextWithTimeout(t, b)
	utils.Bail("failed to deserialise", err)
	mustBeSame(ctx1, ctx2)
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
	fmt.Printf("Contexts match\n")
}

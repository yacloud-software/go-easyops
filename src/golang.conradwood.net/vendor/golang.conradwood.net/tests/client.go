package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/create"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	direct_rpc = flag.Bool("direct_rpc", false, "do not go through interceptors and go-easyops, open connection to service directly (with LB)")
	cl         ge.EasyOpsTestClient
)

func initClient() {
	if !*direct_rpc {
		cl = create.NewEasyOpsTestClient()
		return
	}
	serviceName := "goeasyops.EasyOpsTest"
	conn, err := grpc.Dial(
		"go-easyops://"+serviceName+"/"+serviceName,
		grpc.WithBlock(),
		grpc.WithBalancerName("fancybalancer"),
		grpc.WithTransportCredentials(client.GetClientCreds()),
	)
	utils.Bail("Failed to dial", err)
	cl = ge.NewEasyOpsTestClient(conn)
}
func StartClient() {
	initClient()
xloop:
	ctx := tokens.ContextWithToken()
	me := auth.WhoAmI()
	if me == nil {
		fmt.Printf("authtest-client: pinging as 'nobody'. (new context did not provide a user)\n")
	} else {
		fmt.Printf("authtest-client: pinging as %s\n", me.Email)
	}
	r, err := cl.Ping(ctx, &ge.Chain{})
	if err != nil {
		fmt.Printf("failed to ping with standard token: %s\n", err)
		if *loop {
			time.Sleep(1 * time.Second)
			goto xloop
		}
		os.Exit(10)

	}
	fmt.Printf("%d reports\n", len(r.Calls))
	ft := "%5s | %10s | %10s | %s\n"
	fmt.Printf(ft, "#", "reqid", "userid", "serviceid")
	for _, c := range r.Calls {
		user, err := auth.GetUserByID(ctx, c.UserID)
		utils.Bail("Failed to get user", err)
		service, err := auth.GetUserByID(ctx, c.ServiceID)
		sn := c.ServiceID
		if err == nil {
			sn = auth.Description(service)
		}
		fmt.Printf(ft,
			fmt.Sprintf("%d", c.Position),
			c.RequestID,
			auth.Description(user),
			sn,
		)
	}
	if *loop {
		time.Sleep(time.Duration(300) * time.Millisecond)
		goto xloop
	}
	fmt.Printf("OK\n")
}

package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/helloworld"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

func main() {
	flag.Parse()
	fmt.Printf("go-easyops test client\n")
	ctx := authremote.Context()
	if ctx == nil {
		fmt.Printf("ERROR: authremote.Context() created no context\n")
		os.Exit(10)
	}
	fmt.Printf("Pinging with default client...\n")
	started := time.Now()
	r, err := helloworld.GetHelloWorldClient().Ping(ctx, &common.Void{})
	utils.Bail("failed to ping", err)
	fmt.Printf("Pinged (%0.2fs), User=%s, Service=%s\n", time.Since(started).Seconds(), auth.Description(r.CallingUser), auth.Description(r.CallingService))

	fmt.Printf("Pinging with lookup...\n")
	started = time.Now()
	con := client.Connect("helloworld.HelloWorld")
	c := helloworld.NewHelloWorldClient(con)
	_, err = c.Ping(ctx, &common.Void{})
	utils.Bail("failed to ping", err)
	fmt.Printf("Pinged (%0.2fs), User=%s, Service=%s\n", time.Since(started).Seconds(), auth.Description(r.CallingUser), auth.Description(r.CallingService))
	/*
		fmt.Printf("Pinging stream...\n")
		psreq := &helloworld.PingStreamRequest{DelayInMillis: 500}
		srv, err := helloworld.GetHelloWorldClient().PingStream(ctx, psreq)
		utils.Bail("failed to set up pingstream", err)
		for {
			pr, err := srv.Recv()
			if err != nil {
				fmt.Printf("error received: %s\n", err)
				break
			}
			fmt.Printf("Received Sequence %d\n", pr.SequenceNumber)
		}
	*/
}

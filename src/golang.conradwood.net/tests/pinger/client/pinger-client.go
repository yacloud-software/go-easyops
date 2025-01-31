package main

import (
	"flag"
	"fmt"
	"time"

	pb "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	for {

		ctx := authremote.Context()
		response, err := pb.GetEchoServiceClient().Ping(ctx, &pb.PingRequest{})
		utils.Bail("failed to ping", err)
		fmt.Printf("%s Response: %#v\n", utils.TimeString(time.Now()), response)
		time.Sleep(1 * time.Second)
		// only services can impersonate
		/*
			ctx, err = authremote.ContextForUserIDWithTimeout("1", time.Duration(30)*time.Second)
			utils.Bail("failed to get context for user 1", err)
			response, err = pb.GetEchoServiceClient().Ping(ctx, &common.Void{})
			utils.Bail("failed to ping", err)
			fmt.Printf("Response: %#v\n", response)
			time.Sleep(1 * time.Second)
		*/
	}
}

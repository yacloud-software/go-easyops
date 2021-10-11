package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/echoservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	flag.Parse()
	for {

		ctx := authremote.Context()
		response, err := pb.GetEchoServiceClient().Ping(ctx, &common.Void{})
		utils.Bail("failed to ping", err)
		fmt.Printf("Response: %#v\n", response)
		time.Sleep(1 * time.Second)
		ctx, err = authremote.ContextForUserIDWithTimeout("1", time.Duration(30)*time.Second)
		utils.Bail("failed to get context for user 1", err)
		response, err = pb.GetEchoServiceClient().Ping(ctx, &common.Void{})
		utils.Bail("failed to ping", err)
		fmt.Printf("Response: %#v\n", response)
		time.Sleep(1 * time.Second)
	}
}

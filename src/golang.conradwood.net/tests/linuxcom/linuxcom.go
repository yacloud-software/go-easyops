package main

import (
	"flag"
	"fmt"
	"time"

	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	fmt.Printf("linuxcom tests starting\n")
	com := NewCommand()
	ctx := authremote.Context()
	//	ci, err := com.Start(ctx, "/usr/bin/md5sum")
	ci, err := com.Start(ctx, "./test_com.sh")
	utils.Bail("failed to start", err)
	started := time.Now()
	err = ci.Wait(ctx)
	dur := time.Since(started)
	fmt.Printf("\n\nStopped after %0.1fs\n", dur.Seconds())
	utils.Bail("failed to wait", err)

}

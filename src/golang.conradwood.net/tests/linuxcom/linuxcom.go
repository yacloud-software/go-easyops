package main

import (
	"flag"
	"fmt"
	"time"

	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	fmt.Printf("linuxcom tests starting\n")
	mc, err := linux.MyCgroup()
	utils.Bail("failed to get my cgroup", err)
	fmt.Printf("My cgroup: %s\n", mc)
	com := linux.NewCommand()
	ctx := authremote.ContextWithTimeout(time.Duration(30) * time.Second)
	//	ci, err := com.Start(ctx, "/usr/bin/md5sum")
	ci, err := com.Start(ctx, "./test_com.sh")
	utils.Bail("failed to start", err)
	go func(c linux.Command) {
		time.Sleep(time.Duration(10) * time.Second)
		utils.Bail("sigint failed", c.SigInt())
		time.Sleep(time.Duration(15) * time.Second)
		utils.Bail("sigkill failed", c.SigKill())
	}(com)
	started := time.Now()
	err = ci.WaitAll(ctx)
	dur := time.Since(started)
	fmt.Printf("\n\nStopped after %0.1fs\n", dur.Seconds())
	utils.Bail("failed to wait", err)

}

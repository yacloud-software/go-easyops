package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	fmt.Printf("linuxcom tests starting\n")
	com := NewCommand()
	com.SetExecutable("/usr/bin/md5sum")
	utils.Bail("failed to start", com.Start())
	utils.Bail("failed to wait", com.Wait())
}

package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"os"
)

const (
	NEW_CONTEXT_VERSION = 2 // if new context version,set this to new version
	OLD_CONTEXT_VERSION = 2
)

var (
	run_server = flag.Bool("server", false, "if true run server, othwerise client")
	test_html  = flag.Bool("test_html", false, "if true do a test html render")
	ser        = flag.Bool("serialise", false, "serialiser test")
)

func main() {
	flag.Parse()
	if *ser {
		utils.Bail("failed to serialise", serialise())
		os.Exit(0)
	}
	if *run_server {
		start_server()
	} else if *test_html {
		testrenderer_rendertest()
	} else {
		client()
	}
	fmt.Printf("Done\n")
	os.Exit(0)
}

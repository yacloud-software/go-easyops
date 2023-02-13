package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	CONTEXT_VERSION = 2
)

var (
	run_server = flag.Bool("server", false, "if true run server, othwerise client")
)

func main() {
	flag.Parse()
	if *run_server {
		start_server()
	} else {
		client()
	}
	fmt.Printf("Done\n")
	os.Exit(0)
}

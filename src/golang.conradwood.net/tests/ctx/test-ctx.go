package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	run_server = flag.Bool("server", false, "if true run server, othwerise client")
)

func main() {
	flag.Parse()
	if *run_server {
		server()
	} else {
		client()
	}
	fmt.Printf("Done\n")
	os.Exit(0)
}

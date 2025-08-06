package main

import (
	"flag"

	"golang.conradwood.net/go-easyops/utils/functionchain"
)

func main() {
	flag.Parse()
	functionchain.NewFunctionChain()
}

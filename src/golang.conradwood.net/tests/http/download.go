package main

import (
	"flag"
	"golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	url := flag.Args()[0]
	h := http.HTTP{}
	hr := h.Get(url)
	err := hr.Error()
	utils.Bail("failed to get url", err)

}

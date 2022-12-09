package main

import (
	"flag"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

var (
	direct = flag.Bool("direct", false, "if true use direct access mode instead of urlcacher")
)

func main() {
	flag.Parse()
	url := flag.Args()[0]
	var h http.HTTPIF
	if *direct {
		h = http.NewDirectClient()
	} else {
		ctx := authremote.ContextWithTimeout(time.Duration(180) * time.Second)
		h = http.NewCachingClient(ctx)
	}
	hr := h.Get(url)
	err := hr.Error()
	utils.Bail("failed to get url", err)

}

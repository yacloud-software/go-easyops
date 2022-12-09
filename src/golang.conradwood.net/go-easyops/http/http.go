package http

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"
)

var (
	debug = flag.Bool("ge_debug_http", false, "debug http code")
)

type HTTPIF interface {
	Cookie(name string) *http.Cookie
	Cookies() []*http.Cookie
	Delete(url string, body []byte) *HTTPResponse
	Get(url string) *HTTPResponse
	GetStream(url string) *HTTPResponse
	Head(url string) *HTTPResponse
	Post(url string, body []byte) *HTTPResponse
	Put(url string, body string) *HTTPResponse
	SetHeader(key string, value string)
	SetTimeout(dur time.Duration)
}

// use urlcacher for the url (needs ctx to authenticate)
func NewCachingClient(ctx context.Context) HTTPIF {
	if *debug {
		fmt.Printf("New caching client..\n")
	}
	return cHTTP{ctx: ctx}
}

// retrieve directly from source
func NewDirectClient() HTTPIF {
	if *debug {
		fmt.Printf("New direct client..\n")
	}
	return HTTP{}
}

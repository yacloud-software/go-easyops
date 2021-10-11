package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/cache"
	"golang.conradwood.net/go-easyops/server"
	"time"
)

var (
	ch = cache.New("testcache", time.Duration(1)*time.Second, 999)
)

func main() {
	flag.Parse()
	// prometheus picks us up please
	server.StartFakeService("fakeservice_testcache")
	key := "FOOKEY"
	started := time.Now()
	for {
		diff := time.Since(started).Seconds()
		ds := fmt.Sprintf("%00.2f ", diff)
		o := ch.Get(key)
		if o == nil {
			fmt.Printf("%sNo cache entry\n", ds)
			ch.Put(key, "foo")
		} else {
			fmt.Printf("%sGot value\n", ds)
		}
		time.Sleep(time.Duration(250) * time.Millisecond)
	}
}

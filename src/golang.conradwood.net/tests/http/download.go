package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"time"
)

var (
	direct   = flag.Bool("direct", false, "if true use direct access mode instead of urlcacher")
	dur      = flag.Duration("duration", time.Duration(10)*time.Second, "max duration of http request context")
	timeout  = flag.Duration("timeout", time.Duration(5)*time.Second, "timeout of http request")
	testfile = flag.String("testfile", "", "if set, use this file as a list of urls to download from cache and directly and compare")
)

func main() {
	flag.Parse()
	if *testfile != "" {
		utils.Bail("failed test", TestFile())
		os.Exit(0)
	}
	url := flag.Args()[0]
	var h http.HTTPIF
	if *direct {
		h = http.NewDirectClient()
	} else {
		ctx := authremote.ContextWithTimeout(*dur)
		h = http.NewCachingClient(ctx)
	}
	started := time.Now()
	hr := h.Get(url)
	err := hr.Error()
	utils.Bail("failed to get url", err)
	dur := time.Since(started)
	fmt.Printf("Duration: %0.2fs\n", dur.Seconds())
}

func TestFile() error {
	b, err := utils.ReadFile(*testfile)
	if err != nil {
		return err
	}
	sx := strings.Split(string(b), "\n")
	for _, line := range sx {
		if len(line) < 3 {
			continue
		}
		if strings.Contains(line, "latest") {
			continue
		}
		if strings.Contains(line, "list") {
			continue
		}
		err = compare(line)
		if err != nil {
			return fmt.Errorf("url %s failed: %s", line, err)
		}
	}
	return nil
}
func compare(url string) error {
	fmt.Printf("Comparing %s..", url)
	fmt.Printf("fetching direct...")
	h := http.NewDirectClient()
	h.SetHeader("accept-encoding", "*")
	h.SetTimeout(*timeout)
	hr := h.Get(url)
	err := hr.Error()
	if err != nil {
		fmt.Printf("Body: %s\n", hr.Body())
		return fmt.Errorf("Unable to retrieve %s direct: %s", url, err)
	}
	b1 := hr.Body()

	fmt.Printf("fetching cached #1...")
	ctx := authremote.ContextWithTimeout(*dur)
	h = http.NewCachingClient(ctx)
	h.SetTimeout(*timeout)
	h.SetTimeout(*timeout)
	hr = h.Get(url)
	err = hr.Error()
	if err != nil {
		return fmt.Errorf("Unable to retrieve %s via 1st cached attempt: %s", url, err)
	}
	b2 := hr.Body()

	fmt.Printf("Comparing 1/2...")
	if !bytes.Equal(b1, b2) {
		return fmt.Errorf("URL %s - b1 (%d bytes)/b2 (%d bytes) mismatch", url, len(b1), len(b2))
	}

	fmt.Printf("fetching cached #2...")
	ctx = authremote.ContextWithTimeout(*dur)
	h = http.NewCachingClient(ctx)
	h.SetTimeout(*timeout)
	hr = h.Get(url)
	err = hr.Error()
	if err != nil {
		return fmt.Errorf("Unable to retrieve %s via 2nd cached attempt: %s", url, err)
	}
	b3 := hr.Body()

	fmt.Printf("Comparing 1/3...")
	if !bytes.Equal(b1, b3) {
		return fmt.Errorf("URL %s - b1 (%d bytes)/b3 (%d bytes) mismatch", url, len(b1), len(b3))
	}
	fmt.Printf("OK\n")
	return nil
}

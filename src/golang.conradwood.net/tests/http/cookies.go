package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/http"
	ht "net/http"
)

func TestCookie() error {
	if len(flag.Args()) == 0 {
		return fmt.Errorf("missing url")
	}
	url := flag.Args()[0]
	h := http.NewDirectClient()
	h.SetDebug(true)
	hr := h.Get(url)
	err := hr.Error()
	if err != nil {
		return err
	}
	ck := hr.Cookies()
	print_cookies(ck)

	hr = h.Get(url)
	err = hr.Error()
	if err != nil {
		return err
	}
	ck = hr.Cookies()
	print_cookies(ck)
	return nil
}
func print_cookies(ck []*ht.Cookie) {
	fmt.Printf("Received %d cookies:\n", len(ck))
	for _, c := range ck {
		fmt.Printf("%s == %s\n", c.Name, c.Value)
	}
}

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	service = flag.String("service", "", "Service to shutdown, either host:port or name")
)

func main() {
	flag.Parse()
	target := *service
	if target == "" {
		target = flag.Arg(0)
	}
	if target == "" {
		fmt.Printf("Please provide service name or host:port with --service option\n")
		os.Exit(10)
	}
	fmt.Printf("Shutting down service: %s\n", target)
	if strings.Contains(target, ":") {
		doHttp(target)
	} else {
		doRegistry()
	}
}
func doHttp(target string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://%s/internal/pleaseshutdown", target)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %s\n", url, err)
	} else {
		fmt.Printf("Response: %v\n", resp)
	}
}
func doRegistry() {

}

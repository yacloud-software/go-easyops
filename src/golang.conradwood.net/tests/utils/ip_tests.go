package main

import (
	"fmt"
	"os"

	"golang.conradwood.net/go-easyops/utils"
)

func TestIPs() {
	check_ip_private("192.168.1.1:5000", true)
	check_ip_private("192.168.1.0/24", true)
	check_ip_private("8.8.8.8/32", false)
	check_ip_private("8.8.8.8/24", false)
	check_ip_private("fe80::9e6b:ff:fe10:52c5/64", true)
	check_ip_private("::1/128", true)
	check_ip_private("2a00:1450:4009:823::2004", false)
	check_ip_private("2a00:1450:4009:823::2004/64", false)
	check_ip_private("[2a00:1450:4009:823::2004]:5000", false)
}
func check_ip_private(ip string, expected bool) {
	fmt.Printf("Checking ip \"%s\"\n", ip)
	b, err := utils.IsPrivateIP(ip)
	utils.Bail(fmt.Sprintf("could not parse \"%s\"", ip), err)
	if b != expected {
		fmt.Printf("Expected private==%v for ip \"%s\" but result is %v", expected, ip, b)
		os.Exit(10)
	}
}

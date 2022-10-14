package linux

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	pref_nic_if = []string{"dummy", "wg", "tun"}
)

func fail(txt string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("ERROR: %s: %s\n", txt, err)
	os.Exit(10)
}

// find my ip (which is not localhost)
func (*linux) MyIP() string {
	ifaces, err := net.Interfaces()
	fail("could not get interfaces", err)
	var cur_ip *net.IPNet
	cur_pref := 0

	for _, iface := range ifaces {
		fmt.Printf("Iface %s:\n", iface.Name)
		addrs, err := iface.Addrs()
		fail("cannot get interface address", err)
		for _, adr := range addrs {
			use := true
			ipnet, ok := adr.(*net.IPNet)
			if !ok {
				continue
			}
			s := ""
			if ipnet.IP.IsLoopback() {
				use = false
				s = " [loopback]"
			}
			if ipnet.IP.To4() == nil {
				use = false
				s = " [not ip4]"
			}
			fmt.Printf("   %s%s\n", adr.String(), s)
			if use {
				if cur_ip == nil || nic_name_pref(iface.Name) > cur_pref {
					cur_ip = ipnet
				}
			}
		}
	}
	if cur_ip != nil {
		return cur_ip.IP.String()
	}
	return ""
}

func nic_name_pref(name string) int {
	for i, n := range pref_nic_if {
		if strings.HasPrefix(name, n) {
			return i + 1
		}
	}
	return len(pref_nic_if) + 2
}

package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
Parse an IPAddress. this may be either IPv4 or IPv6.
The result is the IPaddress (perhaps normalised), the port (if any), the version
or an error (if unparseable)
*/
func ParseIP(ip_s string) (string, uint32, int, error) {
	if len(ip_s) < 4 {
		return "", 0, 0, fmt.Errorf("\"%s\" is not a valid ip", ip_s)
	}
	ct := strings.Count(ip_s, ":")
	if ct == 1 {
		idx := strings.Index(ip_s, ":")
		ip := ip_s[:idx]
		port, err := strconv.Atoi(ip_s[idx+1:])
		if err != nil {
			return "", 0, 0, err
		}
		return ip, uint32(port), 4, nil
	}

	if ct == 0 {
		return ip_s, 0, 4, nil
	}

	// must be ip6:
	if ip_s[0] != '[' {
		// without port
		return ip_s, 0, 6, nil
	}

	ep := ip_s[1:]
	idx := strings.Index(ep, "]")
	if idx == -1 {
		return "", 0, 0, fmt.Errorf("not a valid ipv6 with port: \"%s\"", ip_s)
	}
	ip := ep[:idx]
	port, err := strconv.Atoi(ep[idx+2:]) // skip "]" and ":"
	if err != nil {
		return "", 0, 0, err
	}
	return ip, uint32(port), 6, nil
}

// error if not a valid ip
// true if it is a non-routeable IP, such as link-local, loopback, or rfc 1918
func IsPrivateIP(ip_s string) (bool, error) {
	ip, masksize, err := GetIPAndNet(ip_s)
	if err != nil {
		return false, err
	}
	ipa, _, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, masksize))
	if err != nil {
		return false, err
	}

	if ipa.IsPrivate() {
		return true, nil
	}
	if ipa.IsLoopback() {
		return true, nil
	}
	if ipa.IsLinkLocalUnicast() {
		return true, nil
	}
	if ipa.IsLinkLocalMulticast() {
		return true, nil
	}
	if ipa.IsInterfaceLocalMulticast() {
		return true, nil
	}
	return false, nil

}

/*
splits things like so:

		172.29.1.0/24  => "172.29.1.0" and 24
		172.29.1.5:5000  => "172.29.1.5" and 32
	        2a01:4b00:ab0f:5100:5::5/64 => "2a01:4b00:ab0f:5100:5::5" and 64
*/
func GetIPAndNet(ip_s string) (string, int, error) {
	if strings.Contains(ip_s, "/") {
		// with CIDR
		ip, ipnet, err := net.ParseCIDR(ip_s)
		if err != nil {
			return "", 0, err
		}
		ones, _ := ipnet.Mask.Size()
		return ip.String(), ones, nil
	}
	// no mask (but possibly with port)
	ip_ns := ""
	host, _, err := net.SplitHostPort(ip_s)
	if err == nil {
		ip_ns = host
	} else {
		ip_ns = ip_s
	}
	ip, _, t, err := ParseIP(ip_ns)
	if err != nil {
		return "", 0, err
	}
	if t == 4 {
		return ip, 32, nil
	} else if t == 6 {
		return ip, 128, nil
	}
	return "", 0, fmt.Errorf("invalid ip type %d for \"%s\"", t, ip_s)
}

// return true if this is a IPv6 link local ip
func IsLinkLocal(s string) bool {
	ips, _, version, err := ParseIP(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ip %s", s))
	}
	if version != 6 {
		return false
	}
	if strings.HasPrefix(ips, "fe80:") {
		return true
	}
	return false
}

// returns true if this is an IPv4 or IPv6 loopback address
func IsLoopback(s string) bool {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "127.") {
		return true
	}
	if s == "::1" {
		return true
	}
	return false
}

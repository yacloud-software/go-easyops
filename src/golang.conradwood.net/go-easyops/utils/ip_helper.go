package utils

import (
	"fmt"
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

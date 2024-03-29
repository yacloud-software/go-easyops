package client

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/standalone"
	"net"
	"strings"
)

const (
	DIRECT_PREFIX = "direct://"
	PROXY_PREFIX  = "proxy://"
)

// this is called by grpc to get a connection
func CustomDialer(ctx context.Context, name string) (net.Conn, error) {
	if cmdline.IsStandalone() {
		return standalone.DialService(ctx, name)
	}

	if *dialer_debug {
		fmt.Printf("[go-easyops] custom dialling \"%s\"\n", name)
	}
	t := name
	if strings.HasPrefix(name, PROXY_PREFIX) {
		sid := t[len(PROXY_PREFIX):]
		pt, err := GetProxyTarget(ctx, sid)
		if err != nil {
			return nil, err
		}
		if pt == nil {
			return nil, fmt.Errorf("no such proxy service: \"%s\"", name)
		}
		return pt.tcpConn, nil
	}
	if strings.HasPrefix(name, DIRECT_PREFIX) {
		t = t[len(DIRECT_PREFIX):]
	}
	if *dialer_debug {
		fmt.Printf("Dialing: %s (%s)\n", name, t)
	}
	return (&net.Dialer{}).DialContext(ctx, "tcp", t)
}

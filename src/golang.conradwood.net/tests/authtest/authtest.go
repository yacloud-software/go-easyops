package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	_ "golang.conradwood.net/go-easyops/http"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"time"
)

var (
	loop          = flag.Bool("loop", false, "continously ping")
	ctr           = 0
	port          = flag.Int("port", 12322, "grpc port")
	role          = flag.String("role", "", "role of this binary (simple|authproxy|client|service|simpleping|load|register|myip")
	servicetokens []string
)

func main() {
	flag.Parse()
	b, err := utils.ReadFile("/etc/cnw/servicetokens")
	utils.Bail("failed to read tokens", err)
	for _, s := range strings.Split(string(b), "\n") {
		if len(s) > 5 {
			servicetokens = append(servicetokens, s)
		}
	}
	if *role == "register" {
		Register()
	} else if *role == "myip" {
		MyIP()
	} else if *role == "authproxy" {
		StartAuthProxy()
	} else if *role == "client" {
		StartClient()
	} else if *role == "load" {
		LoadClient()
	} else if *role == "simple" {
		StartSimple()
	} else if *role == "simpleping" {
		SimplePing()
	} else if *role == "service" {
		tokens.DisableUserToken() // make sure we don't pick up users' stuff
		//	p := utils.RandomInt(100)
		//		ge.GetEasyOpsTestClient(&testserver{}, *port+p)
		panic("cannot do server anymore")
	} else {
		fmt.Printf("Invalid -role specified: %s\n", *role)
		os.Exit(10)
	}
	fmt.Printf("Done (%s)\n", *role)
}

type testserver struct{}

func (t *testserver) SimplePing(ctx context.Context, req *common.Void) (*common.Void, error) {
	fmt.Printf("%s Simpleping received\n", utils.TimeString(time.Now()))
	return &common.Void{}, nil
}
func (t *testserver) Ping(ctx context.Context, req *ge.Chain) (*ge.Chain, error) {
	var err error
	s := auth.GetService(ctx)
	u := auth.GetUser(ctx)
	sn := "[NONE]"
	if s != nil {
		sn = s.ID
	}
	reqid := ""
	ctr++
	fmt.Printf("%s %04d Pinged by user %s\n", utils.TimeString(time.Now()), ctr, u)
	if u == nil {
		return req, nil
	}
	call := &ge.Call{
		RequestID: reqid,
		Position:  req.Position,
		UserID:    u.ID,
		ServiceID: sn}
	req.Calls = append(req.Calls, call)
	//fmt.Printf("Adding: %#v\n", call)
	if len(req.Calls) < 10 {
		req.Position = req.Position + 1
		svtok := servicetokens[int(req.Position)%len(servicetokens)]
		tokens.SetServiceTokenParameter(svtok)
		req, err = ge.GetEasyOpsTestClient().Ping(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

// use a wellknown service to do a simple call
func StartSimple() {
	ctx := authremote.Context()
	lreq := &registry.V2ListRequest{}
	resp, err := client.GetRegistryClient().ListRegistrations(ctx, lreq)
	utils.Bail("Failed to query registry", err)
	for _, gr := range resp.Registrations {
		fmt.Printf("Result: %s (%s)\n", gr.Target.ServiceName, gr.Target.RoutingInfo)
		fmt.Printf("   %s:%d\n", gr.Target.IP, gr.Target.Port)
	}
	os.Exit(0)
}

func MyIP() {
	l := linux.New()
	s := l.MyIP()
	fmt.Printf("My IP: \"%s\"\n", s)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/create"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

var (
	loop          = flag.Bool("loop", false, "continously ping")
	ctr           = 0
	port          = flag.Int("port", 12322, "grpc port")
	role          = flag.String("role", "", "role of this binary (simple|authproxy|client|service|simpleping|load")
	servicetokens = []string{
		"AQbqyBooiPnRazNPKkmAXsRwTACKuiwPFKKbVkkQIYeWoEOsYPKKWXSnJhniMQI",
		"WEVPXnZYPzlcNcIvZrQBXxEjNOHNwYZDknUDNKKsjcXyywwsZovXqbWLSdrouZwH",
		"aDwbuIDniulJjBpzJqobrLqKFTpgbhUseuMhScZWGOQmOwxPXufAVHeTqVhyzGqp",
	}
)

func main() {
	flag.Parse()
	if *role == "authproxy" {
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
		p := utils.RandomInt(100)
		create.NewEasyOpsTest(&testserver{}, *port+p)
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
	cs := rpc.CallStateFromContext(ctx)
	if cs != nil && cs.RPCIResponse != nil {
		reqid = cs.RPCIResponse.RequestID
	}
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
		req, err = create.NewEasyOpsTestClient().Ping(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

// use a wellknown service to do a simple call
func StartSimple() {
	ctx := tokens.ContextWithToken()
	lreq := &registry.ListRequest{}
	resp, err := create.GetRegistryClient().ListServices(ctx, lreq)
	utils.Bail("Failed to query registry", err)
	for _, gr := range resp.Service {
		fmt.Printf("Result: %#s (%s)\n", gr.Service.Name, gr.Service.Gurupath)
		for _, a := range gr.Location.Address {
			fmt.Printf("   %s:%d\n", a.Host, a.Port)
		}
	}
	os.Exit(0)
}

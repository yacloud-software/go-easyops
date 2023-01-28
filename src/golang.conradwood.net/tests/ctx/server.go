package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	gs "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
)

func start_server() {
	sd := server.NewServerDef()
	sd.Port = 3005
	sd.SetOnStartupCallback(run_tests)
	sd.Register = server.Register(
		func(g *grpc.Server) error {
			gs.RegisterCtxTestServer(g, &geServer{})
			return nil
		},
	)
	err := server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
}

type geServer struct{}

func (g *geServer) TestFork(ctx context.Context, req *common.Void) (*common.Void, error) {
	u1 := auth.GetUser(ctx)
	nctx, err := auth.ForkContext(ctx)
	if err != nil {
		return nil, err
	}
	u2 := auth.GetUser(nctx)
	if (u1 == nil && u2 != nil) || (u2 == nil && u1 != nil) || (u1.ID != u2.ID) {
		return nil, fmt.Errorf("u1 (%s) != u2 (%s)", auth.UserIDString(u1), auth.UserIDString(u2))
	}
	return &common.Void{}, nil
}

func run_tests() {
	fmt.Printf("Running tests...\n")

	fmt.Printf("Running fork test...\n")
	ctx := authremote.Context()
	_, err := gs.GetCtxTestClient().TestFork(ctx, &common.Void{})
	utils.Bail("failed fork test\n", err)

	fmt.Printf("Running fork test...\n")
	cmdline.SetContextWithBuilder(true)
	ctx = authremote.Context()
	_, err = gs.GetCtxTestClient().TestFork(ctx, &common.Void{})
	utils.Bail("failed fork test\n", err)

	fmt.Printf("Done\n")
	os.Exit(0)
}

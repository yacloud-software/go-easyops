package main

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	gs "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	gcm "golang.conradwood.net/go-easyops/common"
	pctx "golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
	"time"
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
	/*
		u1 := auth.GetUser(ctx)
		fmt.Printf("[testfork] invoked as user %s\n", auth.UserIDString(u1))
	*/
	nctx, err := auth.ForkContext(ctx)
	if err != nil {
		return nil, err
	}
	err = AssertEqualContexts(ctx, nctx)
	if err != nil {
		return nil, err
	}

	return &common.Void{}, nil
}

func (g *geServer) TestDeSer(ctx context.Context, req *common.Void) (*gs.SerialisedContext, error) {
	//	fmt.Printf("b\n")
	b, err := auth.SerialiseContext(ctx)
	if err != nil {
		return nil, err
	}

	s, err := auth.SerialiseContextToString(ctx)
	if err != nil {
		return nil, err
	}

	ictx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, b)
	if err != nil {
		return nil, err
	}

	err = AssertEqualContexts(ctx, ictx)
	if err != nil {
		return nil, err
	}

	if cmdline.ContextWithBuilder() {
		ictx, err = pctx.DeserialiseContextFromString(s)
		if err != nil {
			return nil, err
		}

		err = AssertEqualContexts(ctx, ictx)
		if err != nil {
			return nil, err
		}
		if !pctx.IsSerialisedByBuilder(b) {
			return nil, fmt.Errorf("ctx does not recognise its own (byte) serialised context")
		}
		if !pctx.IsSerialisedByBuilder([]byte(s)) {
			return nil, fmt.Errorf("ctx does not recognise its own (byte) serialised context")
		}
	}

	res := &gs.SerialisedContext{
		Data:    b,
		SData:   s,
		User:    auth.GetSignedUser(ctx),
		Service: auth.GetSignedService(ctx),
	}
	return res, nil
}
func run_tests() {
	fmt.Printf("Running tests...\n")

	cmdline.SetContextWithBuilder(false)
	t := NewTest("fork test")
	ctx := authremote.Context()
	_, err := gs.GetCtxTestClient().TestFork(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(true)
	t = NewTest("fork test")
	ctx = authremote.Context()
	_, err = gs.GetCtxTestClient().TestFork(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(false)
	t = NewTest("(de)serialise")
	ctx = authremote.Context()
	_, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(true)
	t = NewTest("(de)serialise")
	ctx = authremote.Context()
	dctx, err := gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(true)
	t = NewTest("use serialised context to access service")
	if !pctx.IsSerialisedByBuilder(dctx.Data) {
		t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
	}
	ctx, err = pctx.DeserialiseContext(dctx.Data)
	t.Error(err)
	dctx, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	if dctx == nil || dctx.User == nil {
		t.Error(fmt.Errorf("No user in service with serialised context"))
	}
	t.Done()

	cmdline.SetContextWithBuilder(false)
	t = NewTest("use serialised context to access service")
	if !pctx.IsSerialisedByBuilder(dctx.Data) {
		t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
	}
	ctx, err = pctx.DeserialiseContext(dctx.Data)
	t.Error(err)
	dctx, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	if dctx == nil || dctx.User == nil {
		t.Error(fmt.Errorf("No user in service with serialised context"))
	}
	t.Done()

	cmdline.SetContextWithBuilder(false)
	t = NewTest("serialise old, deserialise new")
	ctx = authremote.Context()
	r, err := gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	if err == nil {
		cmdline.SetContextWithBuilder(true)
		nctx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, r.Data)
		t.Error(err)
		if err != nil {
			err = AssertEqualContexts(ctx, nctx)
			t.Error(err)
		}
	}
	t.Done()

	//	fmt.Printf("a\n")
	cmdline.SetContextWithBuilder(false)
	t = NewTest("old context, call new service")
	ctx = authremote.Context()
	cmdline.SetContextWithBuilder(true)
	//	fmt.Printf("c\n")
	_, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	fmt.Printf("Done\n")
	PrintResult()
	os.Exit(0)
}

func AssertEqualContexts(ctx1, ctx2 context.Context) error {
	su1 := auth.GetSignedUser(ctx1)
	su2 := auth.GetSignedUser(ctx2)
	if !CompareUsers(su1, su2) {
		u1 := gcm.VerifySignedUser(su1)
		u2 := gcm.VerifySignedUser(su2)
		return fmt.Errorf("u1 (%s) != u2 (%s)", auth.UserIDString(u1), auth.UserIDString(u2))
	}

	return nil
}

func CompareUsers(su1, su2 *apb.SignedUser) bool {
	u1 := gcm.VerifySignedUser(su1)
	u2 := gcm.VerifySignedUser(su2)
	if u1 == nil && u2 == nil {
		return true
	}
	if u1 == nil && u2 != nil {
		return false
	}
	if u1 != nil && u2 == nil {
		return false
	}

	if u1.ID != u2.ID {
		return false
	}
	return true
}

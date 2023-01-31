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
	b, err := auth.SerialiseContext(ctx)
	if err != nil {
		return nil, err
	}
	ictx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, b)
	if err != nil {
		return nil, err
	}
	res := &gs.SerialisedContext{
		Data:    b,
		User:    auth.GetSignedUser(ctx),
		Service: auth.GetSignedService(ctx),
	}
	err = AssertEqualContexts(ctx, ictx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func NewTest(format string, args ...interface{}) *test {
	t := &test{prefix: fmt.Sprintf(format, args...)}
	fmt.Printf("%s -------- STARTING\n", t.Prefix())
	return t
}

type test struct {
	err    error
	prefix string
}

func (t *test) Prefix() string {
	v := fmt.Sprintf("%v", cmdline.ContextWithBuilder())
	return fmt.Sprintf("[%s (builder=%5s)]", t.prefix, v)
}

func (t *test) Printf(format string, args ...interface{}) {
	fmt.Printf(t.Prefix()+" "+format, args...)
}
func (t *test) Error(err error) {
	if err == nil {
		return
	}
	t.err = err
	fmt.Printf("%s Failed (%s)\n", t.Prefix(), err)
}
func (t *test) Done() {
	if t.err != nil {
		fmt.Printf("%s -------- FAILURE\n", t.Prefix())
		return
	}
	fmt.Printf("%s -------- SUCCESS\n", t.Prefix())
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
	_, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	t.Done()

	fmt.Printf("Done\n")
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
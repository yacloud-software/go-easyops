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
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"io"
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
func (g *geServer) TestStream(req *common.Void, srv gs.CtxTest_TestStreamServer) error {
	ctx := srv.Context()
	u := auth.GetUser(ctx)
	if u == nil {
		return errors.Unauthenticated(ctx, "no user")
	}
	_, err := g.TestDeSer(ctx, req)
	return err
}
func (g *geServer) TestDeSer(ctx context.Context, req *common.Void) (*gs.SerialisedContext, error) {
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "no user")
	}
	m := map[string]string{"foo": "bar"}
	ictx := authremote.DerivedContextWithRouting(ctx, m, true)
	u = auth.GetUser(ictx)
	if u == nil {
		return nil, fmt.Errorf("DerivedContextWithRouting lost user information!")
	}
	err := AssertEqualContexts(ctx, ictx)
	if err != nil {
		return nil, err
	}

	//	fmt.Printf("b\n")
	b, err := auth.SerialiseContext(ctx)
	if err != nil {
		return nil, err
	}

	s, err := auth.SerialiseContextToString(ctx)
	if err != nil {
		return nil, err
	}

	ictx, err = auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, b)
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
		_, err = pctx.DeserialiseContext(b)
		if err != nil {
			return nil, err
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
	cmdline.SetDatacenter(false)
	run_all_tests()
	cmdline.SetDatacenter(true)
	run_all_tests()
	fmt.Printf("Done\n")
	PrintResult()
	os.Exit(0)

}
func run_all_tests() {
	fmt.Printf("Running tests...\n")

	cmdline.SetContextWithBuilder(false)
	t := NewTest("stream test")
	ctx := authremote.Context()
	srv, err := gs.GetCtxTestClient().TestStream(ctx, &common.Void{})
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextWithBuilder(true)
	t = NewTest("stream test")
	ctx = authremote.Context()
	srv, err = gs.GetCtxTestClient().TestStream(ctx, &common.Void{})
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	checkStream("CallUnaryFromStream", func(ctx context.Context) (recv, error) {
		return gs.GetCtxTestClient().CallUnaryFromStream(ctx, &common.Void{})
	})

	checkStream("CallStreamFromStream", func(ctx context.Context) (recv, error) {
		return gs.GetCtxTestClient().CallStreamFromStream(ctx, &common.Void{})
	})
	checkUnary("CallStreamFromUnary", func(ctx context.Context) error {
		_, err := gs.GetCtxTestClient().CallStreamFromUnary(ctx, &common.Void{})
		return err
	})
	checkUnary("CallUnaryFromUnary", func(ctx context.Context) error {
		_, err := gs.GetCtxTestClient().CallUnaryFromUnary(ctx, &common.Void{})
		return err
	})

	cmdline.SetContextWithBuilder(false)
	t = NewTest("fork test")
	ctx = authremote.Context()
	_, err = gs.GetCtxTestClient().TestFork(ctx, &common.Void{})
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

	if dctx != nil {
		cmdline.SetContextWithBuilder(true)
		t = NewTest("use serialised context to access service")
		if !pctx.IsSerialisedByBuilder(dctx.Data) {
			t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
		}
		ctx, err = pctx.DeserialiseContext(dctx.Data)
		t.Error(err)
	}
	dctx, err = gs.GetCtxTestClient().TestDeSer(ctx, &common.Void{})
	t.Error(err)
	if dctx == nil || dctx.User == nil {
		t.Error(fmt.Errorf("No user in service with serialised context"))
	}
	t.Done()

	if dctx != nil {
		cmdline.SetContextWithBuilder(false)
		t = NewTest("use serialised context to access service")
		if !pctx.IsSerialisedByBuilder(dctx.Data) {
			t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
		}
		ctx, err = pctx.DeserialiseContext(dctx.Data)
		t.Error(err)
	}
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

type recv interface {
	Recv() (*common.Void, error)
}

func checkSrv(r recv) error {
	for {
		_, err := r.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func checkStream(name string, f func(ctx context.Context) (recv, error)) {
	cmdline.SetContextWithBuilder(true)
	t := NewTest("stream-bouncer %s test", name)
	ctx := authremote.Context()
	srv, err := f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextWithBuilder(false)
	t = NewTest("stream-bouncer %s test", name)
	ctx = authremote.Context()
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextWithBuilder(true)
	ctx = authremote.Context()
	cmdline.SetContextWithBuilder(false)
	t = NewTest("stream-bouncer %s test", name)
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextWithBuilder(false)
	ctx = authremote.Context()
	cmdline.SetContextWithBuilder(true)
	t = NewTest("stream-bouncer %s test", name)
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

}
func checkUnary(name string, f func(ctx context.Context) error) {
	cmdline.SetContextWithBuilder(true)
	t := NewTest("unary-bouncer %s test", name)
	ctx := authremote.Context()
	err := f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(false)
	t = NewTest("unary-bouncer %s test", name)
	ctx = authremote.Context()
	err = f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(true)
	ctx = authremote.Context()
	cmdline.SetContextWithBuilder(false)
	t = NewTest("unary-bouncer %s test", name)
	err = f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextWithBuilder(false)
	ctx = authremote.Context()
	cmdline.SetContextWithBuilder(true)
	t = NewTest("unary-bouncer %s test", name)
	err = f(ctx)
	t.Error(err)
	t.Done()

}

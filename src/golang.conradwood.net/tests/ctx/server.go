package main

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	gcm "golang.conradwood.net/go-easyops/common"
	pctx "golang.conradwood.net/go-easyops/ctx"
	//"golang.conradwood.net/go-easyops/errors"
	"flag"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"io"
	"os"
	"sync"
	"time"
)

var (
	run_sleep_tests = flag.Bool("run_sleep_tests", true, "extra tests to test if timeouts and cancels propagate accurately, they take a while...")
	runlock         sync.Mutex
	didrun          = false
)

func start_server() {
	sd := server.NewServerDef()
	sd.SetPort(3005)
	sd.SetOnStartupCallback(run_tests)
	sd.SetRegister(server.Register(
		func(g *grpc.Server) error {
			ge.RegisterCtxTestServer(g, &geServer{})
			return nil
		},
	))
	err := server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
}

type geServer struct{}

func (g *geServer) Sleep(ctx context.Context, req *ge.SleepRequest) (*common.Void, error) {
	t := time.Duration(req.Seconds) * time.Second
	fmt.Printf("Sleeping for %0.2f seconds\n", t.Seconds())
	time.Sleep(t)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return &common.Void{}, nil
}

func (g *geServer) TestFork(ctx context.Context, req *ge.RequiredContext) (*common.Void, error) {
	err := AssertRequiredContext(ctx, req)
	if err != nil {
		return nil, err
	}

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
func (g *geServer) TestUnary(ctx context.Context, req *ge.RequiredContext) (*common.Void, error) {
	err := AssertRequiredContext(ctx, req)
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil

}
func (g *geServer) TestStream(req *ge.RequiredContext, srv ge.CtxTest_TestStreamServer) error {
	ctx := srv.Context()
	err := AssertRequiredContext(ctx, req)
	if err != nil {
		return err
	}

	_, err = g.TestDeSer(ctx, req)
	return err
}
func (g *geServer) TestDeSer(ctx context.Context, req *ge.RequiredContext) (*ge.SerialisedContext, error) {
	err := AssertRequiredContext(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("TestDeSer: required (and received) context calling service: %s\n", auth.UserIDString(auth.GetService(ctx)))
	m := map[string]string{"foo": "bar"}
	ictx := authremote.DerivedContextWithRouting(ctx, m, true)
	err = AssertRequiredContext(ictx, req)
	if err != nil {
		return nil, fmt.Errorf("broken derived context (%s)", err)
	}

	err = AssertEqualContexts(ctx, ictx)
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
	err = AssertRequiredContext(ictx, req)
	if err != nil {
		return nil, err
	}

	err = AssertEqualContexts(ctx, ictx)
	if err != nil {
		return nil, err
	}

	sctx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, []byte(s))
	if err != nil {
		return nil, err
	}
	err = AssertRequiredContext(sctx, req)
	if err != nil {
		return nil, err
	}

	err = AssertEqualContexts(ctx, sctx)
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

	res := &ge.SerialisedContext{
		Data:    b,
		SData:   s,
		User:    auth.GetSignedUser(ctx),
		Service: auth.GetSignedService(ctx),
	}
	return res, nil
}
func run_tests() {
	runlock.Lock()
	if didrun == true {
		runlock.Unlock()
		return
	}
	didrun = true
	runlock.Unlock()
	// first run a couple of very quick tests..
	svc := authremote.GetLocalServiceAccount()
	ctx := authremote.Context()
	if ctx == nil {
		panic("no context")
	}
	if svc != nil && auth.GetSignedService(ctx) == nil {
		panic("missing service")
	}
	ctx = authremote.DerivedContextWithRouting(ctx, make(map[string]string), true)
	if ctx == nil {
		panic("no context")
	}
	if svc != nil && auth.GetSignedService(ctx) == nil {
		panic("missing service")
	}

	cmdline.SetDatacenter(false)
	run_all_tests()
	cmdline.SetDatacenter(true)
	run_all_tests()
	if *run_sleep_tests {
		sleepTests()
	}
	fmt.Printf("Done\n")
	PrintResult()
	os.Exit(0)

}
func run_all_tests() {
	fmt.Printf("Running tests...\n")

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t := NewTest("simple unary test")
	ctx := authremote.Context()
	_, err := ge.GetCtxTestClient().TestUnary(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t = NewTest("simple unary test")
	ctx = authremote.Context()
	_, err = ge.GetCtxTestClient().TestUnary(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("stream test")
	ctx = authremote.Context()
	srv, err := ge.GetCtxTestClient().TestStream(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t = NewTest("stream test")
	ctx = authremote.Context()
	srv, err = ge.GetCtxTestClient().TestStream(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	checkStream("CallUnaryFromStream", func(ctx context.Context) (recv, error) {
		return ge.GetCtxTestClient().CallUnaryFromStream(ctx, CreateContextObject(ctx))
	})

	checkStream("CallStreamFromStream", func(ctx context.Context) (recv, error) {
		return ge.GetCtxTestClient().CallStreamFromStream(ctx, CreateContextObject(ctx))
	})
	checkUnary("CallStreamFromUnary", func(ctx context.Context) error {
		_, err := ge.GetCtxTestClient().CallStreamFromUnary(ctx, CreateContextObject(ctx))
		return err
	})
	checkUnary("CallUnaryFromUnary", func(ctx context.Context) error {
		_, err := ge.GetCtxTestClient().CallUnaryFromUnary(ctx, CreateContextObject(ctx))
		return err
	})

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("fork test")
	ctx = authremote.Context()
	_, err = ge.GetCtxTestClient().TestFork(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t = NewTest("fork test")
	ctx = authremote.Context()
	_, err = ge.GetCtxTestClient().TestFork(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("(de)serialise")
	ctx = authremote.Context()
	_, err = ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t = NewTest("(de)serialise")
	ctx = authremote.Context()
	dctx, err := ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
	t.Error(err)
	t.Done()

	if dctx != nil {
		cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
		t = NewTest("use serialised context to access service")
		if !pctx.IsSerialisedByBuilder(dctx.Data) {
			t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
		}
		ctx, err = pctx.DeserialiseContext(dctx.Data)
		t.Error(err)
	}
	dctx, err = ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
	t.Error(err)
	if dctx == nil || (!CompareUsers(dctx.User, auth.GetSignedUser(ctx))) {
		t.Error(fmt.Errorf("No user in service with serialised context"))
	}
	t.Done()

	if dctx != nil {
		cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
		t = NewTest("use serialised context to access service")
		if !pctx.IsSerialisedByBuilder(dctx.Data) {
			t.Error(fmt.Errorf("ctx failed to recognise it as a context"))
		}
		ctx, err = pctx.DeserialiseContext(dctx.Data)
		t.Error(err)
	}
	dctx, err = ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
	t.Error(err)
	if dctx == nil || (!CompareUsers(dctx.User, auth.GetSignedUser(ctx))) {
		t.Error(fmt.Errorf("No user in service with serialised context"))
	}
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("serialise old, deserialise new")
	ctx = authremote.Context()
	r, err := ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
	t.Error(err)
	if err == nil {
		cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
		nctx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, r.Data)
		t.Error(err)
		if err != nil {
			err = AssertEqualContexts(ctx, nctx)
			t.Error(err)
		}
	}
	t.Done()

	//	fmt.Printf("a\n")
	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("old context, call new service")
	ctx = authremote.Context()
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	//	fmt.Printf("c\n")
	_, err = ge.GetCtxTestClient().TestDeSer(ctx, CreateContextObject(ctx))
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
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t := NewTest("stream-bouncer %s (new ctx)", name)
	ctx := authremote.Context()
	srv, err := f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("stream-bouncer %s (old ctx)", name)
	ctx = authremote.Context()
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	ctx = authremote.Context()
	t = NewTest("stream-bouncer %s (new/old ctx)", name)
	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	ctx = authremote.Context()
	t = NewTest("stream-bouncer %s (old/new ctx)", name)
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	srv, err = f(ctx)
	t.Error(err)
	t.Error(checkSrv(srv))
	t.Done()

}
func checkUnary(name string, f func(ctx context.Context) error) {
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t := NewTest("unary-bouncer with_cb %s (new ctx)", name)
	ctx := authremote.Context()
	err := f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("unary-bouncer wo_cb %s (old ctx)", name)
	ctx = authremote.Context()
	err = f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	t = NewTest("unary-bouncer %s cb_tf (new/old ctx)", name)
	ctx = authremote.Context()
	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	err = f(ctx)
	t.Error(err)
	t.Done()

	cmdline.SetContextBuilderVersion(OLD_CONTEXT_VERSION)
	t = NewTest("unary-bouncer %s cb_ft test (old/new ctx)", name)
	ctx = authremote.Context()
	cmdline.SetContextBuilderVersion(NEW_CONTEXT_VERSION)
	err = f(ctx)
	t.Error(err)
	t.Done()

}

func CreateContextObject(ctx context.Context) *ge.RequiredContext {
	res := &ge.RequiredContext{
		User:    auth.GetSignedUser(ctx),
		Service: auth.GetSignedService(ctx),
	}
	return res
}
func AssertRequiredContext(ctx context.Context, rc *ge.RequiredContext) error {
	if ctx == nil {
		return fmt.Errorf("no context to assert")
	}
	u := auth.GetSignedUser(ctx)
	s := auth.GetSignedService(ctx)
	err := AssertSameUser("user", rc.User, u)
	if err != nil {
		fmt.Println("Mismatched context:" + pctx.Context2String(ctx))
		return err
	}
	if rc.Service != nil {
		err = AssertSameUser("service", rc.Service, s)
		if err != nil {
			fmt.Println("Mismatched context:" + pctx.Context2String(ctx))
			return err
		}
	}
	return nil
}
func AssertSameUser(s string, u1, u2 *apb.SignedUser) error {
	if !CompareUsers(u1, u2) {
		uu1 := gcm.VerifySignedUser(u1)
		uu2 := gcm.VerifySignedUser(u2)
		utils.PrintStack("%s Mismatch: expected=%s, actual=%s", s, auth.UserIDString(uu1), auth.UserIDString(uu2))
		return fmt.Errorf("%s Mismatch: expected=%s, actual=%s", s, auth.UserIDString(uu1), auth.UserIDString(uu2))
	}
	return nil
}

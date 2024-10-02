package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.conradwood.net/apis/common"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/ctx"
	gctx "golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/ctx/shared"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
)

const (
	ACTION_CALL_YOURSELF      = 1
	ACTION_SEND_ACCESS_DENIED = 2
)

var (
	TEST_SERVICE_IDS = []string{"33", "31", "22", "29", "27", "35", "57", "39"}
)

func main() {
	flag.Parse()
	server.SetHealth(common.Health_READY)
	sd := server.NewServerDef()
	sd.SetPort(3006)
	sd.SetOnStartupCallback(run_tests)
	sd.SetRegister(server.Register(
		func(g *grpc.Server) error {
			ge.RegisterCtx2TestServer(g, &geServer{})
			return nil
		},
	))
	err := server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
}
func run_tests() {
	fmt.Printf("Starting tests...\n")
	ctx := authremote.Context()
	//trr := &ge.TriggerRPCRequest{Action: ACTION_SEND_ACCESS_DENIED}
	trr := &ge.TriggerRPCRequest{Action: ACTION_CALL_YOURSELF}
	_, err := ge.GetCtx2TestClient().TriggerRPC(ctx, trr)
	utils.Bail("test failed", err)
	fmt.Printf("Tests completed\n")
}

type geServer struct {
}

func (ges *geServer) TriggerRPC(ctx context.Context, req *ge.TriggerRPCRequest) (*common.Void, error) {
	cmdline.SetDebugContext()
	fmt.Printf("------------------------------- In TRIGGERRPC %d --------------------\n", req.Counter)
	fmt.Printf("LocalState: %s\n", shared.LocalState2string(gctx.GetLocalState(ctx)))
	if req.Action == ACTION_SEND_ACCESS_DENIED {
		return nil, errors.AccessDenied(ctx, "told to return access denied")
	}
	if req.Action == ACTION_CALL_YOURSELF {
		if req.Counter == 0 {
			ctx = buildContext(req.Counter)
		}
		assert_correct_ctx(ctx, req)
		if req.Counter < 5 {
			req.Counter++
			fmt.Printf("------------------------------- Calling TRIGGERRPC %d --------------------\n", req.Counter)
			return ge.GetCtx2TestClient().TriggerRPC(ctx, req)
		}
	} else {
		return nil, errors.Errorf("Invalid action \"%d\"", req.Action)
	}
	return &common.Void{}, nil
}
func assert_correct_ctx(ctx context.Context, req *ge.TriggerRPCRequest) {
	ls := gctx.GetLocalState(ctx)
	if ls == nil {
		fail(ctx, req, "No localstate at all")
	}
	svc := ls.CallingService()
	if svc == nil {
		fail(ctx, req, "no calling service")
	}
	svc = ls.CreatorService()
	if svc == nil {
		fail(ctx, req, "no creatorservice")
	}
	if ls.SudoUser() == nil {
		fail(ctx, req, "no sudouser")
	}
	if req.Counter == 1 {
	}
}
func fail(ctx context.Context, req *ge.TriggerRPCRequest, format string, args ...interface{}) {
	fmt.Printf("****** FAILED\n")
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[counter=%d] %s\n", req.Counter, s)
	os.Exit(10)
}

func buildContext(ct uint32) context.Context {
	cb := ctx.NewContextBuilder()
	cb.WithTimeout(time.Duration(3) * time.Second)
	//	cb.WithSession(f.session)
	u, err := authremote.GetSignedUserByID(authremote.Context(), "1")
	utils.Bail("failed to get testuserid", err)
	if u == nil {
		panic("no user to build context with")
	}
	_, s := authremote.GetLocalUsers()
	if s == nil {
		panic("no service user. forgot -token option?")
	}
	sx, err := authremote.GetSignedUserByID(authremote.Context(), TEST_SERVICE_IDS[ct])
	utils.Bail("failed to get user by id", err)
	cb.WithUser(u)
	cb.WithSudoUser(u)
	cb.WithCreatorService(s)
	cb.WithCallingService(sx)
	cb.WithDebug()
	cb.WithTrace()
	cb.WithRequestID("foo-requestid")
	cb.EnableExperiment("debug_context1")
	cb.EnableExperiment("debug_context2")

	res_ctx := cb.ContextWithAutoCancel()
	fmt.Printf("Created new context\n")

	return res_ctx
}

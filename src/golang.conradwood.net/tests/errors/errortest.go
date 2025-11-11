package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 0, "grpc port")
)

func main() {
	flag.Parse()
	fmt.Printf("\n\n\n\n\n\n")
	fmt.Printf("Error testing\n")
	fmt.Printf("\n\n\n\n\n\n")
	ctx := authremote.Context()
	errs := []error{
		errors.AccessDenied(ctx, "foo"),
		errors.Errorf("errorf"),
		fmt.Errorf("fmterror"),
		errors.Wrap(fmt.Errorf("wraperror")),
		errors.Wrap(errors.Wrap(fmt.Errorf("wrap-wrap-error"))),
	}
	for _, err := range errs {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("withstacktrace: %s\n", errors.ErrorStringWithStackTrace(err))
		fmt.Printf("short message: %s\n", errors.ShortMessage(err))
		fmt.Printf("\n\n")
	}
	var err error
	sd := server.NewServerDef()
	p := *port
	if p == 0 {
		p = 4100 + utils.RandomInt(50)
	}
	sd.SetPort(p)
	sd.SetRegister(server.Register(
		func(g *grpc.Server) error {
			pb.RegisterEchoServiceServer(g, &echoServer{})
			return nil
		},
	))
	sd.SetOnStartupCallback(startup)
	err = server.ServerStartup(sd)
	utils.Bail("failed to start server", err)
	fmt.Printf("Done\n")
}

func errstr(err error) string {
	return errors.ErrorStringWithStackTrace(err)
}

func startup() {
	ctx := authremote.Context()
	_, err := pb.GetEchoServiceClient().Ping(ctx, &pb.PingRequest{PleaseFail: true})
	utils.Bail("failed to ping", err)
	time.Sleep(time.Duration(20) * time.Second)
}

type echoServer struct {
}

func (e *echoServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	e.Printf("pinged\n")
	if req.PleaseFail {
		return nil, errors.AccessDenied(ctx, "failed")
		//		return nil, errors.Errorf("failed")
	}
	res := &pb.PingResponse{}
	return res, nil
}

func (e *echoServer) Printf(format string, args ...interface{}) {
	prefix := "[server] "
	x := fmt.Sprintf(format, args...)
	fmt.Print(prefix + x)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/echoservice"
	"time"
	//	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
)

var (
	port    = flag.Int("port", 4106, "grpc port")
	verbose = flag.Bool("verbose", false, "verbose mode")
)

func main() {
	flag.Parse()
	sd := server.NewServerDef()
	sd.Port = *port
	sd.Register = server.Register(
		func(server *grpc.Server) error {
			pb.RegisterEchoServiceServer(server, &echoserver{})
			return nil
		},
	)
	err := server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)

}

type echoserver struct {
}

func (e *echoserver) Ping(ctx context.Context, req *common.Void) (*pb.PingResponse, error) {
	if *verbose {
		fmt.Printf("%s Pinged\n", utils.TimeString(time.Now()))
	}
	res := &pb.PingResponse{Response: "servertext"}
	return res, nil
}

package main

import (
	"context"
	"flag"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/echoservice"
	//	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
)

var (
	port = flag.Int("port", 4106, "grpc port")
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
	res := &pb.PingResponse{Response: "servertext"}
	return res, nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "golang.conradwood.net/apis/getestservice"

	//	"golang.conradwood.net/go-easyops/errors"
	"os"

	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
)

var (
	port    = flag.Int("port", 4106, "grpc port")
	verbose = flag.Bool("verbose", false, "verbose mode")
)

func main() {
	flag.Parse()
	sd := server.NewServerDef()
	sd.SetPort(*port)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			pb.RegisterEchoServiceServer(server, &echoserver{})
			return nil
		},
	))
	fmt.Printf("Sleeping...\n")
	//	time.Sleep(time.Duration(4) * time.Second)
	err := server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)

}

type echoserver struct {
}

func (e *echoserver) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if *verbose {
		fmt.Printf("%s Pinged\n", utils.TimeString(time.Now()))
	}
	res := &pb.PingResponse{Response: req}
	return res, nil
}

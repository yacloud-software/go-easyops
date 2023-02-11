package main

import (
	"context"
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/auth"
	//	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
	"strings"
)

var (
	port          = flag.Int("port", 4106, "The grpc server port")
	ping          = flag.Bool("ping", false, "ping continously")
	ping_once     = flag.Bool("ping_once", false, "ping once")
	tag           = flag.String("tag", "", "key=value tag optional")
	fallback      = flag.Bool("fallback", true, "if true, fallback allowed (for client)")
	inject_errors = flag.Bool("inject_errors", false, "if true inject some errors in the rpc")
	ttl           = flag.Int("ttl", 0, "if >0 the server will ping itself until ttl is 0")
	ctr           = 0
)

// create a simple standard server
type echoServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("GO-EASYOPS Echo test server/client\n")
	if *ping || *ping_once {
		utils.Bail("failed", do_client())
		os.Exit(0)
	}

	sd := server.NewServerDef()

	if *tag != "" {
		kv := strings.SplitN(*tag, "=", 2)
		if len(kv) != 2 {
			fmt.Printf("tags not a key=value line\n")
			os.Exit(10)
		}
		sd.AddTag(kv[0], kv[1])
		fmt.Printf("Added tag \"%s\" with value \"%s\"\n", kv[0], kv[1])
	}

	p := *port
	p = p + utils.RandomInt(50)
	//	sd.AddTag("foo", "bar")
	sd.Port = p
	sd.Register = server.Register(
		func(g *grpc.Server) error {
			pb.RegisterEchoServiceServer(g, &echoServer{})
			return nil
		},
	)
	err := server.ServerStartup(sd)
	//	err := create.NewEchoServiceServer(&echoServer{}, p)
	utils.Bail("Unable to start server", err)
}

func (e *echoServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	u := auth.GetUser(ctx)
	fmt.Printf("    tagserver %s: %d Pinged SEQ=%d, by %s (TTL:%d)\n", printTags(), ctr, req.SequenceNumber, auth.Description(u), req.TTL)
	if req.TTL > 0 {
		req.TTL--
		_, err := pb.GetEchoClient().Ping(ctx, req)
		if err != nil {
			return nil, err
		}

	}
	ctr++
	if *inject_errors {
		i := utils.RandomInt(10)
		if i > 3 {
			return nil, errors.Unavailable(ctx, "Ping()")
		}
	}
	return &pb.PingResponse{ServerTags: parse_tags()}, nil
}

func parse_tags() map[string]string {
	res := make(map[string]string)
	if *tag == "" {
		return res
	}
	kv := strings.SplitN(*tag, "=", 2)
	if len(kv) != 2 {
		fmt.Printf("tags not a key=value line\n")
		os.Exit(10)
	}
	res[kv[0]] = kv[1]
	//	fmt.Printf("Added tag \"%s\" with value \"%s\"\n", kv[0], kv[1])
	return res
}

package main

import (
	"context"
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/getestservice"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
	"strings"
	"time"
)

var (
	port          = flag.Int("port", 4106, "The grpc server port")
	ping          = flag.Bool("ping", false, "ping continously")
	ping_once     = flag.Bool("ping_once", false, "ping once")
	tag           = flag.String("tag", "", "key=value tag optional")
	inject_errors = flag.Bool("inject_errors", false, "if true inject some errors in the rpc")
	ctr           = 0
)

// create a simple standard server
type echoServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("GO-EASYOPS Echo test server/client\n")
	if *ping || *ping_once {
		c := pb.GetEchoClient()
		seq := uint32(0)
		for {
			now := time.Now()
			ctx := clientContext()
			u := auth.GetUser(ctx)
			fmt.Printf("   pinging as %s\n", auth.Description(u))
			seq++
			_, err := c.Ping(ctx, &pb.PingRequest{SequenceNumber: seq, TTL: 1})
			if err != nil {
				fmt.Printf("Error :%s\n", utils.ErrorString(err))
			}
			dur := time.Since(now).Milliseconds()
			fmt.Printf("%d Pinged (%d milliseconds)\n", ctr, dur)
			ctr++
			if !*ping {
				return
			}
			time.Sleep(time.Duration(300) * time.Millisecond)
		}
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
	fmt.Printf("    %d Pinged SEQ=%d, by %s (TTL:%d)\n", ctr, req.SequenceNumber, auth.Description(u), req.TTL)
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
	return &pb.PingResponse{}, nil
}

func clientContext() context.Context {
	if *tag == "" {
		return authremote.Context()
	}
	rt := &rc.CTXRoutingTags{
		Tags:            parse_tags(),
		FallbackToPlain: true,
		Propagate:       false,
	}
	ctx := authremote.ContextWithTimeoutAndTags(time.Duration(2)*time.Second, rt)
	return ctx
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
	fmt.Printf("Added tag \"%s\" with value \"%s\"\n", kv[0], kv[1])
	return res
}

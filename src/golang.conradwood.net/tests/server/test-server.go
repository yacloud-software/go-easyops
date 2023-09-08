package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/echoservice"
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
	rand_port     = flag.Bool("random_port", true, "if true, add random number to port")
	port          = flag.Int("port", 4106, "The grpc server port")
	ping          = flag.Bool("ping", false, "ping continously")
	ping_once     = flag.Bool("ping_once", false, "ping once")
	ping_self     = flag.Bool("ping_self", false, "start server and ping self")
	inject_errors = flag.Bool("inject_errors", false, "inject errors")
	tag           = flag.String("tag", "", "key=value tag optional")
	ctr           = 0
)

// create a simple standard server
type echoServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("GO-EASYOPS Echo test server/client\n")
	flag.VisitAll(func(f *flag.Flag) {
		s := "SET"
		if fmt.Sprintf("%v", f.Value) == fmt.Sprintf("%v", f.DefValue) {
			s = "DEFAULT"
		}
		fmt.Printf("%s %s %s %s\n", "STRING", s, f.Name, f.Value)
	})
	os.Exit(0)
	if *ping || *ping_once {
		c := pb.GetEchoClient()
		for {
			now := time.Now()
			ctx := authremote.Context()
			ctx = authremote.Context()
			ctx = authremote.Context()
			u := auth.GetUser(ctx)
			fmt.Printf("   pinging as %s\n", auth.Description(u))
			_, err := c.Ping(ctx, &common.Void{})
			s := "Pinged"
			if err != nil {
				fmt.Printf("Error :%s\n", utils.ErrorString(err))
				s = "failed"
			}
			dur := time.Since(now).Milliseconds()
			fmt.Printf("%d %s (%d milliseconds)\n", ctr, s, dur)
			ctr++
			if !*ping {
				return
			}
			time.Sleep(time.Duration(300) * time.Millisecond)
		}
	}

	server.SetHealth(server.STARTING)
	sd := server.NewServerDef()
	//	sd.SetPublic()
	sd.DontRegister()
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
	if *rand_port {
		p = p + utils.RandomInt(50)
	}
	if *ping_self {
		go PingSelf()
	}
	sd.AddTag("foo", "bar")
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
func PingSelf() {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		ctx := authremote.Context()
		u := auth.GetUser(ctx)
		fmt.Printf("pinging as %s\n", auth.Description(u))
		c := pb.GetEchoClient()
		_, err := c.Ping(ctx, &common.Void{})
		if err != nil {
			fmt.Printf("Error :%s\n", utils.ErrorString(err))
		} else {
			fmt.Printf("Pinged\n")
		}

	}
}
func (e *echoServer) Ping(ctx context.Context, req *common.Void) (*pb.PingResponse, error) {
	u := auth.GetUser(ctx)
	s := auth.GetService(ctx)
	fmt.Printf("   %03d Pinged by user %s, service %s\n", ctr, auth.Description(u), auth.Description(s))
	ctr++
	if *inject_errors {
		i := utils.RandomInt(10)
		if i > 3 {
			return nil, errors.Unavailable(ctx, "Ping()")
		}
	}
	return &pb.PingResponse{}, nil
}

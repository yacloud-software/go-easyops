package main

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/echoservice"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"os"
	"strings"
	"time"
)

/*
tests:
1) restart postgresql - should reconnect
2) stop postgresql for >5 minutes - should reconnect
3) "blackhole" traffic"
ip addr add 172.29.1.235/24 dev br0
ip addr del 172.29.1.235/24 dev br0


*/

const (
	INTERVAL = time.Duration(300) * time.Millisecond
)

var (
	port      = flag.Int("port", 4106, "The grpc server port")
	ping      = flag.Bool("ping", false, "ping continously")
	ping_once = flag.Bool("ping_once", false, "ping once")
	tag       = flag.String("tag", "", "key=value tag optional")
	ctr       = 0
	dbcon     *sql.DB
)

// create a simple standard server
type echoServer struct {
}

func main() {
	flag.Parse()
	var err error
	dbcon, err = sql.Open()
	utils.Bail("failed to open db: %s", err)
	fmt.Printf("GO-EASYOPS Echo test server/client\n")
	if *ping || *ping_once {
		for {
			do_ping()
			time.Sleep(INTERVAL)
			fmt.Printf("Successes: %d out of %d. Failures: %d out of %d\n",
				dbcon.GetFailureCounter().GetCounter(0),
				dbcon.GetFailureCounter().GetCounts(0),
				dbcon.GetFailureCounter().GetCounter(1),
				dbcon.GetFailureCounter().GetCounts(1),
			)
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
	sd.AddTag("foo", "bar")
	sd.Port = p
	sd.Register = server.Register(
		func(g *grpc.Server) error {
			pb.RegisterEchoServiceServer(g, &echoServer{})
			return nil
		},
	)
	err = server.ServerStartup(sd)
	//	err := create.NewEchoServiceServer(&echoServer{}, p)
	utils.Bail("Unable to start server", err)
}

func (e *echoServer) Ping(ctx context.Context, req *common.Void) (*pb.PingResponse, error) {
	u := auth.GetUser(ctx)
	fmt.Printf("    %d Pinged by %s\n", ctr, auth.Description(u))
	ctr++
	i := utils.RandomInt(10)
	if i > 3 {
		return nil, errors.Unavailable(ctx, "Ping()")
	}
	return &pb.PingResponse{}, nil
}

func do_ping() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	var now time.Time
	rows, err := dbcon.QueryContext(ctx, "nowquery", "SELECT NOW() as now")
	if err != nil {
		fmt.Printf("Query Error: %v\n", err)
		return
	}
	if !rows.Next() {
		fmt.Printf("Next error (no rows)\n")
		return
	}
	err = rows.Scan(&now)
	if err != nil {
		fmt.Printf("Scan Error: %v\n", err)
		return
	}
	rows.Close()

	fmt.Printf("Result: %v\n", now)

}

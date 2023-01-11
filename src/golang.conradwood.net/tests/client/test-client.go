package main

import (
	"flag"
	"fmt"
	au "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/helloworld"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	cm "golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

func main() {
	flag.Parse()
	fmt.Printf("go-easyops test client\n")
	u, s := authremote.GetLocalUsers()
	fmt.Printf("Local User account   : %s\n", user2string(u))
	fmt.Printf("Local Service account: %s\n", user2string(s))

	pingAs("7")
	pingLoop()
	pingLookup()
	pingStream()
}
func pingAs(userid string) {
	fmt.Printf("Pinging with user \"%s\"...\n", userid)
	ctx, err := authremote.ContextForUserID(userid)
	utils.Bail("failed to get context", err)
	started := time.Now()
	con := client.Connect("helloworld.HelloWorld")
	c := helloworld.NewHelloWorldClient(con)
	r, err := c.Ping(ctx, &common.Void{})
	utils.Bail("failed to ping", err)
	fmt.Printf("Pinged (%0.2fs), User=%s, Service=%s, Creator=%s\n", time.Since(started).Seconds(), auth.UserIDString(r.CallingUser), auth.Description(r.CallingService), auth.Description(r.CreatingService))
	reu := r.CallingUser
	if reu == nil || reu.ID != userid {
		fmt.Printf("Creatd context for user \"%s\", but server reported user \"%s\"\n", userid, auth.Description(reu))
		os.Exit(10)
	}
}

func pingLoop() {
	ctx := authremote.Context()
	if ctx == nil {
		fmt.Printf("ERROR: authremote.Context() created no context\n")
		os.Exit(10)
	}
	fmt.Printf("Pinging with default client...\n")
	started := time.Now()
	res, err := helloworld.GetHelloWorldClient().PingLoop(ctx, &helloworld.PingRequest{Loops: 5})
	utils.Bail("failed to ping", err)
	fmt.Printf("Pinged (%0.2fs)\n", time.Since(started).Seconds())
	t := &utils.Table{}
	t.AddHeaders("#", "User", "Service", "Creator")
	for i, r := range res.Responses {
		t.AddInt(i + 1)
		t.AddString(auth.Description(r.CallingUser))
		t.AddString(auth.Description(r.CallingService))
		t.AddString(auth.Description(r.CallingService))
		t.NewRow()
	}
	fmt.Println(t.ToPrettyString())
}
func pingLookup() {
	fmt.Printf("Pinging with lookup...\n")
	ctx := authremote.Context()
	started := time.Now()
	con := client.Connect("helloworld.HelloWorld")
	c := helloworld.NewHelloWorldClient(con)
	r, err := c.Ping(ctx, &common.Void{})
	utils.Bail("failed to ping", err)
	fmt.Printf("Pinged (%0.2fs), User=%s, Service=%s, Creator=%s\n", time.Since(started).Seconds(), auth.Description(r.CallingUser), auth.Description(r.CallingService), auth.Description(r.CreatingService))
}

func pingStream() {
	fmt.Printf("Pinging stream...\n")
	ctx := authremote.Context()
	psreq := &helloworld.PingStreamRequest{DelayInMillis: 500}
	started := time.Now()
	srv, err := helloworld.GetHelloWorldClient().PingStream(ctx, psreq)
	utils.Bail("failed to set up pingstream", err)
	pings := 0
	var user, service, cservice *au.User
	for {
		pr, err := srv.Recv()
		if err != nil {
			fmt.Printf("error received: %s\n", err)
			break
		}
		pings++
		fmt.Printf("Received Sequence %d (Stream %d)\n", pr.SequenceNumber, pr.StreamID)
		print := false
		if pr.User != nil && pr.User != user {
			user = pr.User
			print = true
		}
		if pr.CallingService != nil && pr.CallingService != service {
			service = pr.CallingService
			print = true
		}
		if pr.CreatorService != nil && pr.CreatorService != cservice {
			cservice = pr.CreatorService
			print = true
		}
		if print {
			fmt.Printf("User:%s, Service:%s, Creator:%s\n", auth.Description(user), auth.Description(service), auth.Description(cservice))
		}
	}
	dur := time.Since(started)
	fmt.Printf("Stream ended after %0.2fs and %d pings\n", dur.Seconds(), pings)
	fmt.Printf("Sleeping a few seconds after stream end...\n")
	for i := 5; i > 0; i-- {
		fmt.Printf(" %d\r", i)
		time.Sleep(time.Duration(1) * time.Second)
	}

}
func user2string(u *au.SignedUser) string {
	uu := cm.VerifySignedUser(u)
	return auth.Description(uu)
}

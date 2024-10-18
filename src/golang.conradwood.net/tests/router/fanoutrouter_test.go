package router

import (
	"fmt"
	"testing"
	"time"

	"golang.conradwood.net/apis/apitest"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/router"
)

func TestSingle(t *testing.T) {
	cm := router.NewConnectionManager("apitest.ApiTestService")
	cm.AllowMultipleInstancesPerIP()
	fr := router.NewFanoutRouter(cm, func(p *router.ProcessRequest) error {
		client := apitest.NewApiTestServiceClient(p.GRPCConnection())
		ctx := authremote.Context()
		_, err := client.SlowPing(ctx, &apitest.PingRequest{})
		return err
	},
		func(c *router.CompletionNotification) {
		},
	)
	started := time.Now()
	i := 0
	for time.Since(started) < time.Duration(30)*time.Second {
		i++
		fmt.Printf("submitting %d\n", i)
		fr.SubmitWork("foo")
	}
	fr.Stop()
	dur := time.Since(started)
	ps := float64(i) / dur.Seconds()
	fmt.Printf("Processed %d requests in %0.1fs (%0.1f per second)\n", i, dur.Seconds(), ps)
}

package client

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	pp "golang.conradwood.net/go-easyops/profiling"
	"google.golang.org/grpc"
	"time"
)

func unaryStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	pp.ClientRpcEntered()
	var err error
	s := "X"
	m := "Y"
	if *debug_rpc_client {
		s, m, err = splitMethodAndService(method)
		if err != nil {
			fmt.Printf("[go-easyops] failed to split method \"%s\": %s\n", method, err)
		}
		fmt.Printf("[go-easyops] invoking streaming rpc \"%s/%s\" as user %s\n", s, m, auth.UserIDString(auth.GetUser(ctx)))
	}
	started := time.Now()
	cs, err := streamer(ctx, desc, cc, method, opts...)
	pp.ClientRpcDone()
	dur := time.Since(started)
	if *debug_rpc_client {
		if err != nil {
			fmt.Printf("[go-easyops] streaming rpc \"%s/%s\" failed after %0.2fs with error %s\n", s, m, dur.Seconds(), err)
		} else {
			fmt.Printf("[go-easyops] streaming rpc \"%s/%s\" returned after %0.2fs\n", s, m, dur.Seconds())
		}
	}
	return cs, err
}

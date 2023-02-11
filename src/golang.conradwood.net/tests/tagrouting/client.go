package main

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/getestservice"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func do_client() error {
	c := pb.GetEchoClient()
	ctx := clientContext()
	c.Ping(ctx, &pb.PingRequest{SequenceNumber: 0, TTL: 0})
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("--------------------------------------------------\n")
	seq := uint32(0)
	my_tags := parse_tags()
	for {
		now := time.Now()
		ctx := clientContext()
		u := auth.GetUser(ctx)
		seq++
		fmt.Printf("  %s pinging, SEQ=%d, as %s...", printTags(), seq, auth.Description(u))
		resp, err := c.Ping(ctx, &pb.PingRequest{SequenceNumber: seq, TTL: uint32(*ttl)})
		if err != nil {
			fmt.Printf("Error :%s\n", utils.ErrorString(err))
		}
		if !tagsEqual(my_tags, resp.ServerTags) {
			fmt.Printf("TAG mismatch (local:%s vs server:%s)\n", printMap(my_tags), printMap(resp.ServerTags))
			return fmt.Errorf("tag mismatch")
		}
		dur := time.Since(now).Milliseconds()
		fmt.Printf("%d milliseconds\n", dur)
		if !*ping {
			return nil
		}
		time.Sleep(time.Duration(300) * time.Millisecond)
	}
}

func clientContext() context.Context {
	if *tag == "" {
		return authremote.Context()
	}
	rt := &rc.CTXRoutingTags{
		Tags:            parse_tags(),
		FallbackToPlain: *fallback,
		Propagate:       false,
	}
	ctx := authremote.ContextWithTimeoutAndTags(time.Duration(2)*time.Second, rt)
	return ctx
}
func printTags() string {
	return fmt.Sprintf("\"%s\"", *tag)
}
func printMap(m map[string]string) string {
	deli := ""
	s := ""
	for k, v := range m {
		s = s + deli + k + "=" + v
		deli = ", "
	}
	return "\"" + s + "\""
}

func tagsEqual(cl, srv map[string]string) bool {
	if *fallback {
		if len(srv) == 0 {
			return true
		}
	}
	if len(cl) != len(srv) {
		return false
	}
	for k, v := range cl {
		if srv[k] != v {
			return false
		}
	}
	return true
}

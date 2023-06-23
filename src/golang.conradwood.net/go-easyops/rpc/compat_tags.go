package rpc

import (
// ge "golang.conradwood.net/apis/goeasyops"
)

/*
// helper to transition away from rpcinterceptor (old:rpcinterceptor,new:goeasyops)
func Tags_rpc_to_ge(r *rc.CTXRoutingTags) *ge.CTXRoutingTags {
	if r == nil {
		return nil
	}
	res := &ge.CTXRoutingTags{
		FallbackToPlain: r.FallbackToPlain,
		Propagate:       r.Propagate,
		Tags:            r.Tags,
	}
	return res
}

// helper to transition away from rpcinterceptor (old:rpcinterceptor,new:goeasyops)
func Tags_ge_to_rpc(r *ge.CTXRoutingTags) *rc.CTXRoutingTags {
	if r == nil {
		return nil
	}
	res := &rc.CTXRoutingTags{
		FallbackToPlain: r.FallbackToPlain,
		Propagate:       r.Propagate,
		Tags:            r.Tags,
	}
	return res
}
*/

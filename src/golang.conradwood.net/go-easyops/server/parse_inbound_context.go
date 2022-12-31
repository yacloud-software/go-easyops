package server

import (
	"context"
	"golang.conradwood.net/go-easyops/rpc"
)

/*
	given a context (as received by gRPC) parse it into CallState.

returns:

	true - if v2 context was found, false if not
	error - if v2 context was found, but failed to parse.
	(no v2 context is not an error)
*/
func parse_inbound_context_v2(ctx context.Context, out *rpc.CallState) (bool, error) {
	return false, nil
}

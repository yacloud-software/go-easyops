package server

import (
	"context"
	"fmt"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
)

/*
 given a context (as received by gRPC) parse it into CallState.
returns:
 true - if v2 context was found, false if not
 error - if v2 context was found, but failed to parse.
 (no v2 context is not an error)
*/
func parse_inbound_context_v2(ctx context.Context, out *rpc.CallState) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("No metadata at all\n")
		return false, nil
	}
	gemd, ok := md[tokens.METANAME2]
	if !ok {
		return false, nil
	}
	if len(gemd) == 0 {
		return true, fmt.Errorf("go-easyops metadata has 0 entries")
	}
	ges := gemd[0]
	if len(ges) == 0 {
		return true, fmt.Errorf("go-easyops metadata has 0 length")
	}
	inctx := &ge.InContext{}
	err := utils.Unmarshal(ges, inctx)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal received metadata: %s", err)
	}
	fmt.Printf("Metadata: \"%v\"\n", inctx)
	out.SetV2(rpc.NewCallStateV2(inctx))
	ctx = context.WithValue(ctx, rpc.LOCALCONTEXTNAME, out)
	out.Context = ctx
	return true, nil
}

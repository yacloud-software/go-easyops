package client

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/rpc"
	"google.golang.org/grpc/balancer"
	//	"google.golang.org/grpc/metadata"
)

type FancyPicker struct {
	addresslist *FancyAddressList
	failAll     bool // if true all RPCs will fail
}

var (
	ctr         uint32
	honour_tags = flag.Bool("ge_honour_tags", true, "whether or not to honour tag-based routing")
)

// Pick returns the connection to use for this RPC and related information.
//
// Pick should not block.  If the balancer needs to do I/O or any blocking
// or time-consuming work to service this call, it should return
// ErrNoSubConnAvailable, and the Pick call will be repeated by gRPC when
// the Picker is updated (using ClientConn.UpdateState).
//
// If an error is returned:
//
// - If the error is ErrNoSubConnAvailable, gRPC will block until a new
//   Picker is provided by the balancer (using ClientConn.UpdateState).
//
// - If the error implements IsTransientFailure() bool, returning true,
//   wait for ready RPCs will wait, but non-wait for ready RPCs will be
//   terminated with this error's Error() string and status code
//   Unavailable.
//
// - Any other errors terminate all RPCs with the code and message
//   provided.  If the error is not a status error, it will be converted by
//   gRPC to a status error with code Unknown.
func (f *FancyPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if f.failAll {
		// the balancer created a special "failing" picker because it did not have any
		// instances for this service for a long time (so it is not transient anymore, is it?)
		// in this case we don't want to build up a queue of RPCs, we just want to fail-fast them
		fancyPrintf(f, "Picker - failing connections for \"%s\" w/o instance\n", info.FullMethodName)
		cs := rpc.CallStateFromContext(info.Ctx)
		sn := "[unknown rpc]"
		if cs != nil {
			sn = fmt.Sprintf("%s.%s()", cs.ServiceName, cs.MethodName)
		}
		return balancer.PickResult{}, fmt.Errorf("failure in %s whilst calling %s - no backend available", sn, info.FullMethodName)
	}
	if f.addresslist.IsEmpty() {
		// no instances, transient problem though. we tell gRPC to retry the call once we got a new picker
		fancyPrintf(f, "Picker - No connections for %s\n", info.FullMethodName)
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	lf := f.addresslist
	if *honour_tags {
		value := info.Ctx.Value("routingtags")
		if value != nil {
			// convert tags to map[string]string, returning empty if invalid type assertion
			tags, ok := value.(map[string]string)
			if !ok {
				return balancer.PickResult{}, fmt.Errorf("Invalid tags object supplied (%v)", value)
			}
			adrs := lf.ByMatchingTags(tags)
			if len(adrs) == 0 {
				fancyPrintf(f, "Picker - No connection matched all required tags (%v)\n", tags)
				return balancer.PickResult{}, fmt.Errorf("No addresses matched all supplied tags (%v)", tags)
			}
			lf = &FancyAddressList{Name: lf.Name, addresses: adrs}
		}
	}

	// build up a list of valid (e.g. state Ready, match user/context/routing) connections
	matching := lf.SelectValid(info.Ctx)

	if len(matching) == 0 {
		fancyPrintf(f, "Picker - No valid connections for %s\n", info.FullMethodName)
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	ctr++ // overflows at 0xFFFFFFFF, that's ok

	indx := ctr % uint32(len(matching))
	fa := matching[indx]
	if *debug_fancy {
		u := auth.GetUser(info.Ctx)
		fancyPrintf(f, "Picking: %s [%s] for user %s to serve %s from %d connections (%d matching))\n",
			fa.addr, fa.state.String(),
			auth.Description(u),
			info.FullMethodName,
			lf.Count(), len(matching))
		fancyPrintf(f, "         RoutingInfo: %#v\n", fa.Target.RoutingInfo)
	}

	res := balancer.PickResult{SubConn: fa.subcon}
	return res, nil
}

func (f *FancyPicker) ServiceName() string {
	if f.addresslist != nil {
		return f.addresslist.Name
	}
	return "fancy_picker.go"
}

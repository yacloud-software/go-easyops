package client

import (
	"fmt"
	//	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/auth"
	//	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	//	"google.golang.org/grpc/metadata"
)

type FancyPicker struct {
	addresslist *FancyAddressList
	failAll     bool // if true all RPCs will fail
	ctr         uint32
}

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
	cs := rpc.CallStateFromContext(info.Ctx)
	if f.failAll {
		// the balancer created a special "failing" picker because it did not have any
		// instances for this service for a long time (so it is not transient anymore, is it?)
		// in this case we don't want to build up a queue of RPCs, we just want to fail-fast them
		fancyPrintf(f, "Picker - failing connections for \"%s\" w/o instance\n", info.FullMethodName)
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

	cri := cs.RoutingTags()
	if cri != nil {
		fancyPrintf(f, "Picking by tags (%v)\n", cri.Tags)
		// convert tags to map[string]string, returning empty if invalid type assertion
		adrs := lf.ByMatchingTags(cri.Tags)
		if len(adrs) == 0 {
			fancyPrintf(f, "Picker - No connection matched all required tags (%v)\n", cri.Tags)
			if !cri.FallbackToPlain {
				return balancer.PickResult{}, fmt.Errorf("No addresses matched all supplied tags (%v) for %s", cri.Tags, info.FullMethodName)
			} else {
				lf = f.addresslist
				lf = &FancyAddressList{Name: lf.Name, addresses: lf.ByWithoutTags()}
			}
		} else {
			lf = &FancyAddressList{Name: lf.Name, addresses: adrs}
		}
	}

	// build up a list of valid (e.g. state Ready, match user/context/routing) connections
	matching := lf.SelectValid(info.Ctx)

	if len(matching) == 0 {
		for _, a := range lf.addresses {
			// this is not right here. We probably should periodically keep them alive rather than wait until
			// we have no more READY ones
			// but this is a 'hotfix' to stop breakage
			if a.state == connectivity.Idle {
				// this doesn't do the trick. it just makes it worse actually,
				// it covers for quick reconnects on the same port only, but breaks after long periods too
				//a.subcon.Connect()
			}
			fancyPrintf(f, "picker address: %s\n", a.String())
		}
		fancyPrintf(f, "Picker - No valid connections for %s\n", info.FullMethodName)
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	f.ctr++ // overflows at 0xFFFFFFFF, that's ok

	indx := f.ctr % uint32(len(matching))
	fa := matching[indx]
	if *debug_fancy {
		u := auth.GetUser(info.Ctx)
		fancyPrintf(f, "Picking: %s [%s] for user %s to serve %s from %d connections (%d matching) (ctr=%d))\n",
			fa.addr, fa.state.String(),
			auth.Description(u),
			info.FullMethodName,
			lf.Count(), len(matching), f.ctr)
		fancyPrintf(f, "         RoutingInfo: %#v\n", fa.Target.RoutingInfo)
	}

	res := balancer.PickResult{SubConn: fa.subcon}
	fa.subcon.Connect()
	return res, nil
}

func (f *FancyPicker) ServiceName() string {
	if f.addresslist != nil {
		return f.addresslist.Name
	}
	return "fancy_picker.go"
}

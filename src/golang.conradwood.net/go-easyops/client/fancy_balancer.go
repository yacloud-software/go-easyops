package client

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

func init() {
	balancer.Register(&FancyBuilder{})
	go balancer_thread()
}

var (
	balancers []*FancyBalancer
	bal_lock  sync.Mutex
	maxblock  = flag.Float64("ge_max_block", 30, "max `seconds` to block rpcs for if backends are not available (fail afterwards)")
)

/*********** the builder for our balancer *****************/
type FancyBuilder struct {
}

// Build creates a new balancer for the (new) ClientConn.
func (f *FancyBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	fancyPrintf(f, "Building Balancer for %s\n", opts.Target.Endpoint)
	fal := &FancyAddressList{Name: opts.Target.Endpoint}
	cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            &FancyPicker{addresslist: fal}, // not failing - initially we wait
	})
	res := &FancyBalancer{cc: cc,
		target:       opts.Target.Authority,
		blockedSince: time.Now(),
		addresslist:  fal,
	}
	if res.target == "" {
		s := fmt.Sprintf("cannot create fancy-balancer without servicename (opts=%#v). Dial must be in the format 'go-easyops://servicename/servicename@registry'", opts)
		panic(s)
	}
	// looks a bit dumb. we really should reuse slots from closed ones
	bal_lock.Lock()
	defer bal_lock.Unlock()
	balancers = append(balancers, res)
	return res
}

// Name returns the name of balancers built by this builder.
// It will be used to pick balancers (for example in service config).
func (f *FancyBuilder) Name() string {
	return "fancybalancer"
}

/*********** the balancer *****************/

type FancyBalancer struct {
	cc           balancer.ClientConn
	target       string
	addresslist  *FancyAddressList
	closed       bool
	failing      bool
	blockedSince time.Time
}

// EXPERIMENTAL: this is the new-style grpc callback
func (f *FancyBalancer) ResolverError(err error) {
	fmt.Printf("[go-easyops] Resolver reported an error, which is not handled yet: %s\n", err)
}

// EXPERIMENTAL: this is the new-style grpc callback, called by the resolver when a state changes
// it feeds us new addresses
func (f *FancyBalancer) UpdateClientConnState(bc balancer.ClientConnState) error {
	fancyPrintf(f, "balancer: updateclientconnstate (ResolverState: %d addresses)\n", len(bc.ResolverState.Addresses))
	f.HandleResolvedAddrs(bc.ResolverState.Addresses, nil)
	return nil
}

// EXPERIMENTAL: this is the new-style grpc callback
func (f *FancyBalancer) UpdateSubConnState(sc balancer.SubConn, bc balancer.SubConnState) {
	f.HandleSubConnStateChange(sc, bc.ConnectivityState)
}

// DEPRECATED - old version of grpc
// HandleSubConnStateChange is called by gRPC when the connectivity state
// of sc has changed.
// Balancer is expected to aggregate all the state of SubConn and report
// that back to gRPC.
// Balancer should also generate and update Pickers when its internal state has
// been changed by the new state.
func (f *FancyBalancer) HandleSubConnStateChange(sc balancer.SubConn, state connectivity.State) {
	fa := f.addresslist.BySubCon(sc)
	if fa == nil {
		fancyPrintf(f, "balancer: SubConnState on a subconnection we don't know (%#v)!\n", sc)
		return
	}
	oldstate := fa.state
	fa.state = state

	fancyPrintf(f, "balancer: Handlesubstate service %s at %s transitioned from %s to %s\n", f.target, fa.addr, oldstate.String(), state.String())
	f.failing = false
	f.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            f.Picker(),
	})
}

// DEPRECATED - old version of grpc
// HandleResolvedAddrs is called by gRPC to send updated resolved addresses to
// balancers.
// Balancer can create new SubConn or remove SubConn with the addresses.
// An empty address slice and a non-nil error will be passed if the resolver returns
// non-nil error to gRPC.
// Note that each address MUST have an attribute with a registry.ServiceAddress
// because we cannot transport all the information in just ip/port
func (f *FancyBalancer) HandleResolvedAddrs(addresses []resolver.Address, err error) {
	//	fancyPrintf(f,"balancer: HandleResolveAddrs addressed = %d (err=%v)\n", len(addresses), err)

	// create new ones:
	added := false
	for _, resolverAddr := range addresses {
		var sa *registry.Target
		if resolverAddr.Attributes != nil {
			o := resolverAddr.Attributes.Value(RESOLVER_ATTRIBUTE_SERVICE_ADDRESS)
			if o != nil {
				sa = o.(*registry.Target)
			}
		}
		if sa == nil {
			// see note above. serviceAddress is required!
			s := fmt.Sprintf("fancy balancer received a very unfancy address without registry.ServiceAddress attribute for \"%s\". Weird resolver?", f.target)
			panic(s)
		}
		rf := ""
		ri := sa.RoutingInfo
		if ri != nil {
			u := ri.RunningAs
			if u != nil {
				rf = fmt.Sprintf("runningas=%s/#%s", auth.Description(u), u.ID)
			}
		}
		fr := f.addresslist.ByAddr(resolverAddr.Addr)
		if fr != nil {
			fr.Target = sa
			f.addresslist.Updated()
			fancyPrintf(f, "balancer: %s, conn %s known as state %s\n", f.target, resolverAddr.Addr, fr.state.String())
			continue
		}
		fancyPrintf(f, "balancer: New Address %s (%s)\n", resolverAddr.Addr, rf)
		// not yet known - create a new sub connection
		sco, err := f.cc.NewSubConn([]resolver.Address{resolverAddr}, balancer.NewSubConnOptions{})
		if err != nil {
			fancyPrintf(f, "Failed to create subconn: %s\n", err)
			continue
		}
		//	sc = append(sc, sco)
		f.addresslist.Add(&fancy_adr{
			state:  connectivity.Ready, // docs say use CONNECTING here, but that never calls the picker nor the stateupdate. how does that work?
			addr:   resolverAddr.Addr,
			subcon: sco,
			Target: sa,
		})
		added = true
	}
	// we also need to remove connections which are no longer valid for this service:
	remlist := f.addresslist.RequiredList(addresses)
	for _, fa := range remlist {
		f.cc.RemoveSubConn(fa.subcon)
	}
	removed := len(remlist) != 0
	if !added && !removed {
		fancyPrintf(f, "balancer: no state change for \"%s\"\n", f.target)
		return
	}
	f.failing = false
	fancyPrintf(f, "balancer: Sending state update for \"%s\", we got %d subconnections now\n", f.target, f.addresslist.Count())

	if f.addresslist.IsEmpty() {
		f.blockedSince = time.Now()
		f.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.Ready,
			//ConnectivityState: connectivity.TransientFailure,
			Picker: f.Picker(),
		})
		return
	}
	f.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            f.Picker(),
	})

}

// Close closes the balancer. The balancer is not required to call
// ClientConn.RemoveSubConn for its existing SubConns.
func (f *FancyBalancer) Close() {
	f.closed = true
	bal_lock.Lock()
	defer bal_lock.Unlock()
	var res []*FancyBalancer
	// looks a bit dumb. we really should reuse slots from closed ones
	for _, b := range balancers {
		if b.closed {
			continue
		}
		res = append(res, b)
	}
	balancers = res
	fancyPrintf(f, "Close\n")
}

// create a new picker
func (f *FancyBalancer) Picker() *FancyPicker {
	res := &FancyPicker{addresslist: f.addresslist}
	return res
}

/********************************* thread checking for hung pickers ************************
we want RPCs to fail (rather than indefinitely hang)
Behaviour should be like this:
* if we have a "transient failure" we report it as such
* if the "transient failure" remains for some time, we start failing RPCs
* (this is to avoid them all queueing up and blocking up our service(s))
*/

func balancer_thread() {
	for {
		for _, b := range balancers {
			b.Check()
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

}

// periodically called by go routine, checks if it's blocking for too long
func (f *FancyBalancer) Check() {
	if !f.addresslist.IsEmpty() {
		return // not blocked
	}
	if f.failing {
		return // already failing
	}
	sc := time.Since(f.blockedSince)
	fancyPrintf(f, "Blocked since: %v (%v)\n", f.blockedSince, sc)
	if sc < (time.Duration(*maxblock) * time.Second) {
		return // not failing for long enough to take action yet
	}

	f.failing = true
	fp := f.Picker()
	fp.failAll = true
	f.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Ready,
		Picker:            fp,
	})

}
func (f *FancyBalancer) ServiceName() string {
	return f.target
}
func (f *FancyBuilder) ServiceName() string {
	return "fancy_builder.go"
}

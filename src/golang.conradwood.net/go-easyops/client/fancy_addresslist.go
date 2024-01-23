package client

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
	"strings"
)

/*
this is essentially a list of addresses.
the balancer removes/adds/updates addresses and the
picker reads/chooses/sorts them.
this struct synchronises access between them.

The list is semi-uptodate, that means, it is cached, but updated if go-easyops determines that the registry
has better information than its cache.
The Addresses in this list are still subject to the filtering done in the registry. The Registry "prefers" certain targets, for example, higher buildids
*/
type FancyAddressList struct {
	Name      string
	addresses []*FancyAddr
}

type FancyAddr struct {
	addr     string
	subcon   balancer.SubConn
	state    connectivity.State
	removed  bool
	Target   *registry.Target
	grpc_con *grpc.ClientConn // only used if client calls Connect() on this
}

// a key that can be used in maps to find this particular fancyaddress.
func (fa *FancyAddr) Key() string {
	return "faddr_" + fa.addr
}

// address, including port, e.g. 10.1.1.1:6000
func (fa *FancyAddr) Address() string {
	return fa.addr
}
func (fa *FancyAddr) String() string {
	return fmt.Sprintf("%s: %s[%s] removed=%v", fa.Target.ServiceName, fa.addr, fa.state.String(), fa.removed)
}

// return true if this is _actually_ available. e.g. a TCP reset will cause this connection to be "not ready", but still be listed in the registry and caches
func (fa *FancyAddr) IsReady() bool {
	return fa.state == connectivity.Ready
}

func (fa *FancyAddr) disconnect() {
	if fa.grpc_con == nil {
		return
	}
	fa.grpc_con.Close()
	fa.grpc_con = nil
}

/*
 open and maintain a connection to this peer. This can help to build custom load-balancers, but is not intented for general-use.
 use with caution - using this method required in-depth knowledge of grpc and go-easyops
*/

func (fa *FancyAddr) Connection() (*grpc.ClientConn, error) {
	if fa.grpc_con != nil {
		return fa.grpc_con, nil
	}
	gc, err := grpc.Dial(fmt.Sprintf("ipv4:%s", fa.addr), grpc.WithBlock(),
		grpc.WithTransportCredentials(GetClientCreds()),
		grpc.WithUnaryInterceptor(ClientMetricsUnaryInterceptor),
		grpc.WithStreamInterceptor(unaryStreamInterceptor),
	)
	if err != nil {
		return nil, err
	}
	fa.grpc_con = gc
	return gc, nil
}
func (fal *FancyAddressList) Count() int {
	return len(fal.addresses)
}
func (fal *FancyAddressList) IsEmpty() bool {
	return len(fal.addresses) == 0
}

// called by the balancer when a FancyAddr has been updated. (or anyone updating FancyAddr)
// we may need to clear some caches (now or in future...)
func (fal *FancyAddressList) Updated() {
}

// perhaps should check/panic on duplicates here?
func (fal *FancyAddressList) Add(f *FancyAddr) {
	fal.addresses = append(fal.addresses, f)
	fal.Updated()
}

// returns the fancyaddress that matches the key. see fancyaddress.Key(). this might return nil if no such fancyaddress is known (any more)
func (fal *FancyAddressList) ByKey(key string) *FancyAddr {
	for _, fa := range fal.addresses {
		if fa.Key() == key {
			return fa
		}
	}
	return nil
}

// removes all addresses which are NOT in the array and returns the removed ones
func (fal *FancyAddressList) RequiredList(addresses []resolver.Address) []*FancyAddr {
	var res []*FancyAddr
	removed := false
	for _, fa := range fal.addresses {
		stillgood := false
		for _, r := range addresses {
			if r.Addr == fa.addr {
				stillgood = true
				break
			}
		}
		if stillgood {
			continue
		}
		fancyPrintf(fal, "balancer: removed %s\n", fa.addr)
		fa.removed = true
		removed = true
		fa.disconnect()
	}
	if removed {
		var fa []*FancyAddr
		for _, foa := range fal.addresses {
			if foa.removed {
				res = append(res, foa)
				continue
			}
			fa = append(fa, foa)
		}
		fal.addresses = fa
		fal.Updated()
	}
	return res
}

/******************************** find entries by various keys *************************/

func (fal *FancyAddressList) ByAddr(adr string) *FancyAddr {
	for _, fa := range fal.addresses {
		if fa.addr == adr {
			return fa
		}
	}
	return nil
}

func (fal *FancyAddressList) BySubCon(sc balancer.SubConn) *FancyAddr {
	var fa *FancyAddr
	for _, a := range fal.addresses {
		if a.subcon == sc {
			fa = a
			break
		}
	}
	return fa
}

// return all addresses the fancyaddresslist knows about.
func (fal *FancyAddressList) AllAddresses() []*FancyAddr {
	var valids []*FancyAddr
	for _, a := range fal.addresses {
		valids = append(valids, a)
	}
	return valids
}

// return all "ready" addresses (those with a TCP connection in ready state)
func (fal *FancyAddressList) AllReadyAddresses() []*FancyAddr {
	var valids []*FancyAddr
	for _, a := range fal.addresses {
		if !a.IsReady() {
			continue
		}
		valids = append(valids, a)
	}
	return valids
}

// return addresses with 0 tags
func (fal *FancyAddressList) ByWithoutTags() []*FancyAddr {
	var valids []*FancyAddr
	// filter addresses to include only those which contain required all tags
	for _, a := range fal.addresses {
		if a.Target == nil {
			continue
		}
		if a.Target.RoutingInfo == nil || a.Target.RoutingInfo.Tags == nil || len(a.Target.RoutingInfo.Tags) == 0 {
			valids = append(valids, a)
		}
	}
	return valids
}

/*
called for _every_ rpc call when ge_honour_tags flag is true, adjusts the
list of matches by checking whether the addresses matches all the routing tags
supplied.
if no tags are supplied, return _ALL_ targets (including those with tags)
*/

func (fal *FancyAddressList) ByMatchingTags(tags map[string]string) []*FancyAddr {
	fancyPrintf(fal, "Filtering (%d) addresses by tags\n", len(fal.addresses))
	if len(tags) == 0 {
		// no point iterating over all addresses if we have no tags
		// also - we treat "empty list" the same as "no tags specified", that is, return all connections instead of none
		fancyPrintf(fal, "empty list for filterbytags!")
		return fal.addresses
	}
	var valids []*FancyAddr
	// filter addresses to include only those which contain required all tags
	for _, a := range fal.addresses {
		valid := true
		if a.Target == nil || a.Target.RoutingInfo == nil || a.Target.RoutingInfo.Tags == nil {
			fancyPrintf(fal, "tag in %s does have special routing\n", a.addr)
			continue
		}
		for k, v := range tags {
			tgv := a.Target.RoutingInfo.Tags[k]
			if tgv != v {
				fancyPrintf(fal, "tag in %s does not match. \"%s\" != \"%s\"\n", a.addr, tgv, v)
				valid = false
				break
			}
		}
		if valid {
			valids = append(valids, a)
		}
	}
	return valids
}

// get all those without routinginfo or no routinginfo.user
func (fal *FancyAddressList) ByNoUserRoutingInfo() []*FancyAddr {
	var res []*FancyAddr
	for _, a := range fal.addresses {
		ri := a.Target.RoutingInfo
		if ri != nil && ri.RunningAs != nil {
			continue
		}
		res = append(res, a)
	}
	return res
}

// get all those with a routinginfo RunningAs user
func (fal *FancyAddressList) ByUser(userid string) []*FancyAddr {
	var res []*FancyAddr
	for _, a := range fal.addresses {
		ri := a.Target.RoutingInfo
		if ri == nil || ri.RunningAs == nil {
			continue
		}
		if ri.RunningAs.ID != userid {
			continue
		}
		res = append(res, a)
	}
	return res
}

func (fal *FancyAddressList) readyOnly(in []*FancyAddr) []*FancyAddr {
	var valids []*FancyAddr
	bal_state_lock.Lock()
	for _, fa := range in {
		if fa.state != connectivity.Ready {
			continue
		}
		valids = append(valids, fa)
	}
	bal_state_lock.Unlock()
	return valids

}

/*
	this is called for _every_ rpc call. it should be performance optimised
	this returns a list of addresses for the picker to pick from.
	this function is what the picker calls. if loadbalancing is implemented by the user, this
	function should be used

the rules are:
First: Never return any addresses which are not in connectivty state READY.
Then from the remaining addresses (in ready state), follow these rules:
1. If we have 0 addresses with routinginfo for a user, return all.
2. if context has no user, return those without routinginfo.user
3. if 1 or more addresses have a routinginfo.user that matches user in context, return only those
 4. otherwise return those without routinguser.info
*/
func (fal *FancyAddressList) SelectValid(ctx context.Context) []*FancyAddr {
	nro := fal.ByNoUserRoutingInfo()
	if len(nro) == len(fal.addresses) {
		// ALL addresses have routinginfo, so we have 0 addresses WITHOUT routinginfo
		fancyPrintf(fal, "all addresses for %s have routinginfo\n", fal.Name)
		return fal.readyOnly(nro)
	}
	u := auth.GetUser(ctx)
	if u == nil {
		// user has no context, return those without routinginfo
		fancyPrintf(fal, "user-less rpc\n")
		if len(nro) == 0 && len(fal.addresses) > 0 {
			fmt.Printf("[go-easyops] Warning - of %d targets, all require a user in outbound context (but none provided)\n", len(fal.addresses))
		}
		return fal.readyOnly(nro)
	}
	bu := fal.ByUser(u.ID)
	if len(bu) == 0 {
		// none for this user, return those without routinginfo
		fancyPrintf(fal, "no connections specifically for user %s\n", u.Email)
		return fal.readyOnly(nro)
	}
	fancyPrintf(fal, "%d connections specifically for user %s\n", len(bu), u.Email)
	return fal.readyOnly(bu)
}

// servicename, e.g. "registry.Registry"
func (fal *FancyAddressList) ServiceName() string {
	s := fal.Name
	idx := strings.Index(s, "@")
	if idx != -1 {
		s = s[:idx]
	}
	return s
}

func GetAllFancyAddressLists() []*FancyAddressList {
	var res []*FancyAddressList
	for _, bal := range balancers {
		if bal.addresslist != nil {
			res = append(res, bal.addresslist)
		}
	}
	return res
}

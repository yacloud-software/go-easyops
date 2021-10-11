package client

import (
	"context"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

// this is essentially a list of addresses
// the balancer removes/adds/updates addresses and the
// picker reads/chooses/sorts them.
// this struct synchronises access between them
type FancyAddressList struct {
	Name      string
	addresses []*fancy_adr
}

type fancy_adr struct {
	addr    string
	subcon  balancer.SubConn
	state   connectivity.State
	removed bool
	Target  *registry.Target
}

func (fal *FancyAddressList) Count() int {
	return len(fal.addresses)
}
func (fal *FancyAddressList) IsEmpty() bool {
	return len(fal.addresses) == 0
}

// called by the balancer when a fancy_adr has been updated. (or anyone updating fancy_adr)
// we may need to clear some caches (now or in future...)
func (fal *FancyAddressList) Updated() {
}

// perhaps should check/panic on duplicates here?
func (fal *FancyAddressList) Add(f *fancy_adr) {
	fal.addresses = append(fal.addresses, f)
	fal.Updated()
}

// removes all addresses which are NOT in the array and returns the removed ones
func (fal *FancyAddressList) RequiredList(addresses []resolver.Address) []*fancy_adr {
	var res []*fancy_adr
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
	}
	if removed {
		var fa []*fancy_adr
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

func (fal *FancyAddressList) ByAddr(adr string) *fancy_adr {
	for _, fa := range fal.addresses {
		if fa.addr == adr {
			return fa
		}
	}
	return nil
}

func (fal *FancyAddressList) BySubCon(sc balancer.SubConn) *fancy_adr {
	var fa *fancy_adr
	for _, a := range fal.addresses {
		if a.subcon == sc {
			fa = a
			break
		}
	}
	return fa
}

/*
called for _every_ rpc call when ge_honour_tags flag is true, adjusts the
list of matches by checking whether the addresses matches all the routing tags
supplied
*/
//func filterByTags(sn serviceNamer, in []*fancy_adr, tags map[string]string) []*fancy_adr {
func (fal *FancyAddressList) ByMatchingTags(tags map[string]string) []*fancy_adr {
	fancyPrintf(fal, "Filtering (%d) addresses by tags\n", len(fal.addresses))
	if len(tags) == 0 {
		// no point iterating over all addresses if we have no tags
		// also - we treat "empty list" the same as "no tags specified", that is, return all connections instead of none
		fancyPrintf(fal, "empty list for filterbytags!")
		return fal.addresses
	}
	var valids []*fancy_adr
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
func (fal *FancyAddressList) ByNoUserRoutingInfo() []*fancy_adr {
	var res []*fancy_adr
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
func (fal *FancyAddressList) ByUser(userid string) []*fancy_adr {
	var res []*fancy_adr
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

func (fal *FancyAddressList) readyOnly(in []*fancy_adr) []*fancy_adr {
	var valids []*fancy_adr
	for _, fa := range in {
		if fa.state != connectivity.Ready {
			continue
		}
		valids = append(valids, fa)
	}
	return valids

}

/*
 this is called for _every_ rpc call. it should be performance optimised
 this returns a list of addresses for the picker to pick from.
 the rules are:
 First: Never return any addresses which are not in connectivty state READY.
 Then from the remaining addresses (in ready state), follow these rules:
 1. If we have 0 addresses with routinginfo for a user, return all.
 2. if context has no user, return those without routinginfo.user
 3. if 1 or more addresses have a routinginfo.user that matches user in context, return only those
 4. otherwise return those without routinguser.info
*/
func (fal *FancyAddressList) SelectValid(ctx context.Context) []*fancy_adr {
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

func (fal *FancyAddressList) ServiceName() string {
	return fal.Name
}

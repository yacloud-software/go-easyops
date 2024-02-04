package client

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"
)

const (
	REFRESH = time.Duration(5) * time.Second
)

var (
	debug_custom_resolver    = flag.Bool("ge_debug_custom_fancy_resolver", false, "debug the custom fancy resolver, if any are being used")
	custom_fancyadrlist_lock sync.Mutex
)

type custom_fancy_addresslist struct {
	fal      *FancyAddressList
	resolver func(registryadr, servicename string) ([]*registry.Target, error)
}

/*
this creates a fancy address list, which does not automatically maintain its list of targets. it also does not (and cannot) be used for loadbalancing, because it has no connection to a grpc LoadBalancer or Picker.
sometimes maybe might need a custom resolver for a fancy address list.
*/
func NewFancyAddressListWithResolver(servicename string, resolver func(registryadr, servicename string) ([]*registry.Target, error)) (*FancyAddressList, error) {
	fal := &FancyAddressList{Name: servicename}
	cfa := &custom_fancy_addresslist{fal: fal, resolver: resolver}
	targets, err := resolver(cmdline.GetClientRegistryAddress(), servicename)
	if err != nil {
		return nil, err
	}
	for _, fa := range targets_to_fancyadr(targets) {
		fal.Add(fa)
	}
	go cfa.custom_fal_updater(REFRESH)
	return fal, nil
}

func (cfa *custom_fancy_addresslist) custom_fal_updater(refresh time.Duration) {
	for {
		time.Sleep(refresh)
		targets, err := cfa.resolver(cmdline.GetClientRegistryAddress(), cfa.fal.Name)
		if err != nil {
			fmt.Printf("[go-easyops] custom fancyaddresslist resolver failed: %s\n", err)
			continue
		}
		cfa.debugf("got %d targets from resolver\n", len(targets))
		fas := targets_to_fancyadr(targets)
		cfa.debugf("got %d fancyaddress from %d targets from resolver\n", len(fas), len(targets))
		ac := utils.CompareArray(fas, cfa.fal.addresses, func(i, j int) bool {
			return fas[i].Key() == cfa.fal.addresses[j].Key()
		})
		for _, fa_idx := range ac.ElementsIn1ButNot2() {
			cfa.debugf("adding %s\n", fas[fa_idx])
			cfa.fal.Add(fas[fa_idx])
		}
		for _, fa_idx := range ac.ElementsIn2ButNot1() {
			cfa.debugf("removing %s\n", cfa.fal.addresses[fa_idx])
			cfa.fal.remove(cfa.fal.addresses[fa_idx])
		}
	}
}
func (cfa *custom_fancy_addresslist) debugf(format string, args ...interface{}) {
	if !*debug_custom_resolver {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[go-easyops] fancyaddresslist: %s", s)
}
func targets_to_fancyadr(targets []*registry.Target) []*FancyAddr {
	var res []*FancyAddr
	for _, t := range targets {
		if !hasApi(t.ApiType, registry.Apitype_grpc) {
			// ignore targets without apitype grpc
			continue
		}
		fa := &FancyAddr{state: connectivity.Ready, // docs say use CONNECTING here, but that never calls the picker nor the stateupdate. how does that work?
			addr:   fmt.Sprintf("%s:%d", t.IP, t.Port),
			Target: t,
		}
		res = append(res, fa)
	}

	return res
}

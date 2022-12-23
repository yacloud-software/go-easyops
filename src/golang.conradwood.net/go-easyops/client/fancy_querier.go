package client

import (
	"context"
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/standalone"
	"sync"
	"time"
)

var (
	query_max_age = flag.Duration("ge_registry_query_cache_duration", time.Duration(5)*time.Second, "effectively a limit on how often a given registry will be queried for a given service")
	keylock       sync.Map
)

type instance_cache_entry struct {
	lock      sync.Mutex
	res       []*pb.Target
	err       error
	refreshed time.Time
}

func queryServiceInstances(f_registry, serviceName string) ([]*pb.Target, error) {
	key := f_registry + "_" + serviceName
	ic, _ := keylock.LoadOrStore(key, &instance_cache_entry{})
	ice := ic.(*instance_cache_entry)
	if *dialer_debug {
		fmt.Printf("[go-easyops] Request to resolve service address \"%s\" via registry %s...\n", serviceName, f_registry)
	}
	ice.lock.Lock()
	defer ice.lock.Unlock()
	if (ice.res != nil || ice.err != nil) && time.Since(ice.refreshed) < *query_max_age {
		if *dialer_debug {
			fmt.Printf("[go-easyops] Answering request to Resolve service address \"%s\" via registry %s from cache\n", serviceName, f_registry)
		}
		return ice.res, ice.err
	}
	//	serviceName := f.target
	totalQueryCtr.With(prometheus.Labels{"servicename": serviceName}).Inc()
	if *dialer_debug {
		fmt.Printf("[go-easyops] Resolving service address \"%s\" via registry %s...\n", serviceName, f_registry)
	}
	request := &pb.V2GetTargetRequest{
		ApiType:     pb.Apitype_grpc,
		ServiceName: []string{serviceName},
		Partition:   "",
	}
	var err error
	var list *pb.V2GetTargetResponse
	ctx := context.Background()
	if cmdline.IsStandalone() {
		list, err = standalone.Registry_V2GetTarget(ctx, request)
	} else {
		regClient, xerr := getRegistryClient(f_registry)
		if xerr != nil {
			return nil, xerr
		}
		list, err = regClient.V2GetTarget(ctx, request)
	}
	// error getting stuff from registry
	if err != nil {
		if *dialer_debug {
			fmt.Printf("[go-easyops] error retrieving hosts for %s: %s\n", serviceName, err)
		}
		ice.err = err
		ice.refreshed = time.Now()
		return nil, err
	}
	ice.refreshed = time.Now()
	ice.res = list.Targets
	ice.err = nil
	return list.Targets, nil
}

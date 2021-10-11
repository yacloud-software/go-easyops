package client

import (
	"fmt"
	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type TargetList struct {
	reg      string
	svcname  string
	runOnce  bool
	tstarted bool
	targets  map[string]*TargetWithConnection
	rlock    sync.Mutex
}
type TargetWithConnection struct {
	target *registry.Target
	con    *grpc.ClientConn
}

func NewTargetList(servicename string) *TargetList {
	t := &TargetList{svcname: servicename, reg: cmdline.GetClientRegistryAddress()}
	t.targets = make(map[string]*TargetWithConnection)
	return t
}
func (t *TargetWithConnection) Connection() *grpc.ClientConn {
	return t.con
}
func (t *TargetWithConnection) String() string {
	return fmt.Sprintf("%s@%s:%d", t.target.ServiceName, t.target.IP, t.target.Port)
}

func (t *TargetList) Targets() []*TargetWithConnection {
	if !t.runOnce {
		err := t.refresh()
		if err != nil {
			fmt.Printf("[go-easyops] failed to refresh %s: %s\n", t.svcname, utils.ErrorString(err))
			return nil
		}
	}
	var res []*TargetWithConnection
	for _, v := range t.targets {
		if v.con == nil {
			continue
		}
		res = append(res, v)
	}
	return res
}

// returns a connection to exactly one target
func (t *TargetList) Connections() []*grpc.ClientConn {
	if !t.runOnce {
		err := t.refresh()
		if err != nil {
			fmt.Printf("[go-easyops] failed to refresh %s: %s\n", t.svcname, utils.ErrorString(err))
			return nil
		}
	}
	var res []*grpc.ClientConn
	for _, v := range t.targets {
		if v.con != nil {
			res = append(res, v.con)
		}
	}
	return res
}
func (t *TargetList) ByAddress(address string) []*TargetWithConnection {
	var res []*TargetWithConnection
	for _, v := range t.targets {
		if v.target.IP != address {
			continue
		}
		res = append(res, v)
	}
	return res
}
func (t *TargetList) refresh_loop() {
	for {
		time.Sleep(30 * time.Second)
		err := t.refresh()
		if err != nil {
			fmt.Printf("refresh failed: %s\n", utils.ErrorString(err))
		}
	}
}
func (t *TargetList) refresh() error {
	t.rlock.Lock()
	defer t.rlock.Unlock()
	if !t.tstarted {
		go t.refresh_loop()
		t.tstarted = true
	}
	rc := GetRegistryClient()
	treq := &registry.V2GetTargetRequest{ApiType: registry.Apitype_grpc, ServiceName: []string{t.svcname}}
	ctx := tokens.ContextWithToken()
	tr, err := rc.V2GetTarget(ctx, treq)
	if err != nil {
		return err
	}
	for _, target := range tr.Targets {
		t.add(target)
	}
	for _, tc := range t.targets {
		target := tc.target
		url := fmt.Sprintf("%s:%d", target.IP, target.Port)
		found := false
		for _, yestarget := range tr.Targets {
			yesurl := fmt.Sprintf("%s:%d", yestarget.IP, yestarget.Port)
			if yesurl == url {
				found = true
				break
			}
		}
		if !found {
			t.remove(target)
		}
	}
	t.runOnce = true
	return nil
}

// if target does not exist, it's a noop. otherwise close connection and remove
func (t *TargetList) remove(target *registry.Target) {
	url := fmt.Sprintf("%s:%d", target.IP, target.Port)
	tc, k := t.targets[url]
	if !k {
		return
	}
	if tc.con != nil {
		tc.con.Close()
		tc.con = nil
	}
	delete(t.targets, url)
}

// if target exists, it's a noop
func (t *TargetList) add(target *registry.Target) {
	url := fmt.Sprintf("%s:%d", target.IP, target.Port)
	_, k := t.targets[url]
	if k {
		return
	}
	tc := &TargetWithConnection{target: target}
	t.targets[url] = tc
	c, err := ConnectWithIP(url)
	if err != nil {
		fmt.Printf("[go-easyops] could not connect to %s: %s\n", url, utils.ErrorString(err))
		return
	}
	tc.con = c
}

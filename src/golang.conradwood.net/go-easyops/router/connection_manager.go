package router

import (
	"context"
	"fmt"
	"sync"

	"golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/client"
	"google.golang.org/grpc"
)

type ConnectionManager struct {
	use_all     bool // don't follow registry recommendations, just all registered ones
	one_per_ip  bool // but at most one per ip
	servicename string
}
type ConnectionTarget struct {
	lock       sync.Mutex
	ip         string
	port       uint32
	connection *Connection
}
type Connection struct {
	lock    sync.Mutex
	address string
	gcon    *grpc.ClientConn
}

func NewConnectionManager(servicename string) *ConnectionManager {
	return &ConnectionManager{use_all: true, one_per_ip: true, servicename: servicename}
}
func (cm *ConnectionManager) ServiceName() string {
	return cm.servicename
}
func (cm *ConnectionManager) AllowMultipleInstancesPerIP() {
	cm.one_per_ip = false
}
func (cm *ConnectionManager) GetCurrentTargets(ctx context.Context) []*ConnectionTarget {
	var res []*ConnectionTarget
	if cm.use_all {
		res = cm.getCurrentRegistrationsAsTargets(ctx)
		res = append(res, cm.getCurrentTargets(ctx)...)
	} else {
		res = cm.getCurrentTargets(ctx)
	}
	res = cm.filter(res)
	return res
}
func (cm *ConnectionManager) getCurrentRegistrationsAsTargets(ctx context.Context) []*ConnectionTarget {
	req := &registry.V2ListRequest{NameMatch: cm.servicename}
	targetlist, err := client.GetRegistryClient().ListRegistrations(ctx, req)
	if err != nil {
		fmt.Printf("Failed to get registrations for \"%s\": %s\n", req.NameMatch, err)
		return nil
	}
	var res []*ConnectionTarget
	for _, regi := range targetlist.Registrations {
		if !regi.Targetable {
			continue
		}
		t := regi.Target
		found := false
		for _, at := range t.ApiType {
			if at == registry.Apitype_grpc {
				found = true
				break
			}
		}
		if !found {
			// not a grpc targettype
			continue
		}
		ct := &ConnectionTarget{ip: t.IP, port: t.Port}
		res = append(res, ct)
	}
	return res
}
func (cm *ConnectionManager) getCurrentTargets(ctx context.Context) []*ConnectionTarget {
	req := &registry.V2GetTargetRequest{
		ServiceName: []string{cm.servicename},
		ApiType:     registry.Apitype_grpc,
	}
	targetlist, err := client.GetRegistryClient().V2GetTarget(ctx, req)
	if err != nil {
		fmt.Printf("Failed to get targets for \"%s\": %s\n", req.ServiceName, err)
		return nil
	}
	var res []*ConnectionTarget
	for _, t := range targetlist.Targets {
		ct := &ConnectionTarget{ip: t.IP, port: t.Port}
		res = append(res, ct)
	}
	return res
}
func (ct *ConnectionTarget) Address() string {
	return fmt.Sprintf("%s:%d", ct.ip, ct.port)
}

func (ct *ConnectionTarget) Connection() (*Connection, error) {
	ct.lock.Lock()
	defer ct.lock.Unlock()
	if ct.connection != nil {
		return ct.connection, nil
	}
	c := &Connection{address: ct.Address()}
	ct.connection = c
	return ct.connection, nil
}
func (c *Connection) GRPCConnection() (*grpc.ClientConn, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.gcon != nil {
		return c.gcon, nil
	}
	gcon, err := client.ConnectWithIP(c.address)
	if err != nil {
		return nil, err
	}
	c.gcon = gcon
	return c.gcon, nil
}
func (c *Connection) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.gcon != nil {
		c.gcon.Close()
		c.gcon = nil
	}
}

func (cm *ConnectionManager) filter(input []*ConnectionTarget) []*ConnectionTarget {
	ipmap := make(map[string]*ConnectionTarget)
	for _, ct := range input {
		key := ct.ip
		if !cm.one_per_ip {
			key = fmt.Sprintf("%s:%d", ct.ip, ct.port)
		}

		_, fd := ipmap[key]
		if fd {
			continue
		}
		ipmap[key] = ct
	}
	var res []*ConnectionTarget
	for _, v := range ipmap {
		res = append(res, v)
	}
	return res
}

func (cm *ConnectionManager) debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	prefix := fmt.Sprintf("[go-easyops router/cntmgr %s]", cm.servicename)
	txt := fmt.Sprintf(format, args...)
	fmt.Print(prefix + txt)
}
func (ct *ConnectionTarget) Close() {
	ct.lock.Lock()
	defer ct.lock.Unlock()
	if ct.connection == nil {
		return
	}
	ct.connection.Close()
}

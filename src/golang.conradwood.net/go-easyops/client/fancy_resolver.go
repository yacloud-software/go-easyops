package client

// this does the actual resolving with our registry

import (
	"context"
	//	"crypto/tls"
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RESOLVER_ATTRIBUTE_SERVICE_ADDRESS = "service_address"
)

var (
	query_for_proxies = flag.Bool("ge_support_proxies", true, "if true, supports routing via and to registrymultiplexer proxies")
	reglock           sync.Mutex
	proxyTargetLock   sync.Mutex
	proxiedTargets    = make(map[string]*ProxyTarget)      // serviceid -> proxytarget
	registryclients   = make(map[string]pb.RegistryClient) // map of "ip:port" -> registry
	resolv_chan       = make(chan int)
	resolvers         []*FancyResolver
	totalQueryCtr     = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_loadbalancer_registry_queries",
			Help: "counter incremented each time the loadbalancer queries the registry",
		},
		[]string{"servicename"},
	)
)

type ProxyTarget struct {
	Target    *pb.Target
	created   time.Time
	lastused  time.Time
	goingaway bool
	tcpConn   net.Conn
	tlsConn   net.Conn
}

func (p *ProxyTarget) key() string {
	return fmt.Sprintf("%s_%s_%d_%s_%s",
		p.Target.ServiceName,
		p.Target.IP,
		p.Target.Port,
		p.Target.RoutingInfo.GatewayID,
		p.Target.Partition,
	)
}

func init() {
	go resolver_thread()
	resolver.Register(&FancyResolverBuilder{})
	prometheus.MustRegister(totalQueryCtr)
}

type FancyResolverBuilder struct {
}

// this scheme matches the url
func (f *FancyResolverBuilder) Scheme() string {
	return "go-easyops"
}
func (f *FancyResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	if target.Authority == "" {
		panic("no target")
	}
	var registry string
	if !strings.Contains(target.Endpoint, "@") {
		panic(fmt.Sprintf("Invalid url - no registry in resolver. is \"%s\", missing @host:ip", target.Endpoint))
	}
	rs := strings.Split(target.Endpoint, "@")
	registry = rs[len(rs)-1]
	res := &FancyResolver{cc: cc, target: target.Authority, registry: registry}
	fancyPrintf(res, "fancy_resolver(): Request to build resolver for %#v\n", target)

	common.AddServiceName(res.target)
	resolvers = append(resolvers, res)
	resolv_chan <- 1
	return res, nil
}

type FancyResolver struct {
	registry                 string
	target                   string
	cc                       resolver.ClientConn
	noInstanceWarningPrinted bool
	instances                int
	lastScanned              time.Time
}

func (f *FancyResolver) ResolveNow(opts resolver.ResolveNowOptions) {
	fancyPrintf(f, "ResolveNow() on target %s with opts: %#v\n", f.target, opts)
	resolv_chan <- 1

}

func (f *FancyResolver) Close() {
	return
}

// called sync by the resolver_thread
func (f *FancyResolver) ActionResolve() {

	fancyPrintf(f, "fancy_resolver(): Resolving %s\n", f.target)
	regs, err := f.queryForInstances()
	if err != nil {
		fancyPrintf(f, "Error resolving: %s\n", err)
		f.cc.ReportError(err)
		return
	}
	f.instances = len(regs)
	f.updateCounters(len(regs))
	f.blockedWarning(len(regs))
	var ra []resolver.Address
	for _, a := range regs {
		rad := resolver.Address{
			ServerName: "go-easyops-server-name",
			Addr:       fmt.Sprintf("%s%s:%d", DIRECT_PREFIX, a.IP, a.Port),
			Attributes: attributes.New(RESOLVER_ATTRIBUTE_SERVICE_ADDRESS, a),
		}
		if a.RoutingInfo != nil && a.RoutingInfo.GatewayID != "" && *query_for_proxies {
			pt := &ProxyTarget{Target: a, created: time.Now()}
			proxiedTargets[pt.key()] = pt
			rad.Addr = fmt.Sprintf("%s%s", PROXY_PREFIX, pt.key())
		}

		ra = append(ra, rad)
		//	fancyPrintf(f,"fancy_resolver(): service \"%s\" on address: %s\n", r.target, a)
	}
	f.cc.UpdateState(resolver.State{Addresses: ra})
}

// update prometheus counters for got or not got instances
func (f *FancyResolver) updateCounters(adrcount int) {
	if adrcount == 0 {
		blockCtr.With(prometheus.Labels{"servicename": f.target}).Inc()
	}
	// done in "queryForActiveInstances"
	//	totalQueryCtr.With(prometheus.Labels{"servicename": f.target}).Inc()

}

// print a warning (or cancellation of warning) if no instances are found for a service
func (f *FancyResolver) blockedWarning(adrcount int) {
	if adrcount == 0 && !f.noInstanceWarningPrinted {
		fmt.Printf("WARNING - no instances for \"%s\"\n", f.target)
		f.noInstanceWarningPrinted = true
	}
	if adrcount != 0 && f.noInstanceWarningPrinted {
		fmt.Printf("WARNING CANCEL - %d instances for \"%s\"\n", adrcount, f.target)
		f.noInstanceWarningPrinted = false
	}
}

/********************* this thread monitors the registry and provides regular updates ***********/
func resolver_thread() {
	interval := defaultInterval() // update sleep interval to match flag
	for {
		select {
		case _ = <-resolv_chan:
		//
		case <-time.After(interval):
			//
		}

		if len(resolvers) == 0 {
			continue
		}

		interval = defaultInterval() // update sleep interval in case flags change
		//	fancyPrintf(f,"fancy_resolver(): resolving...\n")
		for _, r := range resolvers {
			if r.instances != 0 && (time.Since(r.lastScanned) < defaultInterval()) {
				// don't scan resolver - it has been scanned recently
				continue
			}
			r.ActionResolve() // get potential instances from registry for this resolvers' target
			if r.instances == 0 {
				// if we have a resolver w/o backends, query that one more frequently
				interval = time.Duration(1) * time.Second
			}
		}
	}
}

func defaultInterval() time.Duration {
	return time.Duration(*normal_sleep_time) * time.Second
}

// get the ip:port listings from the registry for this service
func (f *FancyResolver) queryForInstances() ([]*pb.Target, error) {
	serviceName := f.target
	totalQueryCtr.With(prometheus.Labels{"servicename": serviceName}).Inc()
	if *dialer_debug {
		fmt.Printf("[go-easyops] Resolving service address \"%s\" via registry %s...\n", serviceName, f.registry)
	}
	regClient, err := getRegistryClient(f.registry)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	list, err := regClient.V2GetTarget(ctx, &pb.V2GetTargetRequest{
		ApiType:     pb.Apitype_grpc,
		ServiceName: []string{serviceName},
		Partition:   "",
	})
	//	list, err := regClient.ListRegistrations(ctx, &pb.V2ListRequest{NameMatch: serviceName})
	// error getting stuff from registry
	if err != nil {
		if *dialer_debug {
			fmt.Printf("[go-easyops] error retrieving hosts for %s: %s\n", serviceName, err)
		}
		return nil, err
	}
	return list.Targets, nil
}

/*
func hasGRPC(r *pb.Registration) bool {
	for _, a := range r.Target.ApiType {
		if a == pb.Apitype_grpc {
			return true
		}
	}
	return false
}
*/
func getRegistryClient(registryAddress string) (pb.RegistryClient, error) {
	r := registryclients[registryAddress]
	if r != nil {
		return r, nil
	}
	reglock.Lock()
	defer reglock.Unlock()

	// connect to registry
	//	fmt.Printf("[go-easyops] Connecting to registry at %s...\n", registryAddress)
	var err error
	// try to use tls first
	conn := withTLS(registryAddress)
	if conn == nil {
		conn, err = grpc.Dial(
			registryAddress,
			//grpc.WithTransportCredentials(GetClientCreds()),
			grpc.WithInsecure(),
			//			grpc.WithUnaryInterceptor(unaryClientInterceptor),
			//			grpc.WithStreamInterceptor(unaryStreamInterceptor),
			grpc.WithTimeout(time.Duration(CONST_CALL_TIMEOUT)*time.Second),
		)
	}
	if err != nil {
		fmt.Printf("Failed to connect to registry at %s: %s\n", registryAddress, utils.ErrorString(err))
		return nil, err
	}
	registryClient := pb.NewRegistryClient(conn)
	registryclients[registryAddress] = registryClient
	return registryClient, nil
}

// this is quite a hack. Through tribal knowledge we know that the registry
// exposes RPC as non tls. on "port+1" however it exposes it via TLS
// so we try to connect to that first.
func withTLS(address string) *grpc.ClientConn {
	xs := strings.Split(address, ":")
	if len(xs) < 2 {
		return nil
	}
	xx, err := strconv.Atoi(xs[1])
	if err != nil {
		fmt.Printf("weird registry, not a number \"%s\": %s\n", address, err)
		return nil
	}
	np := fmt.Sprintf("%s:%d", xs[0], xx+1)
	conn, err := grpc.Dial(
		np,
		grpc.WithTransportCredentials(GetClientCreds()),
		//			grpc.WithUnaryInterceptor(unaryClientInterceptor),
		//			grpc.WithStreamInterceptor(unaryStreamInterceptor),
		grpc.WithTimeout(time.Duration(CONST_CALL_TIMEOUT)*time.Second),
	)
	if err != nil {
		fmt.Printf("unable to dial registry with TLS: %s", err)
		return nil
	}
	return conn
}

func GetProxyTarget(ctx context.Context, serviceid string) (*ProxyTarget, error) {
	proxyTargetLock.Lock()
	defer proxyTargetLock.Unlock()
	pt := proxiedTargets[serviceid]
	if pt == nil {
		return nil, fmt.Errorf("Proxy ID %s is not known here", serviceid)
	}
	if pt.tlsConn == nil {
		tcs := fmt.Sprintf("%s:%d", pt.Target.IP, pt.Target.Port)
		scs := fmt.Sprintf("\"%s\"@%s", pt.Target.ServiceName, tcs)
		fmt.Printf("Dialing proxy-connection %s\n", scs)
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", tcs)
		if err != nil {
			fmt.Printf("Failed to connect %s: %s\n", scs, err)
			return nil, err
		}
		err = pt.Start(conn)
		if err != nil {
			conn.Close()
			fmt.Printf("Failed to start connection %s: %s\n", scs, err)
			return nil, err
		}
		//		fmt.Printf("[go-easyops] WARNING client requested serviceid \"%s\", which is not resolvable\n", serviceid)
		return pt, nil
	}
	pt.lastused = time.Now()
	return pt, nil
}

// send the initialisation sequence
func (p *ProxyTarget) Start(c net.Conn) error {
	//	tc := tls.Client(c, &tls.Config{InsecureSkipVerify: true})

	s, err := utils.Marshal(p.Target.RoutingInfo)
	if err != nil {
		return err
	}
	buf := []byte("C" + s + "\n")
	_, err = c.Write(buf)
	if err != nil {
		return err
	}
	//	p.tlsConn = tc
	p.tcpConn = c
	fmt.Printf("Started tcp connection for %s\n", p.Target.ServiceName)
	return err
}
func (f *FancyResolverBuilder) ServiceName() string {
	return "fancyresolverbuilder"
}
func (f *FancyResolver) ServiceName() string {
	return f.target
}

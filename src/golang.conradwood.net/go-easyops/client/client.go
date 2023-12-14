/*
This package facilitates making load-balanced, fail-over, authenticated gRPC calls to server. (it also provides shortcuts to objectstore get/put)

Typically, in the yacloud, a new client is constructed via the proto package. For example

	import "golang.conradwood.net/apis/getestservice"
	...
	getestservice.GetEchoClient()...

# Load balancing

The client will maintain a list of available targets for grpc calls. Each service has a unique list of targets. The list is periodically and aysnchronously updated by polling the registry.

If no targets are available, a call to this service will be blocked for some time, then fail. Once failed, all subsequent calls will be failed immediately until a target becomes available.
This allows for some basic recovery for circular service dependencies on boot. Whilst it is considered bad practice, it is a pattern commonly found and thus go-easyops attempts to make it work as well as can be expected. (Better to avoid circular dependencies altogether!!)

RPC calls are not retried - if they fail (for any reason) they will not be sent to a different server.

# Routing

This client implements several and distinct features to determine where to route rpc calls too.
  - Round-Robin for multiple targets
  - By User: to user specific services
  - By Context-Tag: arbitrary tags in the context

In the absense of both user specific services and context-tags, a simple round-robin strategy is implemented for multiple targets.

# Routing - by user

A service may be registered with a useraccount instead of a service account (the default if started by a user on the command line. Also see command line flag -ge_disable_user_token). The client determines the current useraccount from the context used to invoke the target. If a service running as the same user as is in the context, the client will route the rpc call to the service running as this user.

This is intended for debugging and "live" development. The user object is typically created and propagated and the edge of the system, for example at the webserver proxy. Thus, while developing or debugging any rpc server it is often useful to route some (and only some) calls into the development version. Developers should always start their work-in-progress servers under their own useraccount. Thus all calls the Developer makes are routed to their laptop. Subsequent calls go back into the infrastructure, but remain restricted to the useraccount, thus can be considered safe (subject to bugs in the backends of course).

Note: this is not intended for general production use. For various reasons, it is a really bad idea to fire up rpc servers specific to each user. Instead the rpc server should handle multiple users. The user's token for authentication is considered private, like a password.

# Routing - by Context-Tag

A context-tag is a key with a value .

A service may register itself with one or more tags. A Context may carry one or more tags. On each rpc call the list of targets is iterated through. Any service that has all tags with exactly matching values will be considered for routing, all others dismissed. If after, the first iteration, one or more services remain, those will be used for routing in a round-robin fashion.

If no exact match is found, the context tags are inspected for "FallbackToPlain" option (see ctx package). If set, the list of all services with exactly 0 tags will be used for round-robin. If not set the rpc will be failed with "no target available".

Note: Whilst context-tags are often quite useful, their use is generally discouraged, especially for large sets of servers. It is intented to be used for a standardized, quick and reasonably efficient means to route low-volume (~20/s or less) calls to remote rpc servers. In small setups this can often be useful to send information to remote clusters for remote-processing. (using a tag, for example, cluster=lhr, cluster=fra, cluster=lgw, cluster=cdn etc...)

# Routing - directly

One can bypass the fail-over and connection management altogether with functions such as ConnectWithIP. This is intented for circumstances where a standardised approach (round-robin/user/context) is unsuitable. Beware of dragons: This approach requires development of load-balancers, fail-over strategies, monitoring, recovery, start-up delays and many other features the default routing strategy implements. The complexity is often underestimated, but soon does become significant. (That is the point of the routing implementations above, really)

# Standalone operation

Standalone operation means that no other services are required to make rpc calls between a client and server. Whilst this is very limited (e.g. no load-balancing, no fail-over and NO AUTHENTICATION, it is useful for quickly testing a binary). Both, server and client must be started in standalone mode.

Also see command line parameter -ge_standalone and package server

# ObjectStore

Mostly because an extra package for objectstore seems overkill, this package also provides function to get/put objects into the objectstore.
*/
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/certificates"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
	got_client_creds = false
	client_creds     credentials.TransportCredentials

	cert      = []byte{1, 2, 3}
	errorList []*errorCache
	errorLock sync.Mutex
	reg       pb.RegistryClient
)

type errorCache struct {
	servicename string
	lastOccured time.Time
	lastPrinted time.Time
}

func GetRegistryClient() pb.RegistryClient {
	if reg == nil {
		reg = pb.NewRegistryClient(Connect("registry.Registry"))
	}
	return reg
}

// opens a tcp connection to a servicename
func DialTCPWrapper(serviceName string) (net.Conn, error) {
	if strings.Contains(serviceName, "/") {
		s := fmt.Sprintf("Error: The parameter for DialTCPWrapper needs a servicename. not a path. You passed in %s, which looks very much like a path. The \"old-style\" picoservices required a path at this function, but go-framework does not. Did you recently upgrade and did not upgrade a config?\n", serviceName)
		debug.PrintStack()
		return nil, errors.New(s)
	}
	if reg == nil {
		reg = pb.NewRegistryClient(Connect("registry.Registry"))
	}
	ctx := getContext()
	//ctx := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	targets, err := reg.GetTarget(ctx, &pb.GetTargetRequest{Name: serviceName, ApiType: pb.Apitype_tcp})
	if err != nil {
		return nil, err
	}
	list := targets.Service
	if len(list) == 0 {
		return nil, fmt.Errorf("No tcp connection by service %s", serviceName)
	}
	if len(list[0].Location.Address) == 0 {
		return nil, fmt.Errorf("No tcp location found for name %s - is it running?", serviceName)
	}
	adr := fmt.Sprintf("%s:%d", list[0].Location.Address[0].Host, list[0].Location.Address[0].Port)
	conn, err := net.Dial("tcp", adr)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func hasApi(ar []pb.Apitype, lf pb.Apitype) bool {
	for _, a := range ar {
		if a == lf {
			return true
		}
	}
	return false
}

// get the Client Credentials we use to connect to other RPCs
func GetClientCreds() credentials.TransportCredentials {
	if got_client_creds {
		return client_creds
	}
	roots := x509.NewCertPool()

	frontendCert := certificates.Certificate()

	roots.AppendCertsFromPEM(frontendCert)
	ImCert := certificates.Ca() //ioutil.ReadFile(*clientca)
	roots.AppendCertsFromPEM(ImCert)

	pk := certificates.Privatekey()

	cert, err := tls.X509KeyPair(frontendCert, pk)
	//	cert, err := tls.LoadX509KeyPair(*clientcrt, *clientkey)
	if err != nil {
		fmt.Printf("Failed to create client certificates: %s\n", err)
		fmt.Printf("key:\n%s\n", string(pk))
		return nil
	}
	// verify using the server address in the certificte, not the ACTUAL address
	creds := credentials.NewTLS(&tls.Config{
		ServerName:         certificates.ServerName(),
		Certificates:       []tls.Certificate{cert},
		RootCAs:            roots,
		InsecureSkipVerify: true,
	})
	client_creds = creds
	got_client_creds = true
	return creds
}

func getErrorCacheByName(name string) *errorCache {
	errorLock.Lock()
	defer errorLock.Unlock()
	for _, ec := range errorList {
		if ec.servicename == name {
			return ec
		}
	}
	ec := &errorCache{servicename: name,
		lastOccured: time.Now(),
	}
	errorList = append(errorList, ec)
	return ec
}

func printError(path string, msg string) {
	e := getErrorCacheByName(path)
	if e == nil {
		fmt.Println(msg)
		return
	}
	if !e.needsPrinting() {
		return
	}
	fmt.Println(msg)
}

// returns true if this needs printing
// resets counter if it returns true
func (e *errorCache) needsPrinting() bool {
	now := time.Now()
	if now.Sub(e.lastPrinted) < (time.Duration(5) * time.Minute) {
		return false
	}
	e.lastPrinted = now
	return false
}

// given an inbound Context (e.g. in an RPC call) this creates a new outbound
// context suitable to call other servers
// we keep user & org intact.
// we add/override 'service' token (with our token)
func DIS_OutboundContext(inbound context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(inbound)
	if !ok {
		fmt.Printf("[go-easyops] WARNING -> inbound context has no metadata authentication\n")
		md = metadata.Pairs()
	}
	md = md.Copy()
	return metadata.NewOutgoingContext(inbound, md)

}

// get instances for a service currently being connected to (that is, it will return 0 for services which have not been dialled (yet)
func GetConnectedInstanceCount(servicelookupid string) int {
	for _, r := range resolvers {
		fmt.Printf("lookup=%s, target=%s\n", servicelookupid, r.target)
	}
	return 0

}

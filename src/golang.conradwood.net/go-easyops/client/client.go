package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/certificates"
	"golang.conradwood.net/go-easyops/tokens"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
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
	targets, err := reg.GetTarget(tokens.ContextWithToken(), &pb.GetTargetRequest{Name: serviceName, ApiType: pb.Apitype_tcp})
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

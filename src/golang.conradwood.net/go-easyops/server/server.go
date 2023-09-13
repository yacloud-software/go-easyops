package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	pm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	au "golang.conradwood.net/apis/auth"
	cm "golang.conradwood.net/apis/common"
	echo "golang.conradwood.net/apis/echoservice"
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/auth"
	ar "golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/certificates"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	easyhttp "golang.conradwood.net/go-easyops/http"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/standalone"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	COOKIE_NAME = "Auth-Token"
)

var (
	auto_kill                      = flag.Bool("ge_autokill_instance_on_port", false, "if true, kill an instance on that grpc port before starting service")
	never_register_service_as_user = flag.Bool("ge_never_register_service_as_user", false, "if true, do not register service as user, even if it is run locally with a user token")
	reg_tags                       = flag.String("ge_routing_tags", "", "comma seperated list of key-value pairs. For example -tags=foo=bar,foobar=true")
	debug_internal_serve           = flag.Bool("ge_debug_internal_server", false, "debug the server @ https://.../internal/... (serving metrics amongst other things)")
	debug_rpc_serve                = flag.Bool("ge_debug_rpc_server", false, "debug the grpc server ")
	deployDescriptor               = flag.String("ge_deployment_descriptor", "", "The deployment path by which other programs can refer to this deployment. expected is: a path of the format: \"V1:namespace/groupname/repository/buildid\"")
	register_refresh               = flag.Int("ge_register_refresh", 10, "registration refresh interval in `seconds`")
	serverDefs                     = make(map[string]*serverDef)
	knownServices                  []*serverDef // all services, even not known ones
	stopped                        bool
	ticker                         *time.Ticker
	promHandler                    http.Handler
	//promReg         = prometheus.NewRegistry()
	stdMetrics        = NewServerMetrics()
	startedPreviously = false
	starterLock       sync.Mutex
	rgclient          pb.RegistryClient
	startup_complete  = false
)

type UserCache struct {
	UserID  string
	created time.Time
}

type Register func(server *grpc.Server) error

// serverdef interface
type Server interface {
	AddTag(key, value string)
}

// no longer exported - please use NewServerDef instead
type serverDef struct {
	callback    func() // called if/when server started up successfully
	Port        int
	Certificate []byte
	Key         []byte
	CA          []byte
	Register    Register
	// set to true if this server does NOT require authentication (default: it does need authentication)
	NoAuth bool
	// set to false if this service should not register with the registry initially
	RegisterService bool
	name            string
	types           []pb.Apitype
	registered_id   string
	DeployPath      string
	serviceID       uint64
	asUser          *au.SignedUser // if we're running as a user rather than a server this is the account
	tags            map[string]string
	ErrorHandler    func(ctx context.Context, function_name string, err error)
	local_service   *au.SignedUser // the local service account
	service_user_id string         // the serviceaccount userid
	public          bool
}

func init() {
	if cmdline.IsStandalone() {
		return
	}
	// start period re-registration
	ticker = time.NewTicker(time.Duration(*register_refresh) * time.Second)
	go func() {
		for _ = range ticker.C {
			reRegister()
		}
	}()
}
func (s *serverDef) DontRegister() {
	s.RegisterService = false
}
func (s *serverDef) SetPublic() {
	s.public = true
}

/*
set a callback that is called AFTER grpc server started successfully
*/
func (s *serverDef) SetOnStartupCallback(f func()) {
	s.callback = f
}

// add a routing tag to a serverdef
func (s *serverDef) AddTag(key, value string) {
	s.tags[key] = value
}
func (s *serverDef) toString() string {
	return fmt.Sprintf("Port #%d: %s (%v)", s.Port, s.name, s.types)
}

func NewTCPServerDef(name string) *serverDef {
	sd := NewServerDef()
	sd.tags = make(map[string]string)
	sd.types = sd.types[:0]
	sd.types = append(sd.types, pb.Apitype_tcp)
	sd.name = name
	return sd
}

func NewHTMLServerDef(name string) *serverDef {
	sd := NewServerDef()
	sd.tags = make(map[string]string)
	sd.types = sd.types[:0]
	sd.types = append(sd.types, pb.Apitype_html)
	sd.name = name
	return sd
}

func NewServerDef() *serverDef {
	res := &serverDef{}
	res.tags = make(map[string]string)
	res.registered_id = ""
	/*
		res.Key = Privatekey
		res.Certificate = Certificate
		res.CA = Ca
	*/
	res.DeployPath = deploymentPath()
	res.types = append(res.types, pb.Apitype_status)
	res.types = append(res.types, pb.Apitype_grpc)
	res.RegisterService = true
	return res
}
func deploymentPath() string {
	if *deployDescriptor != "" {
		return (*deployDescriptor)[3:]
	}
	return ""
}

func stopping() {
	starterLock.Lock()
	if stopped {
		starterLock.Unlock()
		return
	}
	stopped = true
	starterLock.Unlock()
	pp.ProfilingStop()
	fancyPrintf("Server shutdown - deregistering services\n")

	c := client.GetRegistryClient()
	/*
		opts := []grpc.DialOption{grpc.WithInsecure()}
		rconn, err := grpc.Dial(cmdline.GetRegistryAddress(), opts...)
		if err != nil {
			fancyPrintf("failed to dial registry server: %v", err)
			return
		}
		defer rconn.Close()
		c = pb.NewRegistryClient(rconn)
	*/
	// value is a serverdef
	for _, sd := range knownServices {
		fancyPrintf("Deregistering Service \"%s\"\n", sd.toString())
		ctx := context_Background()
		ctx, _ = context.WithTimeout(ctx, time.Duration(2)*time.Second) // don't hang on shutdown

		//		ctx := authremote.Context()
		_, err := c.V2DeregisterService(ctx, &pb.DeregisterServiceRequest{ProcessID: sd.registered_id})
		if err != nil {
			fancyPrintf("Failed to deregister Service \"%s\": %s\n", sd.toString(), err)
		}
	}
}

func addTags(sd *serverDef) {
	if *reg_tags == "" {
		return
	}
	vals := strings.Split(*reg_tags, ",")
	for _, v := range vals {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			s := fmt.Sprintf("Invalid keyvalue tag: \"%s\" - it splits into %d parts instead of 2\n", v, len(kv))
			panic(s)
		}
		tk := kv[0]
		tv := kv[1]
		fmt.Printf("Adding tag \"%s\" with value \"%s\"\n", tk, tv)
		sd.AddTag(tk, tv)
	}

}

// this is our typical gRPC server startup
// it sets ourselves up with our own certificates
// which is set for THIS SERVER, so installed/maintained
// together with the server (rather than as part of this software)
// it also configures the rpc server to expect a token to identify
// the user in the rpc metadata call
func ServerStartup(def *serverDef) error {
	if *auto_kill {
		ht := easyhttp.NewDirectClient()
		hr := ht.Get(fmt.Sprintf("https://localhost:%d/internal/pleaseshutdown", def.Port))
		if hr.IsSuccess() {
			for {
				ht := easyhttp.NewDirectClient()
				hr := ht.Get(fmt.Sprintf("https://localhost:%d/internal/pleaseshutdown", def.Port))
				if hr.IsSuccess() {
					break
				}
				time.Sleep(time.Duration(300) * time.Millisecond)
			}
		}
	}
	addTags(def)
	go client.GetSignatureFromAuth() // init pubkey
	go error_handler_startup()
	var tk string
	started := time.Now()
	for {
		if client.GotSig() {
			break
		}
		if time.Since(started) > time.Duration(3)*time.Second {
			fmt.Printf("[go-easyops] WARNING could not retrieve signature in time\n")
			break
		}
	}
	tokname := ""
	tokname = "service"
	tkservice := tokens.GetServiceTokenParameter()
	var u *au.User
	if !cmdline.IsStandalone() {
		tk = tkservice
		if !cmdline.Datacenter() {
			tks := tokens.GetUserTokenParameter()
			if tks != "" {
				tokname = "user"
				tk = tks
			}
		}
		var su *au.SignedUser
		if !def.NoAuth {
			if tk == "" {
				fancyPrintf("*********** AUTHENTICATION CONFIGURATION ERROR ******************\n")
				fancyPrintf("Cannot connect to a server without %s token.\n", tokname)
				//os.Exit(10)
			}
			su = ar.SignedGetByToken(context_Background(), tk)
			if su == nil {
				fancyPrintf("*********** AUTHENTICATION CONFIGURATION ERROR ******************\n")
				fancyPrintf("The authentication %s token is not valid.\n", tokname)
				fancyPrintf("Token: \"%s\"\n", tk)
				//os.Exit(10)
			}
			u = common.VerifySignedUser(su)

		}
		if u != nil {
			if u.ServiceAccount {
				def.local_service = su
			} else {
				if *never_register_service_as_user {
					fancyPrintf("NOT Registering as a user-specific service (disabled by commandline)\n")
				} else {
					fancyPrintf("Registering as a user-specific service, because it is running as:\n")
					auth.PrintUser(u)
					def.asUser = su
				}
			}
		}
	}
	startOnce()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopping()
		os.Exit(0)
	}()
	stopped = false
	defer stopping()
	listenAddr := fmt.Sprintf(":%d", def.Port)
	s := ""
	if u != nil {
		def.service_user_id = u.ID
		s = fmt.Sprintf(" #%s [%s]", u.ID, auth.Description(u))
	}
	fancyPrintf("Starting server%s on %s\n", s, listenAddr)

	if def.tags != nil && len(def.tags) > 0 {
		fancyPrintf("Routing tags: %v\n", def.tags)
	}

	BackendCert := certificates.Certificate()
	BackendKey := certificates.Privatekey()
	ImCert := certificates.Ca()
	cert, err := tls.X509KeyPair(BackendCert, BackendKey)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v\n", err)
	}
	roots := x509.NewCertPool()
	FrontendCert := certificates.Certificate()
	roots.AppendCertsFromPEM(FrontendCert)
	roots.AppendCertsFromPEM(ImCert)

	creds := credentials.NewServerTLSFromCert(&cert)
	var grpcServer *grpc.Server
	// Create the gRPC server with the credentials
	grpcServer = grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(def.UnaryAuthInterceptor),
		grpc.StreamInterceptor(def.StreamAuthInterceptor),
	)

	grpc.EnableTracing = true
	// callback to the callers' specific intialisation
	// (set by the caller of this function)
	if def.Register != nil {
		def.Register(grpcServer)
	}
	if err != nil {
		fancyPrintf("Serverstartup: failed to register service on startup: %s\n", err)
		return fmt.Errorf("grpc register error: %s", err)
	}
	if len(grpcServer.GetServiceInfo()) > 1 {
		return fmt.Errorf("cannot register multiple(%d) names", len(grpcServer.GetServiceInfo()))
	}
	if def.name == "" {
		for name, _ := range grpcServer.GetServiceInfo() {
			def.name = name
		}
	}
	if def.name == "" {
		fmt.Println("Got no server name!")
		return errors.New("Missing servername")
	}

	serverDefs[def.name] = def
	common.AddExportedServiceName(def.name)

	if def.RegisterService {
		fancyPrintf("Adding service %s to registry...\n", def.name)
		AddRegistry(def)
	}
	// something odd?
	if !def.public {
		reflection.Register(grpcServer)
	}
	// Serve and Listen
	// Blocking call!
	err = startHttpServe(def, grpcServer)

	// Create the channel to listen on
	// I don't think this is ever called!
	fancyPrintf("INTERNAL BUG - we should have never, ever arrived here\n")
	os.Exit(3)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("could not listen on %s: %s", listenAddr, err)
	}
	fancyPrintf("Starting service %s...\n", def.name)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("grpc serve error: %s", err)
	}
	return nil
}

func startHttpServe(sd *serverDef, grpcServer *grpc.Server) error {
	mux := http.NewServeMux()
	if !sd.public {
		mux.HandleFunc("/internal/service-info/", func(w http.ResponseWriter, req *http.Request) {
			serveServiceInfo(w, req, sd)
		})
		mux.HandleFunc("/internal/pleaseshutdown", func(w http.ResponseWriter, req *http.Request) {
			pleaseShutdown(w, req, grpcServer)
		})
		mux.HandleFunc("/internal/health", func(w http.ResponseWriter, req *http.Request) {
			healthzHandler(w, req, sd)
		})
		mux.HandleFunc("/internal/help", func(w http.ResponseWriter, req *http.Request) {
			helpHandler(w, req, sd)
		})
		mux.HandleFunc("/internal/clearcache", func(w http.ResponseWriter, req *http.Request) {
			clearCacheHandler(w, req)
		})
		mux.HandleFunc("/internal/parameters", func(w http.ResponseWriter, req *http.Request) {
			paraHandler(w, req, sd)
		})

		nm, _ := prometheus.NonstandMetricNames(pm.DefaultRegisterer.(*pm.Registry))
		if len(nm) > 0 {
			for _, n := range nm {
				fmt.Printf("Reg: \"%s\"\n", n)
			}
			panic("something registered outside go-easyops and will not be exposed")
		}
		gatherer := prometheus.GetGatherer()
		h := promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{})
		mux.Handle("/internal/service-info/metrics", h)
		//	mux.Handle("/internal/service-info/metrics", promhttp.Handler())
	}
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", sd.Port))
	if err != nil {
		panic(err)
	}

	BackendCert := certificates.Certificate()
	BackendKey := certificates.Privatekey()
	cert, err := tls.X509KeyPair(BackendCert, BackendKey)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", sd.Port),
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			NextProtos:         []string{"h2"},
			InsecureSkipVerify: true,
		},
	}

	fancyPrintf("grpc on port: %d\n", sd.Port)
	go callback_attempt(sd)
	startup_complete = true
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	fancyPrintf("Serve failed: %v\n", err)
	return err
}

// attempt to http call into the server to trigger server_started callback
func callback_attempt(sd *serverDef) {
	url := fmt.Sprintf("https://localhost:%d/internal/health", sd.Port)
	for {
		//fmt.Printf("Testing %s\n", url)
		hr := easyhttp.NewDirectClient().Get(url)
		if hr.Error() == nil {
			break
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	fmt.Printf("[go-easyops] Server started on port %d\n", sd.Port)
	if sd.callback != nil {
		sd.callback()
	}
}

// this function is called by http and works out wether it's a grpc or http-serve request
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/internal/debug") {
			if *debug_internal_serve {
				fancyPrintf("Serving debug path %s\n", path)
			}
			debugHandler(w, r)
		} else if strings.HasPrefix(path, "/internal/clearcache") {
			clearCacheHandler(w, r)
		} else if strings.HasPrefix(path, "/internal/") {
			if *debug_internal_serve {
				fancyPrintf("Serving path %s\n", path)
			}
			otherHandler.ServeHTTP(w, r)
		} else {
			grpcServer.ServeHTTP(w, r)
		}
	})
}

// mostly for autodeployer
func UnregisterPortRegistry(port []int) error {
	var err error
	client := client.GetRegistryClient()
	/*
		opts := []grpc.DialOption{grpc.WithInsecure()}
		conn, err := grpc.Dial(cmdline.GetRegistryAddress(), opts...)
		if err != nil {
			fancyPrintf("failed to dial registry server: %v", err)
			return err
		}
		defer conn.Close()
		client := pb.NewRegistryClient(conn)
	*/
	var ps []int32
	for _, p := range port {
		ps = append(ps, int32(p))
	}
	psr := pb.ProcessShutdownRequest{Port: ps}
	_, err = client.InformProcessShutdown(context_Background(), &psr)
	return err
}

func find(port int, name string) *serverDef {
	for _, sd := range knownServices {
		if sd.Port == port && sd.name == name {
			return sd
		}
	}
	return nil
}

func AddRegistry(sd *serverDef) (string, error) {
	if find(sd.Port, sd.name) == nil {
		knownServices = append(knownServices, sd)
	}

	req := pb.ServiceLocation{}
	req.Service = &pb.ServiceDescription{}
	req.Service.Name = sd.name
	req.Service.Path = sd.DeployPath
	sa := &pb.ServiceAddress{Port: int32(sd.Port)}
	req.Address = []*pb.ServiceAddress{sa}

	rsr := &pb.RegisterServiceRequest{
		ProcessID:   cmdline.GetInstanceID(),
		Port:        uint32(sd.Port),
		ApiType:     sd.types,
		ServiceName: sd.name,
		Pid:         cmdline.GetPid(),
		RoutingInfo: &pb.RoutingInfo{},
		UserID:      sd.service_user_id,
	}
	if sd.asUser != nil {
		rsr.RoutingInfo.RunningAs = common.VerifySignedUser(sd.asUser)
	}
	if sd.tags != nil {
		rsr.RoutingInfo.Tags = sd.tags
	}
	if cmdline.IsStandalone() {
		return standalone.RegisterService(rsr)
	}
	if rgclient == nil {
		rgclient = client.GetRegistryClient()
	}
	resp, err := rgclient.V2RegisterService(context_Background(), rsr)
	if err != nil {
		fancyPrintf("RegisterService(%s) failed: %s\n", req.Service.Name, err)
		return "", err
	}
	if resp == nil {
		fmt.Println("Registration failed with no error provided.")
	}
	sd.registered_id = rsr.ProcessID
	//	fancyPrintf("Response to register service: %v\n", resp)
	//	fancyPrintf("Registered: %s\n", sd.registered_id)
	return sd.registered_id, nil
}

func reRegister() {
	// register any that are not yet registered
	for _, sd := range knownServices {
		AddRegistry(sd)
	}
}

func getServerDefByName(name string) *serverDef {
	return serverDefs[name]
}
func MethodNameFromUnaryInfo(info *grpc.UnaryServerInfo) string {
	full := info.FullMethod
	if full[0] == '/' {
		full = full[1:]
	}
	ns := strings.SplitN(full, "/", 2)
	if len(ns) < 2 {
		return ""
	}
	res := ns[1]
	if res[0] == '/' {
		res = res[1:]
	}
	return ns[1]
}
func ServiceNameFromUnaryInfo(info *grpc.UnaryServerInfo) string {
	full := info.FullMethod
	if full[0] == '/' {
		full = full[1:]
	}
	ns := strings.SplitN(full, "/", 2)
	return ns[0]
}

func targetName(name string) string {
	x := strings.Split(name, ".")
	return x[0]
}

func isInternalService(name string) bool {
	if name == "grpc.reflection.v1alpha.ServerReflection" {
		return true
	}
	return false
}

func startOnce() {
	starterLock.Lock()
	if startedPreviously {
		starterLock.Unlock()
		return
	}
	startedPreviously = true
	starterLock.Unlock()
	pp.ProfilingCheckStart()
}

/***************************************************
* convenience function to register stuff with the registry
* useful to register long-running clients, for example
* this allows for metrics to be exposed and scraped automatically
* uses a default RPC
***************************************************/
func StartFakeService(name string) {
	port, err := getFreePort()
	if err != nil {
		s := fmt.Sprintf("Failed to get a free port: %s", err)
		fmt.Println(s)
		panic(s)
	}
	sd := NewServerDef()
	sd.Port = port
	sd.Register = Register(
		func(server *grpc.Server) error {
			e := new(echoServer)
			echo.RegisterEchoServiceServer(server, e)
			return nil
		},
	)
	sd.name = name
	go ServerStartup(sd)
}

type echoServer struct{}

func (e *echoServer) Ping(ctx context.Context, req *cm.Void) (*echo.PingResponse, error) {
	fancyPrintf("I was pinged\n")
	resp := &echo.PingResponse{Response: "goeasyops-server"}
	return resp, nil
}

// ugly race-condition-hack to find a free port
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func incFailure(service string, method string, err error) {
	status := status.Convert(err)
	var code codes.Code
	if status != nil {
		code = status.Code()
	}
	grpc_failed_requests.With(prometheus.Labels{"method": method, "servicename": service, "grpccode": fmt.Sprintf("%d", code)}).Inc()
}

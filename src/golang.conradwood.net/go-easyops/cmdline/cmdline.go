/*
provides a standard configuration mechanism for clients and servers on the commandline

go-easyops itself can be configured through environment variables and command line parameters. any go-easyops program
includes standard command line parameters. some are described below

# -h

prints out build information and command line parameters for the application

# -X

prints out build information and command line parameters to the behaviour of go-easyops

# Environment Variables

-h and -X also print environment variables and a short help text for each. Application developers are encouraged to use
this package to manage environment variables. a typical example

	var (
	  mytext = cmdline.ENV("MYTEXT","specifies the text to display")
	)
	func main() {
	  fmt.Println(mytext.Value())
	}

# Config Files

a config file (typically /tmp/goeasyops.config) provides optional and initial configuration for go-easyops. This is intented to configure developer machines on-the-fly for access to a different cloud and cluster. For example, based on the current path, a git repository url may be used to configure a specific and matching registry. (config file syntax is in yaml, see goeasyops proto, protobuf "Config")
*/
package cmdline

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	pb "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/appinfo"
	"golang.conradwood.net/go-easyops/common"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_FILE      = "/tmp/goeasyops.config"
	REGISTRY_DEFAULT = "localhost:5000"
)

var (
	debug_rpc_client = flag.Bool("ge_debug_rpc_client", false, "set to true to debug remote invokations")
	debug_rpc_serve  = flag.Bool("ge_debug_rpc_server", false, "debug the grpc server ")
	default_timeout  = flag.Duration("ge_ctx_deadline", time.Duration(10)*time.Second, "the default timeout for contexts. do not change in production")
	reg_env          = ENV("GE_REGISTRY", "default registry address")
	e_ctx            = ENV("GE_CTX", "a serialised context to use when creating new ones")
	config           *pb.Config
	// annoyingly, not all go-easyops flags start with ge_
	internal_flag_names   = []string{"token", "registry", "registry_resolver", "AD_started_by_auto_deployer", "X"}
	debug_auth            = flag.Bool("ge_debug_auth", false, "debug auth stuff")
	debug_sig             = flag.Bool("ge_debug_signature", false, "debug signature stuff")
	mlock                 sync.Mutex
	running_in_datacenter = flag.Bool("AD_started_by_auto_deployer", false, "the autodeployer sets this to true to modify the behaviour to make it suitable for general-availability services in the datacenter")

	registry          = flag.String("registry", REGISTRY_DEFAULT, "address of the registry server. This is used for registration as well as resolving unless -registry_resolver is set, in which case this is only used for registration")
	registry_resolver = flag.String("registry_resolver", "", "address of the registry server (for lookups)")
	instance_id       = flag.String("ge_instance_id", "", "autodeployers internal instance id. We may use this to get information about ourselves")
	ext_help          = flag.Bool("X", false, "extended help")
	XXdoappinfo       = ImmediatePara("ge_info", "print application build number", doappinfo)
	print_easyops     = false
	manreg            = ""
	stdalone          = flag.Bool("ge_standalone", false, "if true, do not use a registry, just run stuff standlone")
	//	context_with_builder   = flag.Bool("ge_context_with_builder", true, "a new (experimental) context messaging method")
	context_build_version  = flag.Int("ge_context_builder_version", 2, "the version to create by the context builder (0=do not use context builder)")
	overridden_env_context = ""
	enabled_experiments    = flag.String("ge_enable_experiments", "", "a comma delimited set of names of experiments to enable with this context")
	debug_ctx              = flag.Bool("ge_debug_context", false, "if true debug context stuff")
)

// in the init function we have not yet defined all the flags
// each init() is called in the order of import statements, thus packages imported AFTER this package
// won't have their flags initialized yet
// I have not found a good way of being triggered once flags are parsed, thus we use a timer in the hope that it will work well enough
func init() {
	flag.Usage = PrintUsage
	for _, o := range os.Args {
		if o == "-X" {
			go print_late_usage()
		}
	}

	// read a potential config file
	err := readConfig(CONFIG_FILE)
	if err != nil {
		os.Exit(10)
	}
}
func readConfig(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		return nil // if file does not exist, it's not an error
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("[go-easyops] failed to read file %s: %s\n", filename, err)
		return err
	}
	cfg := &pb.Config{}
	err = yaml.UnmarshalStrict(b, cfg)
	if err != nil {
		fmt.Printf("[go-easyops] invalid file %s: %s\n", filename, err)
		return err
	}
	config = cfg
	return nil
}

// if we have a -X argument we will print extended usage AFTER flags are parsed.
// we know flags areg parsed if ext_help flag (-X) turns true (timeout after 5 secs)
func print_late_usage() {
	fmt.Printf("[go-easyops] Printing extended help after flag.Parse() was called...\n")
	st := time.Now()
	for *ext_help == false {
		if time.Since(st) > time.Duration(5)*time.Second {
			break
		}
	}
	print_easyops = true
	PrintUsage()
	os.Exit(0)
}

func PrintUsage() {
	fmt.Fprintf(os.Stdout, "(C) Conrad Wood.\n")
	fmt.Fprintf(os.Stdout, "  Go-Easyops version          : %d\n", BUILD_NUMBER)
	fmt.Fprintf(os.Stdout, "  Go-Easyops build timestamp  : %d\n", BUILD_TIMESTAMP)
	fmt.Fprintf(os.Stdout, "  Go-Easyops build time       : %s\n", time.Unix(BUILD_TIMESTAMP, 0))
	fmt.Fprintf(os.Stdout, "  Go-Easyops description      : %s\n", BUILD_DESCRIPTION)

	fmt.Fprintf(os.Stdout, "  App version                 : %d\n", appinfo.AppInfo().Number)
	fmt.Fprintf(os.Stdout, "  App build timestamp         : %d\n", appinfo.AppInfo().Timestamp)
	fmt.Fprintf(os.Stdout, "  App build time              : %s\n", time.Unix(appinfo.AppInfo().Timestamp, 0))
	fmt.Fprintf(os.Stdout, "  App description             : %s\n", appinfo.AppInfo().Description)
	fmt.Fprintf(os.Stdout, "  App artefactid              : %d\n", appinfo.AppInfo().ArtefactID)
	fmt.Fprintf(os.Stdout, "  App repository              : %d\n", appinfo.AppInfo().RepositoryID)
	fmt.Fprintf(os.Stdout, "  App repository git url      : %s\n", appinfo.AppInfo().GitURL)
	fmt.Fprintf(os.Stdout, "  Source code path            : %s\n", SourceCodePath())

	PrintDefaults()
}
func PrintDefaults() {
	if print_easyops {
		fmt.Fprintf(os.Stdout, "\nGo-easyops Usage:\n")
	} else {
		fmt.Fprintf(os.Stdout, "\nUsage:\n")
	}
	f := flag.CommandLine
	f.VisitAll(func(fg *flag.Flag) {
		isext := strings.HasPrefix(fg.Name, "ge_")
		if !isext {
			for _, s := range internal_flag_names {
				if fg.Name == s {
					isext = true
					break
				}
			}
		}
		if print_easyops != isext {
			return
		}
		s := fmt.Sprintf("  -%s", fg.Name) // Two spaces before -; see next two comments.
		name, usage := flag.UnquoteUsage(fg)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.ReplaceAll(usage, "\n", "\n    \t")

		s += fmt.Sprintf(" (default %v)", fg.DefValue)

		fmt.Printf("%s\n", s)
	})
	fmt.Printf(`
Yaml Mapping file: /etc/yacloud/config/service_map.yaml
Defaults override file: /tmp/goeasyops.config
Environment Variables:
`)

	s := render_env_help()
	fmt.Println(s)

}
func GetInstanceID() string {
	s := *instance_id
	if s == "" {
		mlock.Lock()
		defer mlock.Unlock()
		if *instance_id != "" {
			return *instance_id
		}
		s = "L-" + RandomString(32)
		*instance_id = s
	}
	return s
}
func GetPid() uint64 {
	p := os.Getpid()
	return uint64(p)
}

// get registry address as per -registry parameter, or if -registry_resolver is set, use that
func GetClientRegistryAddress() string {
	if manreg != "" {
		return manreg
	}
	if *registry_resolver == "" {
		return GetRegistryAddress()
	}
	res := *registry_resolver
	if !strings.Contains(res, ":") {
		res = fmt.Sprintf("%s:5000", res)
	}
	return res
}

// programmatically override -registry_resolver flag
func SetClientRegistryAddress(reg string) {
	if !strings.Contains(reg, ":") {
		reg = fmt.Sprintf("%s:5000", reg)
	}
	manreg = reg
	common.NotifyRegistryChangeListeners()
}

// get registry address as per -registry parameter
func GetRegistryAddress() string {
	res := *registry
	if *registry == REGISTRY_DEFAULT {
		s := reg_env.Value()
		if s != "" {
			res = s
		}
	}
	if *registry == REGISTRY_DEFAULT {
		if config != nil && config.Registry != "" {
			res = config.Registry
		}
	}
	if !strings.Contains(res, ":") {
		res = fmt.Sprintf("%s:5000", res)
	}
	return res
}

func Datacenter() bool {
	return *running_in_datacenter
}

// for testing purposes to mock parameter -AD_started_by_auto_deployer
func SetDatacenter(b bool) {
	*running_in_datacenter = b
}

// if (para != "") { return para }, else return os.GetEnv(envname)
func OptEnvString(para, envname string) string {
	if para != "" {
		return para
	}
	return os.Getenv(envname)
}
func doappinfo() {
	fmt.Printf("%d\n", appinfo.AppInfo().Number)
	os.Exit(0)
}
func IsStandalone() bool {
	return *stdalone
}
func LocalRegistrationDir() string {
	return "/tmp/local_registry"
}
func ContextWithBuilder() bool {
	return true

}

// this is for testing purposes to mock the parameter -ge_context_with_builder
func GetContextBuilderVersion() int {
	version := *context_build_version
	if version != 2 {
		panic(fmt.Sprintf("Unsupported context version (%d)", version))
	}
	return version
}

// this is for testing purposes to mock the parameter -ge_context_with_builder
func SetContextBuilderVersion(version int) {
	if version != 2 {
		panic(fmt.Sprintf("Unsupported context version (%d)", version))
	}
	*context_build_version = version
}

// get a serialised context from environment variable GE_CTX
func GetEnvContext() string {
	if overridden_env_context != "" {
		s := overridden_env_context
		if len(s) > 10 {
			s = s[:10]
		}
		fmt.Printf("[go-easyops] using overriden env context (%s )\n", s)
		return overridden_env_context
	}
	return e_ctx.Value()
}
func DebugSignature() bool {
	return *debug_sig
}
func DebugAuth() bool {
	return *debug_auth
}

// this is for testing purposes to mock the environment variable GE_CTX
func SetEnvContext(s string) {
	overridden_env_context = s
}

// usually returns /opt/yacloud/current
func GetYACloudDir() string {
	dirs := []string{
		"/opt/yacloud/current",
		"/opt/yacloud/",
	}
	for _, res := range dirs {
		dname := res + "/ctools"
		st, err := os.Stat(dname)
		if err != nil {
			continue
		}
		if st.IsDir() {
			return res
		}
	}
	return ""
}

// default timeout for new contexts
func DefaultTimeout() time.Duration {
	return *default_timeout
}

func EnabledExperiments() []string {
	exs := *enabled_experiments
	if len(exs) == 0 {
		return nil
	}
	ex := strings.Split(exs, ",")
	var res []string
	for _, e := range ex {
		e = strings.Trim(e, " ")
		res = append(res, e)
	}
	return res
}

// print context debug stuff
func DebugfContext(format string, args ...interface{}) {
	if !*debug_ctx {
		return
	}
	x := fmt.Sprintf(format, args...)
	fmt.Printf("[go-easyops/debugctx] %s", x)
}

// print context debug stuff
func DebugfRPC(format string, args ...interface{}) {
	if !*debug_ctx {
		return
	}
	x := fmt.Sprintf(format, args...)
	fmt.Printf("[go-easyops/rpc] %s", x)
}

func SetDebugContext() {
	*debug_ctx = true
}
func IsDebugRPCClient() bool {
	return *debug_rpc_client
}
func IsDebugRPCServer() bool {
	return *debug_rpc_serve
}

// is this a flag defined and used by go-easyops?
func IsEasyopsFlag(name string) bool {
	if strings.HasPrefix(name, "ge_") {
		return true
	}
	for _, ifn := range internal_flag_names {
		if ifn == name {
			return true
		}
	}
	return false
}

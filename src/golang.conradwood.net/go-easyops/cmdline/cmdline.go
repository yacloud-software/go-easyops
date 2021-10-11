package cmdline

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// annoyingly, not all go-easyops flags start with ge_
	internal_flag_names = []string{"token", "registry", "registry_resolver", "AD_started_by_auto_deployer"}

	mlock                 sync.Mutex
	running_in_datacenter = flag.Bool("AD_started_by_auto_deployer", false, "the autodeployer sets this to true to modify the behaviour to make it suitable for general-availability services in the datacenter")

	registry          = flag.String("registry", "localhost:5000", "address of the registry server. This is used for registration as well as resolving unless -registry_resolver is set, in which case this is only used for registration")
	registry_resolver = flag.String("registry_resolver", "", "address of the registry server (for lookups)")
	instance_id       = flag.String("ge_instance_id", "", "autodeployers internal instance id. We may use this to get information about ourselves")
	ext_help          = flag.Bool("X", false, "extended help")
	print_easyops     = false
	manreg            = ""
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
}

// if we have a -X argument we will print extended usage AFTER flags are parsed.
// we know flags are parsed if ext_help flag (-X) turns true (timeout after 5 secs)
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

	fmt.Fprintf(os.Stdout, "  App version                 : %d\n", APP_BUILD_NUMBER)
	fmt.Fprintf(os.Stdout, "  App build timestamp         : %d\n", APP_BUILD_TIMESTAMP)
	fmt.Fprintf(os.Stdout, "  App build time              : %s\n", time.Unix(APP_BUILD_TIMESTAMP, 0))
	fmt.Fprintf(os.Stdout, "  App description             : %s\n", APP_BUILD_DESCRIPTION)
	fmt.Fprintf(os.Stdout, "  App repository              : %s\n", APP_BUILD_REPOSITORY)
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

}
func GetInstanceID() string {
	s := *instance_id
	if s == "" {
		mlock.Lock()
		defer mlock.Unlock()
		if *instance_id != "" {
			return *instance_id
		}
		s = "L-" + utils.RandomString(32)
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
}

// get registry address as per -registry parameter
func GetRegistryAddress() string {
	res := *registry
	if !strings.Contains(res, ":") {
		res = fmt.Sprintf("%s:5000", res)
	}
	return res
}

func Datacenter() bool {
	return *running_in_datacenter
}

// return para if non empty or retrieve value from environment and return that instead
func OptEnvString(para, envname string) string {
	if para != "" {
		return para
	}
	return os.Getenv(envname)
}

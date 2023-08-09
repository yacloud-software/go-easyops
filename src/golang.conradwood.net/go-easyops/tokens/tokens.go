package tokens

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"
)

var (
	token                    = flag.String("token", "", "service token")
	disusertoken             = flag.Bool("ge_disable_user_token", false, "if true disable reading of user token (for testing)")
	Deadline                 = flag.Int("ge_ctx_deadline_seconds", 10, "do not change for production services")
	tokenwasread             = false
	usertoken                string
	last_token_read_registry string
	cloudname                = "yacloud.eu"
	debug                    = flag.Bool("ge_debug_tokens", false, "debug user token stuff")
)

const (
	METANAME = "goeasyops_meta" // marshaled proto
// METANAME2 = "goeasyopsv2_meta" // marshaled proto
)

func DisableUserToken() {
	tokenwasread = true
	usertoken = ""
}

// OUTBOUND metadata...
func buildMeta() metadata.MD {
	panic("obsolete codepath")
}

// this builds a *NEW* token (detached from previous contexts)
// if there is neither a -token parameter nor a user token it will
// look at Environment variable GE_CTX and deserialise it
// this function is deprecated, obsolete and broken. use authremote.Context() instead
func DISContextWithToken() context.Context {
	// we need to allow this as long as we have OLD contexts that need deserialising
	/*
		if cmdline.ContextWithBuilder() {
			utils.NotImpl("(context_with_builder) tokens.ContextWithToken - V1 context only\n")
		}
	*/
	md := buildMeta()
	ctx, cnc := context.WithTimeout(context.Background(), time.Duration(*Deadline)*time.Second)
	go func(cf context.CancelFunc) {
		time.Sleep(time.Duration((*Deadline)+5) * time.Second)
		cnc()
	}(cnc)
	return metadata.NewOutgoingContext(ctx, md)

}

// this function is deprecated, obsolete and broken. use authremote.Context() instead
func DISContextWithTokenAndTimeout(seconds uint64) context.Context {
	if cmdline.ContextWithBuilder() {
		utils.NotImpl("contextv2 incomplete")
	}
	md := buildMeta()
	ctx, cnc := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	go func(cf context.CancelFunc, s uint64) {
		time.Sleep(time.Duration(s+5) * time.Second)
		cnc()
	}(cnc, seconds)
	return metadata.NewOutgoingContext(ctx, md)
}

// this function is deprecated, obsolete and broken. use authremote.Context() instead
func DISContext2WithTokenAndTimeout(seconds uint64) (context.Context, context.CancelFunc) {
	if cmdline.ContextWithBuilder() {
		panic("contextv2 incomplete")
	}
	md := buildMeta()
	ctx, cnc := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cnc
}

func SetServiceTokenParameter(tok string) {
	*token = tok
}
func GetServiceTokenParameter() string {
	return cmdline.OptEnvString(*token, "GE_TOKEN")
}
func readToken(token string) string {
	var tok string
	var btok []byte
	var fname string
	fname = "n/a"
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("[go-easyops] Failed to determine current user: %s\n", err)
		return ""
	}
	fnames := []string{
		fmt.Sprintf("%s/.go-easyops/tokens/%s.%s", usr.HomeDir, token, cloudname),
		fmt.Sprintf("%s/.go-easyops/tokens/%s.%s", usr.HomeDir, token, cmdline.GetRegistryAddress()),
		fmt.Sprintf("%s/.go-easyops/tokens/%s", usr.HomeDir, token),
	}
	for _, fname = range fnames {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			if *debug {
				fmt.Printf("File \"%s\" does not exist - ignoring.\n", fname)
			}
			continue
		}
		break
	}

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		// services usually don't have such file. don't throw errors in that case. it's confusing
		return ""
	}
	if *debug {
		fmt.Printf("[go-easyops] Reading file \"%s\", parsing as token\n", fname)
	}
	btok, err = ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("[go-easyops] Failed to read user token: %s\n", err)
		return ""
	} else if len(btok) == 0 {
		fmt.Printf("[go-easyops] Failed to read user token from %s\n", fname)
		return ""

	} else {
		tok = string(btok)
	}
	tok = strings.TrimSpace(tok)
	//	fmt.Printf("[go-easyops] read token from \"%s\"\n", fname)
	if *debug {
		fmt.Printf("[go-easyops] token is empty\n")
		if len(tok) < 5 {
			fmt.Printf("[go-easyops] token is too short\n")
		} else {
			fmt.Printf("[go-easyops] token: %s...\n", tok[:3])
		}
	}
	return tok
}

/*
get a usertoken parameter from:

1. GE_USERTOKEN

2. ~/.go-easyops/user_token

if ge_disable_user_token is true, return "" (empty string)

if GE_TOKEN is set, does not read file (but honour GE_USERTOKEN)
*/
func GetUserTokenParameter() string {
	if *disusertoken {
		//fmt.Printf("[go-easyops] tokens: user token disabled\n")
		return ""
	}
	ut := cmdline.OptEnvString("", "GE_USERTOKEN")
	if ut != "" {
		return ut
	}
	// if token is set either as parameter or as ENV variable GE_TOKEN, then return ""
	// because we are a service (services do not run as users)
	if cmdline.OptEnvString(*token, "GE_TOKEN") != "" {
		return ""
	}
	if tokenwasread && last_token_read_registry == cmdline.GetClientRegistryAddress() {
		return usertoken
	}
	u := readToken("user_token")
	usertoken = u
	tokenwasread = true
	last_token_read_registry = cmdline.GetClientRegistryAddress()
	/*
		if usertoken == "" {
			fmt.Printf("[go-easyops] tokens: reading usertoken, but is empty\n")
		} else {
			fmt.Printf("[go-easyops] got usertoken\n")
		}
	*/
	return usertoken
}

func SaveUserToken(token string) error {
	usr, err := user.Current()
	if err != nil {
		err := fmt.Errorf("Failed to determine current operating system user: %s (are you logged in to your computer?)\n", err)
		return err
	}
	dir := fmt.Sprintf("%s/.go-easyops/tokens", usr.HomeDir)
	os.MkdirAll(dir, 0700)
	fname := fmt.Sprintf("%s/.go-easyops/tokens/user_token", usr.HomeDir)
	err = utils.WriteFile(fname, []byte(token))
	if err != nil {
		fmt.Printf("Failed to save access token: %s\n", err)
		return err
	}
	tokenwasread = false
	return nil
}
func GetCloudName() string {
	return cloudname
}
func SetCloudName(xcloudname string) {
	if xcloudname != "" && xcloudname != cloudname {
		if *debug {
			fmt.Printf("setting cloud name to \"%s\"\n", xcloudname)
		}
		cloudname = xcloudname
		tokenwasread = false // force re-read of token
	}
}

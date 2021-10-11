package tokens

import (
	"context"
	"flag"
	"fmt"
	rc "golang.conradwood.net/apis/rpcinterceptor"
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
	token        = flag.String("token", "", "service token")
	disusertoken = flag.Bool("ge_disable_user_token", false, "if true disable reading of user token (for testing)")
	Deadline     = flag.Int("ge_ctx_deadline_seconds", 10, "do not change for production services")
	tokenwasread = false
	usertoken    string
	cloudname    = "yacloud.eu"
)

const (
	METANAME = "goeasyops_meta" // marshaled proto
)

func DisableUserToken() {
	tokenwasread = true
	usertoken = ""
}

// OUTBOUND metadata...
func buildMeta() metadata.MD {
	user := GetUserTokenParameter()
	im := &rc.InMetadata{ServiceToken: cmdline.OptEnvString(*token, "GE_TOKEN"), UserToken: user, FooBar: "moo_none"}
	ims, err := utils.Marshal(im)
	if err != nil {
		fmt.Printf("[go-easyops] WARNING - failed to marshal metadata (%s)\n", err)
	}
	res := metadata.Pairs(
		METANAME, ims,
	)
	return res
}

// this builds a *NEW* token (detached from previous contexts)
// if there is neither a -token parameter nor a user token it will
// look at Environment variable GE_CTX and deserialise it
func ContextWithToken() context.Context {
	md := buildMeta()
	ctx, cnc := context.WithTimeout(context.Background(), time.Duration(*Deadline)*time.Second)
	go func(cf context.CancelFunc) {
		time.Sleep(time.Duration((*Deadline)+5) * time.Second)
		cnc()
	}(cnc)
	return metadata.NewOutgoingContext(ctx, md)

}
func ContextWithTokenAndTimeout(seconds uint64) context.Context {
	md := buildMeta()
	ctx, cnc := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	go func(cf context.CancelFunc, s uint64) {
		time.Sleep(time.Duration(s+5) * time.Second)
		cnc()
	}(cnc, seconds)
	return metadata.NewOutgoingContext(ctx, md)
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
			continue
		}
		break
	}

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		// services usually don't have such file. don't throw errors in that case. it's confusing
		return ""
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

	return tok
}

// get the user token as was passed in through the commandline
func GetUserTokenParameter() string {
	if *disusertoken {
		return ""
	}
	if cmdline.OptEnvString(*token, "GE_TOKEN") != "" {
		return ""
	}
	if tokenwasread {
		return usertoken
	}
	u := readToken("user_token")
	usertoken = u
	tokenwasread = true
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

func SetCloudName(xcloudname string) {
	if xcloudname != "" && xcloudname != cloudname {
		cloudname = xcloudname
		tokenwasread = false // force re-read of token
	}
}

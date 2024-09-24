package client

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/cmdline"
	cm "golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
)

/*
we attempt to get the public key for authentication from the auth-server
*/
var (
	last_registry           string
	pubkeylock              sync.Mutex
	retrieve_sig_lock       = utils.NewTimeoutLock("retrieve_sig_lock")
	retrieved_sig           = false
	retrieving              = false
	no_retrieve             = flag.Bool("ge_disable_dynamic_pubkey", false, "if true, disable the lookup of the public key from authservice on startup")
	crash_on_race_condition = flag.Bool("ge_development_signature", false, "crash and burn if getsignaturefromauth detects a race condition")
)

// cannot use init function, because flags might not be initialised (for example registry flag)
/*
func init() {
	GetSignatureFromAuth()
}
*/
func GotSig() bool {
	cur := cmdline.GetClientRegistryAddress()
	if cur != last_registry {
		return false
	}
	return retrieved_sig
}
func init() {
	cm.AddRegistryChangeReceiver(GetSignatureFromAuth)
}
func GetSignatureFromAuth() {
	if *crash_on_race_condition {
		if stack_includes("GetSignatureFromAuth", 2) {
			utils.PrintStack("reentry")
			panic("re-entry issue")
		}
	}
	if cmdline.IsStandalone() {
		//TODO: set a "fake" signature
		retrieved_sig = true
		retrieving = false
		last_registry = cmdline.GetClientRegistryAddress()
		return
	}
	if GotSig() {
		return
	}
	if *no_retrieve {
		return
	}
	if retrieving {
		b := retrieve_sig_lock.LockWithTimeout(time.Duration(10) * time.Second) // hack to wait for retrieving to finish
		if b {
			retrieve_sig_lock.Unlock()
			return
		}

		fmt.Printf("[go-easyops] WARNING Signature retrieving seems to be stuck\n")
		return
	}
	pubkeylock.Lock()
	if retrieving {
		pubkeylock.Unlock()
		return
	}

	retrieve_sig_lock.Lock()
	defer retrieve_sig_lock.Unlock()
	retrieving = true

	pubkeylock.Unlock()
	if GotSig() {
		return
	}
	started_sig := time.Now()
	if cmdline.DebugSignature() {
		fmt.Printf("[go-easyops] Retrieving signature and cloudname...\n")
	}
	cb := ctx.NewContextBuilder()
	cctx := cb.ContextWithAutoCancel()
	cctx = context.Background()
	rega := cmdline.GetClientRegistryAddress()
	cn := ConnectAtNoAuth(rega, "auth.AuthenticationService")
	authServer := apb.NewAuthenticationServiceClient(cn)
	pk, err := authServer.GetPublicSigningKey(cctx, &common.Void{})
	if err != nil {
		fmt.Printf("[go-easyops] failed to get public auth key (%s)\n", err)
		cn.Close()
		retrieving = false
		return
	}
	if cmdline.DebugSignature() {
		fmt.Printf("[go-easyops] CloudName=\"%s\" (after %0.1fs) \n", pk.CloudName, time.Since(started_sig).Seconds())
	}
	tokens.SetCloudName(pk.CloudName)
	cm.SetPublicSigningKey(pk)
	cn.Close()
	retrieved_sig = true
	retrieving = false
	last_registry = cmdline.GetClientRegistryAddress()
	if cmdline.DebugSignature() {
		fmt.Printf("[go-easyops] got Signature and cloudname (%s) from registry %s (after %0.1fs)\n", pk.CloudName, last_registry, time.Since(started_sig).Seconds())
	}
}

func stack_includes(pat string, min int) bool {
	s := utils.GetStack("")
	if strings.Count(s, pat) >= min {
		return true
	}
	return false
}

package client

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/cmdline"
	cm "golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/tokens"
	"sync"
)

/*
we attempt to get the public key for authentication from the auth-server
*/
var (
	last_registry = ""
	pubkeylock    sync.Mutex
	retrieved_sig = false
	retrieving    = false
	no_retrieve   = flag.Bool("ge_disable_dynamic_pubkey", false, "if true, disable the lookup of the public key from authservice on startup")
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
		return
	}
	pubkeylock.Lock()
	if retrieving {
		pubkeylock.Unlock()
		return
	}
	retrieving = true
	pubkeylock.Unlock()
	if GotSig() {
		return
	}
	if cmdline.DebugAuth() {
		fmt.Printf("[go-easyops] Retrieving signature and cloudname...\n")
	}
	cb := ctx.NewContextBuilder()
	cctx := cb.ContextWithAutoCancel()
	cctx = context.Background()
	cn := ConnectAt(cmdline.GetClientRegistryAddress(), "auth.AuthenticationService")
	authServer := apb.NewAuthenticationServiceClient(cn)
	pk, err := authServer.GetPublicSigningKey(cctx, &common.Void{})
	if err != nil {
		fmt.Printf("[go-easyops] failed to get public auth key (%s)\n", err)
		cn.Close()
		retrieving = false
		return
	}
	//	fmt.Printf("[go-easyops] CloudName=\"%s\"\n", pk.CloudName)
	tokens.SetCloudName(pk.CloudName)
	cm.SetPublicSigningKey(pk)
	cn.Close()
	retrieved_sig = true
	retrieving = false
	last_registry = cmdline.GetClientRegistryAddress()
	if cmdline.DebugAuth() {
		fmt.Printf("[go-easyops] got Signature and cloudname (%s) from registry %s\n", pk.CloudName, last_registry)
	}
}

package client

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	cm "golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/tokens"
	"sync"
)

/*
 we attempt to get the public key for authentication from the auth-server
*/
var (
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
	return retrieved_sig
}
func GetSignatureFromAuth() {
	if retrieved_sig {
		return
	}
	if *no_retrieve {
		return
	}
	pubkeylock.Lock()
	if retrieving {
		pubkeylock.Unlock()
		return
	}
	retrieving = true
	pubkeylock.Unlock()
	if retrieved_sig {
		return
	}
	ctx := context.Background()
	cn := Connect("auth.AuthenticationService")
	authServer := apb.NewAuthenticationServiceClient(cn)
	pk, err := authServer.GetPublicSigningKey(ctx, &common.Void{})
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
}

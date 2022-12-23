package authremote

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"sync"
	"time"
)

var (
	local_service_lock     sync.Mutex
	local_service          *apb.SignedUser
	local_service_resolved = false
	local_user_lock        sync.Mutex
	local_user             *apb.SignedUser
	local_user_resolved    = false
)

// the local user account (nil if a service)
func getLocalUserAccount() *apb.SignedUser {
	if local_user_resolved {
		return local_user
	}
	local_user_lock.Lock()
	defer local_user_lock.Unlock()
	if local_user_resolved {
		return local_user
	}
	st := tokens.GetUserTokenParameter()
	if st == "" {
		fmt.Printf("[go-easyops] no user account, assuming cli tool\n")
		local_user_resolved = true
		return nil
	}
	atr := &apb.AuthenticateTokenRequest{Token: st}
	timeout := time.Duration(50) * time.Second
	ctx, cnc := context.WithTimeout(context.Background(), timeout)
	go auto_cancel(cnc, timeout)
	fmt.Printf("[go-easyops] verifying and resolving local user account\n")
	ar, err := GetAuthenticationService().SignedGetByToken(ctx, atr)
	if err != nil {
		fmt.Printf("Unable to resolve user token.(%s)\n", utils.ErrorString(err))
		panic("unable to resolve user token")
	}
	if !ar.Valid {
		fmt.Printf("invalid token: %s\n(%s)\n", ar.PublicMessage, ar.LogMessage)
		panic("Invalid user token")
	}
	fmt.Printf("[go-easyops] local user: %s\n", auth.Description(common.VerifySignedUser(ar.User)))
	local_user = ar.User
	local_user_resolved = true
	return ar.User
}

// the local service's useraccount (nil if on commandline or service without useraccount)
func getLocalServiceAccount() *apb.SignedUser {
	if local_service_resolved {
		return local_service
	}
	local_service_lock.Lock()
	defer local_service_lock.Unlock()
	if local_service_resolved {
		return local_service
	}
	st := tokens.GetServiceTokenParameter()
	if st == "" {
		fmt.Printf("[go-easyops] no service account, assuming cli tool\n")
		local_service_resolved = true
		return nil
	}
	fmt.Printf("[go-easyops] verifying and resolving local service account\n")
	atr := &apb.AuthenticateTokenRequest{Token: st}
	timeout := time.Duration(15) * time.Second
	ctx, cnc := context.WithTimeout(context.Background(), timeout)
	go auto_cancel(cnc, timeout)
	ar, err := GetAuthenticationService().SignedGetByToken(ctx, atr)
	if err != nil {
		fmt.Printf("Unable to resolve service token.(%s)\n", utils.ErrorString(err))
		panic("unable to resolve service token")
	}
	if !ar.Valid {
		fmt.Printf("invalid token: %s\n(%s)\n", ar.PublicMessage, ar.LogMessage)
		panic("Invalid service token")
	}
	local_service = ar.User
	local_service_resolved = true
	return ar.User
}

package auth

import (
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	//	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	//	"golang.conradwood.net/go-easyops/tokens"
	"context"
)

func GetUser(ctx context.Context) *apb.User {
	cs := rpc.CallStateFromContext(ctx)
	if cs == nil {
		return nil
	}
	return cs.User()
}

func PrintUser(u *apb.User) {
	if u == nil {
		return
	}
	fmt.Printf("User ID: %s\n", u.ID)
	fmt.Printf("  Email: %s\n", u.Email)
	fmt.Printf("  Abbrev:%s\n", u.Abbrev)
}
func PrintSignedUser(uu *apb.SignedUser) {
	u := common.VerifySignedUser(uu)
	if u == nil {
		return
	}

	fmt.Printf("User ID: %s\n", u.ID)
	fmt.Printf("  Email: %s\n", u.Email)
	fmt.Printf("  Abbrev:%s\n", u.Abbrev)
}

func GetService(ctx context.Context) *apb.User {
	cs := rpc.CallStateFromContext(ctx)
	if cs == nil {
		return nil
	}
	return cs.CallerService()
}

// one line description of the user/caller
func Description(user *apb.User) string {
	if user == nil {
		return "ANONYMOUS"
	}
	if user.Abbrev != "" {
		return user.Abbrev
	}
	if user.Email != "" {
		return user.Email
	}
	return "user #" + user.ID
}
func UserIDString(user *apb.User) string {
	if user == nil {
		return "ANONYMOUS"
	}
	if user.Abbrev != "" {
		return "#" + user.ID + " (" + user.Abbrev + ")"
	}
	if user.Email != "" {
		return "#" + user.ID + " (" + user.Email + ")"
	}
	return "user #" + user.ID
}

// returns  "User ID (email)"
func CurrentUserString(ctx context.Context) string {
	u := GetUser(ctx)
	if u == nil {
		return "ANONYMOUS"
	}
	return fmt.Sprintf("User #%s (%s)", u.ID, u.Email)
}

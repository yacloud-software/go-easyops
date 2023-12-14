/*
This provides information about users and their groups and organisations.

It provides many helper functions to determine if a user has certain permissions, is a root user or a service.
it has helpers to print out user information.

It works on a context to retrieve user and service information.

Note that a 'service account' is also a 'user'. It is just marked as a service. Service accounts are treated differently in some places. For example, a web proxy does not (should not) allow services to log in to a website.

As a rule of thumb, Most calls in this package do not require network I/O and thus are fairly safe and cheap. see package authremote for calls that require network I/O.

It also provides some wrappers to create a new context. That is for historic reasons. Developers should use and port code to use the ctx package instead. In almost all cases, if a context is already present (for example, in a gRPC server), that context must be re-used. New Contexts should only be created after a user has been authenticated.
*/
package auth

import (
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	//	"golang.conradwood.net/go-easyops/client"
	"context"
	//	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/ctx"
	"golang.conradwood.net/go-easyops/rpc"
)

// get the OrganisationID in this context
func GetOrganisationID(uctx context.Context) string {
	return "1" // currently ALWAYS 1
}

// get the user in this context
func GetUser(uctx context.Context) *apb.User {
	u := ctx.GetLocalState(uctx).User()
	us := common.VerifySignedUser(u)
	if us != nil {
		// new path succeeded
		return us
	}
	// code below to be removed:
	cs := rpc.CallStateFromContext(uctx)
	if cs == nil {
		return nil
	}
	return cs.User()
}

// get the user in this context
func GetSignedUser(uctx context.Context) *apb.SignedUser {
	u := ctx.GetLocalState(uctx).User()
	if u != nil {
		// new path succeeded
		return u
	}
	// code below to be removed:
	cs := rpc.CallStateFromContext(uctx)
	if cs == nil {
		return nil
	}
	su := cs.SignedUser()
	if su == nil && cs.User() != nil {
		panic("mismatched old style context, no signed user, but unsigned user present")
	}
	return su
}

// get the user in this context
func GetSignedService(uctx context.Context) *apb.SignedUser {
	u := ctx.GetLocalState(uctx).CallingService()
	if u != nil {
		return u
	}
	// code below to be removed, obsolete path...:
	cs := rpc.CallStateFromContext(uctx)
	if cs == nil {
		return nil
	}
	res := cs.SignedService()
	if res == nil && cs.CallerService() != nil {
		panic("invalid callstate (no signed service, but unsignedservice)")
	}
	return res
}

// get the service which directly called us
func GetService(uctx context.Context) *apb.User {
	return common.VerifySignedUser(GetSignedService(uctx))
}

// get the service which created this context
func GetCreatingService(uctx context.Context) *apb.User {
	u := ctx.GetLocalState(uctx).CreatorService()
	us := common.VerifySignedUser(u)
	return us

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

// one line description of the user/caller
func SignedDescription(user *apb.SignedUser) string {
	u := common.VerifySignedUser(user)
	return Description(u)
}
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

// print the userid and description
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

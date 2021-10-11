package auth

import (
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	//	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/common"
	"golang.conradwood.net/go-easyops/rpc"
	//	"golang.conradwood.net/go-easyops/tokens"
	"context"
	"strings"
)

func IsRoot(ctx context.Context) bool {
	return IsRootUser(rpc.CallStateFromContext(ctx).User())
}
func IsRootUser(user *apb.User) bool {
	return IsInGroupByUser(user, "1")
}

// return true if user is in any of the groups (comma delimited list of ids)
func IsInGroupsByUser(user *apb.User, groupids string) bool {
	for _, g := range strings.Split(groupids, ",") {
		g = strings.Trim(g, " ")
		if IsInGroupByUser(user, g) {
			return true
		}
	}
	return false
}

/*
* return true if user is in this group
 */
func IsInGroupByUser(user *apb.User, groupid string) bool {
	if user == nil || groupid == "" || user.Groups == nil {
		return false
	}
	for _, g := range user.Groups {
		if g.ID == groupid {
			return true
		}
	}

	return false
}

// return true if user (from context) is part of group specified by groupid
func IsInGroup(ctx context.Context, groupid string) bool {
	u := GetUser(ctx)
	return IsInGroupByUser(u, groupid)
}

// return true if user (from context) is part of at least one of the groups specified by groupids. groupids is a comma delimited list of groupids
func IsInGroups(ctx context.Context, groupids string) bool {
	u := GetUser(ctx)
	return IsInGroupsByUser(u, groupids)
}
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

func CurrentUserString(ctx context.Context) string {
	u := GetUser(ctx)
	if u == nil {
		return "ANONYMOUS"
	}
	return fmt.Sprintf("User #%s (%s)", u.ID, u.Email)
}

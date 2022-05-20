package auth

import (
	"context"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/rpc"
	"strings"
)

func IsRoot(ctx context.Context) bool {
	return IsRootUser(rpc.CallStateFromContext(ctx).User())
}
func IsRootUser(user *apb.User) bool {
	return IsInGroupByUser(user, "1")
}

// return true if service in context is one of the serviceids. serviceids comma delimited
func IsService(ctx context.Context, serviceids string) bool {
	svc := GetService(ctx)
	if svc == nil {
		return false
	}
	for _, g := range strings.Split(serviceids, ",") {
		g = strings.Trim(g, " ")
		if svc.ID == g {
			return true
		}
	}
	return false
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

package authremote

import (
	"context"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/common"
	"sync"
)

var (
	cached_user_lock sync.Mutex
	cached_users     []*user_cache_entry
)

type user_cache_entry struct {
	userid      string
	signed_user *apb.SignedUser
	user        *apb.User
}

func (uc *user_cache_entry) isValid() bool {
	return true
}

func usercache_GetSignedUserByID(ctx context.Context, userid string) (*apb.SignedUser, error) {
	for _, uc := range cached_users {
		if uc.isValid() && uc.userid == userid {
			return uc.signed_user, nil
		}
	}

	managerClient()
	res, err := authManager.SignedGetUserByID(ctx, &apb.ByIDRequest{UserID: userid})
	if err != nil {
		return nil, err
	}
	usercache_add(res)
	return res, nil
}
func usercache_GetUserByID(ctx context.Context, userid string) (*apb.User, error) {
	for _, uc := range cached_users {
		if uc.isValid() && uc.userid == userid {
			return uc.user, nil
		}
	}

	managerClient()
	res, err := authManager.SignedGetUserByID(ctx, &apb.ByIDRequest{UserID: userid})
	if err != nil {
		return nil, err
	}
	u := usercache_add(res)
	return u, nil
}

func usercache_add(u *apb.SignedUser) *apb.User {
	user := common.VerifySignedUser(u)
	if user == nil {
		return nil
	}
	userid := user.ID
	cached_user_lock.Lock()
	defer cached_user_lock.Unlock()
	for _, uc := range cached_users {
		if uc.isValid() && uc.userid == userid {
			uc.user = user
			uc.signed_user = u
			return uc.user // nothing to do
		}
	}
	uc := &user_cache_entry{userid: userid, signed_user: u, user: user}
	cached_users = append(cached_users, uc)
	return uc.user
}

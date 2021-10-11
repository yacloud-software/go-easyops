package server

import (
	"sync"
	"time"
)

var (
	excessiveLock sync.Mutex
	usercache     = make(map[string]*UserCache)
)

func addUserToCache(token string, id string) {
	uc := UserCache{UserID: id, created: time.Now()}
	excessiveLock.Lock()
	usercache[token] = &uc
	excessiveLock.Unlock()
}

func getUserFromCacheByToken(token string) *UserCache {
	excessiveLock.Lock()
	res := usercache[token]
	excessiveLock.Unlock()
	return res
}

// return userid
func getUserFromCache(token string) string {
	uc := getUserFromCacheByToken(token)
	if uc == nil {
		return ""
	}
	if time.Since(uc.created) > (time.Minute * 5) {
		return ""
	}
	return uc.UserID

}

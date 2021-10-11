package auth

import (
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/go-easyops/cache"
	"sync"
	"time"
)

var (
	users     = cache.New("goeasyops_usercache", time.Duration(60)*time.Second, 100)
	cacheLock sync.Mutex
)

func cacheAdd(a *apb.User) {
	if a == nil {
		return
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()
	users.Put(a.ID, a)
}

func cacheByID(id string) *apb.User {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	l := users.Get(id)
	if l == nil {
		return nil
	}
	return l.(*apb.User)
}

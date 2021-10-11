package common

// this packge MUST NOT add command line variables.
// it is used by server & client to share common state
// no side effects please!

import (
	"sync"
)

var (
	myservicenames []string
	clients        []*Service
	blocked        []*BlockedService
	addlock        sync.Mutex
)

type Service struct {
	Name string
}

type BlockedService struct {
	Name    string
	Counter int
}

func AddServiceName(s string) {
	addlock.Lock()
	defer addlock.Unlock()
	for _, c := range clients {
		if c.Name == s {
			return
		}
	}
	clients = append(clients, &Service{Name: s})
}
func GetConnectionNames() []*Service {
	return clients
}
func GetBlockedConnectionNames() []*BlockedService {
	addlock.Lock()
	defer addlock.Unlock()
	var res []*BlockedService
	for _, x := range blocked {
		if x.Counter > 0 {
			res = append(res, x)
		}
	}
	return res
}
func AddBlockedServiceName(s string) {
	addlock.Lock()
	defer addlock.Unlock()
	for _, x := range blocked {
		if x.Name == s {
			x.Counter++
			return
		}
	}
	blocked = append(blocked, &BlockedService{Name: s, Counter: 1})
}
func RemoveBlockedServiceName(s string) {
	addlock.Lock()
	defer addlock.Unlock()
	for _, x := range blocked {
		if x.Name == s {
			x.Counter = 0
			return
		}
	}
}

func AddExportedServiceName(name string) {
	addlock.Lock()
	defer addlock.Unlock()
	for _, m := range myservicenames {
		if m == name {
			return
		}
	}
	myservicenames = append(myservicenames, name)
}

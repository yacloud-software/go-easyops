/*
commonly used by other go-easyops packages. May provide useful information on the state of go-easyops.
*/
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

// add a service name to the list of services being used
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

// get all services used
func GetConnectionNames() []*Service {
	return clients
}

// get services that are used, but are currently blocked because no targets are available
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

// mark a service as blocked by name
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

// unmark a service as blocked by name
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

// mark a service as exported.
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

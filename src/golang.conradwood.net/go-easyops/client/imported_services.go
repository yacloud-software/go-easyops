package client

import (
	"sync"
)

var (
	dependencies []string
	dep_lock     sync.Mutex
)

func RegisterDependency(name string) {
	dep_lock.Lock()
	defer dep_lock.Unlock()
	for _, d := range dependencies {
		if d == name {
			return
		}
	}
	dependencies = append(dependencies, name)
}
func GetDependencies() []string {
	return dependencies
}

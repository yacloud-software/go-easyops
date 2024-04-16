package server

import (
	"sync"
)

var (
	active_rpc_lock sync.Mutex
	active_rpcs     = 0
)

func startRPC() {
	active_rpc_lock.Lock()
	active_rpcs++
	defer active_rpc_lock.Unlock()
}
func stopRPC() {
	active_rpc_lock.Lock()
	if active_rpcs == 0 {
		panic("[go-easyops] active_rpcs must never be negative")
	}
	active_rpcs--
	defer active_rpc_lock.Unlock()
}
func ActiveRPCs() int {
	return active_rpcs
}

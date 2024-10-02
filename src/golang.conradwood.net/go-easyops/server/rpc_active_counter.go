package server

import (
	"sync"

	"golang.conradwood.net/go-easyops/cmdline"
)

var (
	active_rpc_lock sync.Mutex
	active_rpcs     = 0
)

func startRPC() {
	active_rpc_lock.Lock()
	active_rpcs++
	active_rpc_lock.Unlock()
	cmdline.DebugfRPC("------------------ RPC ENTERED (%d)  --------------------\n", active_rpcs)
}
func stopRPC() {
	active_rpc_lock.Lock()
	if active_rpcs == 0 {
		panic("[go-easyops] active_rpcs must never be negative")
	}
	active_rpcs--
	active_rpc_lock.Unlock()
	cmdline.DebugfRPC("------------------ RPC COMPLETE (%d) --------------------\n", active_rpcs)
}
func ActiveRPCs() int {
	return active_rpcs
}

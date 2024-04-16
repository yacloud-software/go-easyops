package server

import (
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"sync"
	"time"
)

var (
	stoplock           sync.Mutex
	was_stop_requested = false
)

func IsStopping() bool {
	return was_stop_requested
}

func stop_requested() {
	stoplock.Lock()
	if was_stop_requested {
		stoplock.Unlock()
		return
	}
	was_stop_requested = true
	stoplock.Unlock()
	fmt.Printf("[go-easyops] server received stop request and is shutting down\n")
	SetHealth(common.Health_STOPPING)
	go func() {
		for {
			// no rpcs active
			x := make(chan bool, 10)
			stopping(x)
			for ActiveRPCs() != 0 {
				time.Sleep(time.Duration(1) * time.Second)
				response := &goeasyops.StopUpdate{Stopping: true, ActiveRPCs: uint32(ActiveRPCs())}
				b, err := utils.MarshalBytes(response)
				if err != nil {
					fmt.Printf("[go-easyops] unable to marshal stop response: %s\n", err)
				} else {
					_, err = unixipc_srv.Send("STOPREQUEST", b)
					if err != nil {
						fmt.Printf("[go-easyops] failed to send update to stoprequest: %s\n", err)
					}
				}

			}
			os.Exit(0)
		}
	}()
}

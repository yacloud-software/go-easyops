package server

import (
	"fmt"
	"golang.conradwood.net/apis/common"
)

var (
	health = common.Health_READY
)

/*
Set the health of this service. This is useful for services which have periods where they are unavailable. Typically this is directly after starting, but sometimes also a period where they gather data.

For a service with a delayed startup, the pattern is as follows:

		func main() {
		 server.SetHealth(common.Health_STARTING)
	         go do_initialisation()
		}
	        func do_initialisation() {
	            doing_slow_things() // ...
	            server.SetHealth(common.Health_READY)
	        }
*/
func SetHealth(h common.Health) error {
	rereg := false
	if h != health {
		rereg = true
	}
	health = h
	if rereg && startup_complete {
		reRegister()
	}
	if len(knownServices) == 0 {
		fmt.Printf("[go-easyops] server.SetHealth() must be called after server.ServerStartup().")
		fmt.Printf("[go-easyops] recommended usage:\n")
		fmt.Printf("[go-easyops] sd.SetOnStartupCallback(func() { server.SetHealth(common.Health_STARTING) })\n")
		panic("[go-easyops] attempt to send ipc health before server startup")
	}
	for _, svc := range knownServices {
		err := ipc_send_health(svc, h)
		if err != nil {
			return err
		}
	}
	return nil

}

func getHealthString() string {
	s := common.Health_name[int32(health)]
	return s
}

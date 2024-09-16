package server

import (
	//	"fmt"
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
		err := ipc_send_health(nil, h)
		if err != nil {
			return err
		}
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
func GetHealth() common.Health {
	return health
}

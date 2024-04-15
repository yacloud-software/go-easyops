package server

import (
	"golang.conradwood.net/apis/common"
)

var (
	health = common.Health_READY
)

func SetHealth(h common.Health) error {
	rereg := false
	if h != health {
		rereg = true
	}
	health = h
	if rereg && startup_complete {
		reRegister()
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

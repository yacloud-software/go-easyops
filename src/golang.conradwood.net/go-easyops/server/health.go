package server

import (
	"golang.conradwood.net/apis/common"
)

var (
	health = common.Health_READY
)

func SetHealth(h common.Health) {
	rereg := false
	if h != health {
		rereg = true
	}
	health = h
	if rereg && startup_complete {
		reRegister()
	}

}

func getHealthString() string {
	s := common.Health_name[int32(health)]
	return s
}

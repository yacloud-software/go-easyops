package server

const (
	STARTING = 1
	READY    = 2
	STOPPING = 3
)

type HEALTH int

var (
	health = HEALTH(READY)
)

func SetHealth(h HEALTH) {
	health = h
}

func getHealthString() string {
	if health == STARTING {
		return "STARTING"
	} else if health == READY {
		return "READY"
	} else if health == STOPPING {
		return "STOPPING"
	}
	return ""
}

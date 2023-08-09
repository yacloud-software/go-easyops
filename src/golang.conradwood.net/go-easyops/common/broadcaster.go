package common

import (
	"sync"
)

var (
	receiver_lock sync.Mutex
	receivers     = make(map[string]*broadcast_receiver)
)

type broadcast_receiver struct {
	callbacks []func()
}

func AddRegistryChangeReceiver(f func()) {
	receiver_lock.Lock()
	defer receiver_lock.Unlock()
	br := receivers["REGISTRYCHANGE"]
	if br == nil {
		br = &broadcast_receiver{}
		receivers["REGISTRYCHANGE"] = br
	}
	br.callbacks = append(br.callbacks, f)
}

func NotifyRegistryChangeListeners() {
	wake("REGISTRYCHANGE")
}
func wake(name string) {
	receiver_lock.Lock()
	br := receivers[name]
	receiver_lock.Unlock()
	if br == nil {
		return
	}
	for _, cb := range br.callbacks {
		cb()
	}
}

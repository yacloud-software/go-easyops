package utils

import "sync"

type LockedBool struct {
	sync.Mutex
	val bool
}

func (lb *LockedBool) Set(b bool) {
	lb.Lock()
	lb.val = b
	lb.Unlock()
}
func (lb *LockedBool) Value() bool {
	lb.Lock()
	res := lb.val
	lb.Unlock()
	return res
}

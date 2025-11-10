package utils

import "sync"

type LockedInt struct {
	sync.Mutex
	val int
}

// returns PREVIOUS value
func (li *LockedInt) Set(val int) int {
	li.Lock()
	res := li.val
	li.val = val
	li.Unlock()
	return res
}

// returns NEW value
func (li *LockedInt) Inc() int {
	li.Lock()
	li.val = li.val + 1
	res := li.val
	li.Unlock()
	return res
}

// returns NEW value
func (li *LockedInt) Dec() int {
	li.Lock()
	li.val = li.val - 1
	res := li.val
	li.Unlock()
	return res
}
func (li *LockedInt) Value() int {
	li.Lock()
	res := li.val
	li.Unlock()
	return res
}

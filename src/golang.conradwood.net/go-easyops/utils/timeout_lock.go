package utils

import (
	"fmt"
	"time"
)

type TimeoutLock interface {
	Unlock()
	Lock()
	LockWithTimeout(time.Duration) bool // true if lock was acquired
}

type timeoutlock struct {
	ch       chan bool
	name     string
	lockedat string
}

// create a new lock which can be used like sync.Mutex but also adds LockWithTimeout()
func NewTimeoutLock(name string) TimeoutLock {
	tl := &timeoutlock{
		name: name,
		ch:   make(chan bool, 1),
	}
	return tl
}

// lock - waits indefinitely for a lock to become available
func (tl *timeoutlock) Lock() {
	tl.ch <- true
	tl.lockedat = GetStack("lock %s", tl.name)
}

// lock, return true if it was able to acquire lock within duration t. false if it was not
func (tl *timeoutlock) LockWithTimeout(t time.Duration) bool {
	select {
	case tl.ch <- true:
		tl.lockedat = GetStack("lock %s", tl.name)
		return true
	case <-time.After(t):
		PrintStack("[go-easyops] lock \"%s\" timeout", tl.name)
		fmt.Printf("lock \"%s\" was locked at: %s\n", tl.name, tl.lockedat)
		return false
	}

}

// unlock this lock
func (tl *timeoutlock) Unlock() {
	select {
	case <-tl.ch:
		//
	default:
		// nothing to unlock
		fmt.Printf("[go-easyops] timeoutlock \"%s\" unlocked but was not locked previously\n", tl.name)
		PrintStack("[go-easyops]")
	}
}

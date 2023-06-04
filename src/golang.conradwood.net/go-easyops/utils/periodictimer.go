package utils

import (
	"sort"
	"sync"
	"time"
)

var ()

type PeriodicTimer struct {
	secs                  []uint32 // sorted as Largest first
	started               time.Time
	callback              func(*PeriodicTimer, uint32) error
	thread_running        bool // true if thread is running
	stop_request          bool // true if thread is meant to exit
	wait_chan             chan bool
	lastSuccessfulRunSecs uint32
	wasRunAtStart         bool // true if it was run at least once
	lock                  sync.Mutex
	runLock               sync.Mutex
}

/*
A "PeriodicTimer" executes a callback at certain intervals over a certain period of time.
For example, a periodictimer, defined as NewPeriodicTimer([]uint32{20,15,5}) will run for 20 seconds and call the callback
after 5 seconds, 10 seconds and 20 seconds.
The callback will be retried every second if it returns an error until it returns no error
callback will also be called each time "Start()" is called.
*/
func NewPeriodicTimer(secs []uint32, cb func(pt *PeriodicTimer, secsLapsed uint32) error) *PeriodicTimer {
	if len(secs) == 0 {
		secs = []uint32{0}
	}
	sort.Slice(secs, func(i, j int) bool {
		return secs[i] > secs[j]
	})
	pt := &PeriodicTimer{callback: cb, wait_chan: make(chan bool), secs: secs}
	return pt
}
func (pt *PeriodicTimer) Start() {
	pt.lock.Lock()
	pt.started = time.Now()
	pt.lastSuccessfulRunSecs = 0
	pt.stop_request = false
	pt.wasRunAtStart = false
	if !pt.thread_running {
		go pt.timerLoop()
		pt.thread_running = true
	}
	err := pt.run_callback(0)
	if err == nil {
		pt.wasRunAtStart = true
	}
	pt.lock.Unlock()

}
func (pt *PeriodicTimer) Stop() {
	pt.lock.Lock()
	pt.stop_request = true
	pt.lock.Unlock()
}

// wait for timer to either stop or expire
func (pt *PeriodicTimer) Wait() {
	<-pt.wait_chan
	pt.stop_request = true
}

func (pt *PeriodicTimer) timerLoop() {
	for {
		if pt.stop_request {
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
		pt.checkTime()

		if pt.is_running_for_as_long_as_need_be() {
			break
		}

	}
	pt.lock.Lock()
	pt.wait_chan <- true
	pt.stop_request = false
	if !pt.thread_running {
		pt.thread_running = false
	}
	pt.lock.Unlock()
}
func (pt *PeriodicTimer) is_running_for_as_long_as_need_be() bool {
	rs := time.Since(pt.started)
	if rs.Seconds() >= float64(pt.secs[0]) {
		return true
	}
	return false
}

// periodically called by timer
func (pt *PeriodicTimer) checkTime() {
	secs := time.Since(pt.started).Seconds()
	// work out which period we're in
	period := uint32(0)
	for _, r := range pt.secs {
		if float64(r) > secs {
			continue
		}
		period = r
		break
	}
	if period == pt.lastSuccessfulRunSecs && pt.wasRunAtStart {
		return
	}

	err := pt.run_callback(period)
	if err == nil {
		pt.lastSuccessfulRunSecs = period
		pt.wasRunAtStart = true
	}
}
func (pt *PeriodicTimer) run_callback(period uint32) error {
	pt.runLock.Lock()
	defer pt.runLock.Unlock()
	err := pt.callback(pt, uint32(period))
	return err
}
func (pt *PeriodicTimer) Secs() []uint32 {
	return pt.secs
}
func (pt *PeriodicTimer) LastStarted() time.Time {
	return pt.started
}

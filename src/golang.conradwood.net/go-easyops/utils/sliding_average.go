package utils

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

const (
	COUNTERS = 10
)

var (
	sliding_avg_debug = flag.Bool("ge_debug_sliding_average", false, "debug sliding avg code")
)

/*
it is often useful to take the most recent period of time average, for example last minute average.
however, it is useful to double-buffer this so that there is always a full sample available.
that is where this struct helps.
*/
type SlidingAverage struct {
	lock                sync.Mutex
	calc1               *sacalc
	calc2               *sacalc
	MinAge              time.Duration // minimum age before a counter is valid
	MinSamples          uint64        // minimum number of samples before a counter is valid
	updating_number_one bool          // if true updates calc1, otherwise calc2
}

func NewSlidingAverage() *SlidingAverage {
	res := &SlidingAverage{
		MinAge:     time.Duration(10) * time.Second,
		MinSamples: 10,
	}
	return res
}

// a tracker of an average
type sacalc struct {
	debug_name string
	is_fresh   bool // reset each time it is modified and set when it is cleared or new
	started    time.Time
	counts     []uint64
	counter    []uint64
}

func new_sacalc(debug_name string) *sacalc {
	res := &sacalc{
		debug_name: debug_name,
		is_fresh:   true,
		counts:     make([]uint64, COUNTERS),
		counter:    make([]uint64, COUNTERS),
	}
	return res
}

/**************************** double-buffering stuff ******************************/
// if the to_be_updated one meets criteria, then switch it to be the read one and make the other buffer updateable
func (sa *SlidingAverage) check_for_switch() {
	up := sa.to_be_updated()
	rd := sa.to_be_read()
	if sa.meetsCriteria(up) {
		sa.updating_number_one = !sa.updating_number_one
		if rd != nil {
			rd.make_fresh()
		}
	}
}
func (sa *SlidingAverage) meetsCriteria(sc *sacalc) bool {
	if time.Since(sc.started) < sa.MinAge {
		return false
	}
	samples := uint64(0)
	for i := 0; i < COUNTERS; i++ {
		samples = samples + sc.counts[i]
	}
	if samples < sa.MinSamples {
		return false
	}

	return true
}
func (sc *sacalc) make_fresh() {
	for i := 0; i < COUNTERS; i++ {
		sc.counts[i] = 0
		sc.counter[i] = 0
	}
	sc.is_fresh = true
}

/**************************** update counter stuff ******************************/
func (sa *SlidingAverage) to_be_updated() *sacalc {
	if sa.updating_number_one {
		if sa.calc1 == nil {
			sa.calc1 = new_sacalc("sacalc-1")
		}
		return sa.calc1
	}
	if sa.calc2 == nil {
		sa.calc2 = new_sacalc("sacalc-2")
	}
	return sa.calc2
}
func (sa *SlidingAverage) to_be_read() *sacalc {
	if sa.updating_number_one {
		return sa.calc2
	}
	return sa.calc1
}

// get number of counts
func (sa *SlidingAverage) GetCounts(counter int) uint64 {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sc := sa.to_be_read()
	if sc == nil {
		return 0
	}
	return sc.getCounts(counter)
}

// get a counter
func (sa *SlidingAverage) GetCounter(counter int) uint64 {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sc := sa.to_be_read()
	if sc == nil {
		return 0
	}
	return sc.getCounter(counter)
}

func (sa *SlidingAverage) Add(counter int, a uint64) {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sa.to_be_updated().Add(counter, a)
	sa.check_for_switch()

}

func (sc *sacalc) Add(counter int, a uint64) {
	sc.counts[counter]++
	sc.counter[counter] = sc.counter[counter] + a
	if sc.is_fresh {
		sc.started = time.Now()
		sc.is_fresh = false
	}
	sc.Printf("added: %d to counter #%d, now %d\n", a, counter, sc.counter[counter])
}
func (sc *sacalc) getCounter(counter int) uint64 {
	return sc.counter[counter]
}
func (sc *sacalc) getCounts(counter int) uint64 {
	return sc.counts[counter]
}

func (sc *sacalc) Printf(format string, args ...interface{}) {
	if !*sliding_avg_debug {
		return
	}
	s := "[" + sc.debug_name + "] " + format

	fmt.Printf(s, args...)

}

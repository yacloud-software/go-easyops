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
That is, it is guaranteed that the average is always calculated over at least MinAge. Periodically old average numbers are "dropped", so the average is also reflective of fresh values.

The InitialAge value may be set to start providing averages faster than MinAge upon startup.
*/
type SlidingAverage struct {
	lock                sync.Mutex
	calc1               *sacalc
	calc2               *sacalc
	InitialAge          time.Duration // time to wait before providing the first averages (after startup)
	MinAge              time.Duration // minimum age before a counter is valid
	MinSamples          uint64        // minimum number of samples before a counter is valid
	updating_number_one bool          // if true updates calc1, otherwise calc2
	created             time.Time
	created_via_new     bool // helper to detect instantiation via New()
	switched            bool // true if switched at least once
}

func NewSlidingAverage() *SlidingAverage {
	res := &SlidingAverage{
		InitialAge:      time.Duration(24) * time.Hour,
		created_via_new: true,
		created:         time.Now(),
		MinAge:          time.Duration(10) * time.Second,
		MinSamples:      10,
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
		sa.switched = true
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
	var alt_res *sacalc
	// got at least one full buf
	if sa.updating_number_one {
		alt_res = sa.calc2
	} else {
		alt_res = sa.calc1
	}
	if sa.switched || time.Since(sa.created) < sa.InitialAge {
		return alt_res
	}
	if sa.calc1 == nil && sa.calc2 != nil {
		return sa.calc2
	}
	if sa.calc2 == nil && sa.calc1 != nil {
		return sa.calc1
	}
	if sa.calc1 == nil && sa.calc2 == nil {
		return nil
	}

	if sa.calc1.is_fresh && !sa.calc2.is_fresh {
		alt_res = sa.calc2
	}

	if sa.calc2.is_fresh && !sa.calc1.is_fresh {
		alt_res = sa.calc1
	}

	return alt_res

}

// get number of counts
func (sa *SlidingAverage) GetCounts(counter int) uint64 {
	if !sa.created_via_new {
		panic("[go-easyops] SlidingAverage must be created with function NewSlidingAverage()")
	}
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
	if !sa.created_via_new {
		panic("[go-easyops] SlidingAverage must be created with function NewSlidingAverage()")
	}
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sc := sa.to_be_read()
	if sc == nil {
		return 0
	}
	return sc.getCounter(counter)
}
func (sa *SlidingAverage) GetAverage(counter int) float64 {
	num := sa.GetCounter(counter)
	counts := sa.GetCounts(counter)
	if counts == 0 || num == 0 {
		return 0
	}
	res := float64(num) / float64(counts)
	sa.printf("result: num=%0.1f, counts=%0.1f -> %0.1f\n", num, counts, res)
	return res
}

// per second added rate
func (sa *SlidingAverage) GetRate(counter int) float64 {
	sa.lock.Lock()
	defer sa.lock.Unlock()
	sc := sa.to_be_read()
	if sc == nil {
		return 0.0
	}
	num := float64(sc.getCounter(counter))
	secs := time.Since(sc.started).Seconds()
	if num == 0 || secs == 0 {
		return 0.0
	}
	res := num / secs
	return res
}

func (sa *SlidingAverage) Add(counter int, a uint64) {
	if !sa.created_via_new {
		panic("[go-easyops] SlidingAverage must be created with function NewSlidingAverage()")
	}
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
	if sc == nil {
		return 0
	}
	return sc.counter[counter]
}
func (sc *sacalc) getCounts(counter int) uint64 {
	if sc == nil {
		return 0
	}
	return sc.counts[counter]
}

func (sc *sacalc) Printf(format string, args ...interface{}) {
	if !*sliding_avg_debug {
		return
	}
	s := "[go-easyops " + sc.debug_name + "] " + format

	fmt.Printf(s, args...)
}
func (sa *SlidingAverage) printf(format string, args ...interface{}) {
	if !*sliding_avg_debug {
		return
	}
	s := "[go-easyops slidingavg] " + format
	fmt.Printf(s, args...)
}

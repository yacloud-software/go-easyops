package utils

import (
	"sync"
	"time"
)

type MiniTimeSeries struct {
	sync.Mutex
	max_keep time.Duration
	values   map[int64]float64
}

// create a new rolling timeseries, which keeps values no longer than 'keep' parameter. (older ones will be discarded)
func NewMiniTimeSeries(keep time.Duration) *MiniTimeSeries {
	res := &MiniTimeSeries{max_keep: keep, values: make(map[int64]float64)}
	return res
}

// called automatically occassionally to delete old values
func (mt *MiniTimeSeries) GC() {
	mt.Lock()
	defer mt.Unlock()
	cutoff := time.Now().Add(0 - mt.max_keep).Unix()
	var deletes []int64
	for i, _ := range mt.values {
		if i <= cutoff {
			deletes = append(deletes, i)
		}
	}
	for _, d := range deletes {
		delete(mt.values, d)
	}

}

// total difference between first and last value (latest - earliest)
func (mt *MiniTimeSeries) Difference() float64 {
	mt.GC()
	_, fl := mt.EarliestValue()
	_, fh := mt.LatestValue()
	return fh - fl
}

func (mt *MiniTimeSeries) EarliestValue() (time.Time, float64) {
	mt.Lock()
	defer mt.Unlock()
	cur_ts := int64(0)
	cur_val := 0.0
	for ts, val := range mt.values {
		if cur_ts == 0 || cur_ts >= ts {
			cur_ts = ts
			cur_val = val
		}
	}
	t := time.Unix(cur_ts, 0)
	return t, cur_val
}
func (mt *MiniTimeSeries) LatestValue() (time.Time, float64) {
	mt.Lock()
	defer mt.Unlock()
	cur_ts := int64(0)
	cur_val := 0.0
	for ts, val := range mt.values {
		if cur_ts == 0 || cur_ts <= ts {
			cur_ts = ts
			cur_val = val
		}
	}
	t := time.Unix(cur_ts, 0)
	return t, cur_val
}
func (mt *MiniTimeSeries) Add(value float64) {
	mt.Lock()
	defer mt.Unlock()
	now := time.Now().Unix()
	mt.values[now] = value

}

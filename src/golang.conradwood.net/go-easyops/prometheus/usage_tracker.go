package prometheus

import (
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"sync"
	"time"
)

var (
	tlock          sync.Mutex
	tracker_arrays = make(map[string][]*metrictracker)
)

type metrictracker struct {
	name        string
	labels      map[string]string
	lastUpdated time.Time
}

func label_map_string(m map[string]string) string {
	s := ""
	for k, v := range m {
		s = s + k + "=" + v + ","
	}
	return s
}
func (p *promRegistry) find_or_create_metric(metricname string, labels map[string]string) *metrictracker {
	tlock.Lock()
	defer tlock.Unlock()
	lms := label_map_string(labels)
	trackers := tracker_arrays[lms]
	for _, t := range trackers {
		if t.name != metricname {
			continue
		}
		if isEqualMap(labels, t.labels) {
			return t
		}
	}
	t := &metrictracker{name: metricname, labels: labels, lastUpdated: time.Now()}
	//	trackers = append(trackers, t)
	tracker_arrays[lms] = append(tracker_arrays[lms], t)
	return t
}
func (p *promRegistry) find_metric(metricname string, labels map[string]string) *metrictracker {
	tlock.Lock()
	defer tlock.Unlock()
	lms := label_map_string(labels)
	trackers := tracker_arrays[lms]
	for _, t := range trackers {
		if t.name != metricname {
			continue
		}
		if isEqualMap(labels, t.labels) {
			return t
		}
	}
	return nil
}
func (m *metrictracker) String() string {
	ls := ""
	deli := ""
	for k, v := range m.labels {
		ls = ls + deli + k + "=\"" + v + "\""
		deli = ","
	}
	return fmt.Sprintf("%s{%s}", m.name, ls)
}
func (m *metrictracker) update() {
	m.lastUpdated = time.Now()
}

// a metric calls this if it is modified
func (p *promRegistry) used(metricname string, labels map[string]string) {
	p.find_or_create_metric(metricname, labels).update()
}
func (p *promRegistry) recently_used_family(mf *dto.MetricFamily, maxage time.Duration) *dto.MetricFamily {
	var rm []*dto.Metric
	for _, m := range mf.Metric {
		l := make(map[string]string)
		for _, lp := range m.Label {
			l[*lp.Name] = l[*lp.Value]
		}
		mt := p.find_metric(*mf.Name, l)
		if mt == nil {
			// not one that has been added
			rm = append(rm, m)
			continue
		}
		if time.Since(mt.lastUpdated) > maxage {
			//		fmt.Printf("Filtered: %s\n", mt.String())
			continue
		}
		rm = append(rm, m)
	}
	mf.Metric = rm
	return mf
}

func isEqualMap(a1, a2 map[string]string) bool {
	if a1 == nil && a2 == nil {
		return true
	}
	if a1 == nil || a2 == nil {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	for k, v := range a1 {
		if a2[k] != v {
			return false
		}
	}
	return true
}

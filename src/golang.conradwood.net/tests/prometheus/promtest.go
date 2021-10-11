package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

type Timer struct {
	start time.Time
	txt   string
}

func NewTimer(txt string) *Timer {
	t := &Timer{txt: txt, start: time.Now()}
	fmt.Printf("Starting \"%s\"...\n", txt)
	return t
}
func (t *Timer) Finish(err error) {
	utils.Bail(fmt.Sprintf("\"%s\" failed", t.txt), err)
	d := time.Since(t.start)
	fmt.Printf("Duration of \"%s\": %0.2f seconds\n", t.txt, d.Seconds())
}

func main() {
	flag.Parse()
	fmt.Printf("Prometheus go-easyops registry test code\n")
	server.StartFakeService("foo.Foo")
	i := 0
	ctr := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("prometheus_test_metric_%d", i),
			Help: fmt.Sprintf("desc_prometheus_test_metric_%d", i),
		},
		[]string{"foo"})
	prometheus.MustRegister(ctr)
	prometheus.GetGatherer().Expiry = 0
	t := NewTimer("Gather() - plain")
	g, err := prometheus.GetRegistry().Gather()
	t.Finish(err)
	t = NewTimer("Gather() - arraylen")
	fmt.Printf("Got %d metric familes\n", len(g))
	t.Finish(nil)

	t = NewTimer("Incrementing counters")
	for i := 0; i < 30000; i++ {
		l := prometheus.Labels{"foo": fmt.Sprintf("%d", i)}
		ctr.With(l).Inc()
	}
	t.Finish(nil)

	t = NewTimer("Gather() - with labels")
	g, err = prometheus.GetRegistry().Gather()
	t.Finish(err)
	fmt.Printf("Got %d metric familes\n", len(g))
	select {}
}

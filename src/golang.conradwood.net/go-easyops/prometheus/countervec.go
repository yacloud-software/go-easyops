package prometheus

import (
	pm "github.com/prometheus/client_golang/prometheus"
)

type CounterVec struct {
	collector
	opts       CounterOpts
	labelnames []string
	gv         *pm.CounterVec
}

func (h *CounterVec) init() *CounterVec {
	h.gv = pm.NewCounterVec(
		pm.CounterOpts{
			Name: h.opts.Name,
			Help: h.opts.Help,
		}, h.labelnames)
	h.setParent(h)
	return h
}
func (g *CounterVec) PMCollector() pm.Collector {
	return g.gv
}
func (g *CounterVec) With(l Labels) pm.Counter {
	promreg.used(g.opts.Name, l)
	return g.gv.With(pm.Labels(l))
}
func (g *CounterVec) Inc() {
	g.gv.With(pm.Labels{}).Inc()
}

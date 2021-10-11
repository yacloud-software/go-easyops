package prometheus

import (
	pm "github.com/prometheus/client_golang/prometheus"
)

type HistogramVec struct {
	collector
	opts       HistogramOpts
	labelnames []string
	gv         *pm.HistogramVec
}

func (h *HistogramVec) init() *HistogramVec {
	h.gv = pm.NewHistogramVec(
		pm.HistogramOpts{
			Name:    h.opts.Name,
			Buckets: h.opts.Buckets,
			Help:    h.opts.Help,
		}, h.labelnames)
	h.setParent(h)
	return h
}
func (g *HistogramVec) PMCollector() pm.Collector {
	return g.gv
}
func (g *HistogramVec) With(l Labels) pm.Observer {
	promreg.used(g.opts.Name, l)
	return g.gv.With(pm.Labels(l))
}

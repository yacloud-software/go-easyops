package prometheus

import (
	pm "github.com/prometheus/client_golang/prometheus"
)

type GaugeVec struct {
	collector
	opts       GaugeOpts
	labelnames []string
	gv         *pm.GaugeVec
}

func (g *GaugeVec) init() *GaugeVec {
	g.gv = pm.NewGaugeVec(
		pm.GaugeOpts{
			Name: g.opts.Name,
			Help: g.opts.Help,
		}, g.labelnames)
	g.setParent(g)
	return g
}
func (g *GaugeVec) PMCollector() pm.Collector {
	return g.gv
}

func (g *GaugeVec) With(l Labels) pm.Gauge {
	promreg.used(g.opts.Name, l)
	return g.gv.With(pm.Labels(l))
}
func (g *GaugeVec) Set(f float64) {
	g.gv.With(pm.Labels{}).Set(f)
}
func (g *GaugeVec) Inc() {
	g.gv.With(pm.Labels{}).Inc()
}
func (g *GaugeVec) Dec() {
	g.gv.With(pm.Labels{}).Dec()
}

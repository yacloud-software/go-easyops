package prometheus

import (
	pm "github.com/prometheus/client_golang/prometheus"
)

type SummaryVec struct {
	collector
	opts       SummaryOpts
	labelnames []string
	gv         *pm.SummaryVec
}

func (h *SummaryVec) init() *SummaryVec {
	h.gv = pm.NewSummaryVec(
		pm.SummaryOpts{
			Name:       h.opts.Name,
			Objectives: h.opts.Objectives,
			Help:       h.opts.Help,
		}, h.labelnames)
	h.setParent(h)
	return h
}
func (g *SummaryVec) PMCollector() pm.Collector {
	return g.gv
}
func (g *SummaryVec) With(l Labels) pm.Observer {
	promreg.used(g.opts.Name, l)
	return g.gv.With(pm.Labels(l))
}
func (g *SummaryVec) Observe(f float64) {
	g.gv.With(pm.Labels{}).Observe(f)
}
func (g *SummaryVec) WithLabelValues(vs ...string) pm.Observer {
	i := 0
	l := make(map[string]string)
	for {
		if i >= len(vs) {
			break
		}
		l[vs[i]] = l[vs[i+1]]
		i = i + 2
	}

	promreg.used(g.opts.Name, l)
	return g.gv.WithLabelValues(vs...)
}

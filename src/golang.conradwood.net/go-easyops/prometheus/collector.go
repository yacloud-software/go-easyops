package prometheus

import (
	pm "github.com/prometheus/client_golang/prometheus"
)

type parent interface {
	PMCollector() pm.Collector
}

type collector struct {
	p   parent
	reg *promRegistry
}

func (c *collector) setParent(p parent) {
	c.p = p
}

func (c *collector) Describe(x chan<- *pm.Desc) {
	c.p.PMCollector().Describe(x)
}
func (c *collector) Collect(x chan<- pm.Metric) {
	c.p.PMCollector().Collect(x)
}

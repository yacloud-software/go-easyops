package prometheus

import (
	"fmt"
	pm "github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	lock    sync.Mutex
	promreg = &promRegistry{Expiry: time.Duration(5) * time.Minute}
)

type Labels pm.Labels

type HistogramOpts pm.HistogramOpts
type GaugeOpts pm.GaugeOpts
type SummaryOpts pm.SummaryOpts
type CounterOpts pm.CounterOpts

/*
type Desc pm.Desc
type Metric pm.Metric

func NewDesc(fqName, help string, variableLabels []string, constLabels Labels) *Desc {
	d := pm.NewDesc(fqName, help, variableLabels, pm.Labels(constLabels))
	return *Desc(d)
}
*/
func NewHistogramVec(opts HistogramOpts, label_names []string) *HistogramVec {
	return (&HistogramVec{opts: opts, labelnames: label_names}).init()
}
func NewHistogram(opts HistogramOpts) *HistogramVec {
	return (&HistogramVec{opts: opts, labelnames: []string{}}).init()
}
func NewSummaryVec(opts SummaryOpts, label_names []string) *SummaryVec {
	return (&SummaryVec{opts: opts, labelnames: label_names}).init()
}
func NewSummary(opts SummaryOpts) *SummaryVec {
	return (&SummaryVec{opts: opts, labelnames: []string{}}).init()
}
func NewCounterVec(opts CounterOpts, label_names []string) *CounterVec {
	return (&CounterVec{opts: opts, labelnames: label_names}).init()
}
func NewCounter(opts CounterOpts) *CounterVec {
	return (&CounterVec{opts: opts, labelnames: []string{}}).init()
}
func NewGaugeVec(opts GaugeOpts, label_names []string) *GaugeVec {
	return (&GaugeVec{opts: opts, labelnames: label_names}).init()
}
func NewGauge(opts GaugeOpts) *GaugeVec {
	return (&GaugeVec{opts: opts, labelnames: []string{}}).init()
}

func MustRegister(cols ...pm.Collector) {
	for _, c := range cols {
		//		c.reg = promreg
		//		promreg.MustRegister(c.p.PMCollector())
		promreg.MustRegister(c)
	}
}
func Register(cols ...pm.Collector) error {
	for _, c := range cols {
		//		c.reg = promreg
		//		promreg.MustRegister(c.p.PMCollector())
		e := promreg.Register(c)
		if e != nil {
			return e
		}
	}
	return nil
}

type promRegistry struct {
	reg    *pm.Registry
	Expiry time.Duration
}

func (p *promRegistry) MustRegister(c pm.Collector) {
	e := p.Register(c)
	if e == nil {
		return
	}
	fmt.Printf("Metric registration of %v failed: %s\n", c, e)
	panic("Metric registration failed")
}

func (p *promRegistry) Register(c pm.Collector) error {
	if p.reg == nil {
		lock.Lock()
		if p.reg == nil {
			p.reg = pm.NewRegistry()
			// install stuff that the normal registry also includes
			p.reg.MustRegister(pm.NewGoCollector())
			p.reg.MustRegister(pm.NewProcessCollector(os.Getpid(), ""))
		}
		lock.Unlock()
	}
	e := p.reg.Register(c)
	return e
}

func GetRegistry() *pm.Registry {
	return promreg.reg
}
func GetGatherer() *promRegistry {
	return promreg
}
func (p *promRegistry) Gather() ([]*dto.MetricFamily, error) {
	var d []*dto.MetricFamily
	dmf, err := p.reg.Gather()
	if err != nil {
		return nil, err
	}
	if p.Expiry == 0 {
		return dmf, nil
	}
	//filter to non expired ones
	for _, mf := range dmf {
		usedonly := p.recently_used_family(mf, p.Expiry)
		if len(usedonly.Metric) > 0 {
			d = append(d, usedonly)
		}
	}

	return d, nil
}

func MetricNames(reg *pm.Registry) ([]string, error) {
	dtm, err := reg.Gather()
	if err != nil {
		return nil, err
	}
	ml := make(map[string]bool)
	for _, mf := range dtm {
		ml[*mf.Name] = true
	}
	var res []string
	for k, _ := range ml {
		res = append(res, k)
	}
	return res, nil
}
func NonstandMetricNames(reg *pm.Registry) ([]string, error) {
	mn, err := MetricNames(reg)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, m := range mn {
		if strings.HasPrefix(m, "go_") {
			continue
		}
		if strings.HasPrefix(m, "process_") {
			continue
		}
		res = append(res, m)
	}
	return res, nil
}

func SetExpiry(expiry time.Duration) {
	promreg.Expiry = expiry
}

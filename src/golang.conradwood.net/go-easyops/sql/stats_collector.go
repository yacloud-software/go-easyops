package sql

import (
	pm "github.com/prometheus/client_golang/prometheus"
	//	"golang.conradwood.net/go-easyops/prometheus"
)

type PoolSizeCollector struct {
	counterDesc *pm.Desc
}

func (c *PoolSizeCollector) Describe(ch chan<- *pm.Desc) {
	ch <- c.counterDesc
}

func (c *PoolSizeCollector) Collect(ch chan<- pm.Metric) {
	for _, db := range databases {
		value := float64(db.dbcon.Stats().OpenConnections)
		ch <- pm.MustNewConstMetric(
			c.counterDesc,
			pm.CounterValue,
			value,
			db.dbname,
		)

	}
}
func NewPoolSizeCollector() *PoolSizeCollector {
	x := &PoolSizeCollector{
		counterDesc: pm.NewDesc("sql_pool_size", "sql pool size",
			[]string{"database"},
			nil),
	}
	return x
}

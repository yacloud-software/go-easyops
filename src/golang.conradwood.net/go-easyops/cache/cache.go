package cache

// a basic local key->value cache with expiry
// the actual implementation of the cache is hidden
// to the user. This is on purpose, so to enable us
// to replace the backend with a distributed cache
// like redis if that becomes benefitial

import (
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/utils"
	"sync"
	"time"
)

var (
	cacheLock   sync.Mutex
	performance = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "goeasyops_cache_performance",
			Help:       "V=1 UNIT=s DESC=Performance of cache lookups",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			MaxAge:     time.Hour,
		},
		[]string{"cachename"},
	)
	size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goeasyops_cache_size",
			Help: "V=1 UNIT=ops DESC=size of cache",
		},
		[]string{"cachename", "et"},
	)
	efficiency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goeasyops_cache_efficiency",
			Help: "V=1 UNIT=ops DESC=hit and miss counters of cache",
		},
		[]string{"cachename", "result"},
	)
	usage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goeasyops_cache_lookups",
			Help: "V=1 UNIT=ops DESC=size of cache",
		},
		[]string{"cachename"},
	)
	caches []*Cache
)

type Cache struct {
	name        string
	mcache      []*cacheEntry
	mlock       sync.Mutex
	MaxLifetime time.Duration
}

type cacheEntry struct {
	free     bool
	created  time.Time
	accessed time.Time
	expiry   *time.Time
	key      string
	value    interface{}
}

func init() {
	prometheus.MustRegister(efficiency, size, performance, usage)
}
func Clear(cacheName string) ([]*Cache, error) {
	fmt.Printf("[go-easyops] Clearing cache \"%s\"\n", cacheName)
	cacheLock.Lock()
	defer cacheLock.Unlock()
	var res []*Cache
	for _, c := range caches {
		if cacheName != "" && c.name != cacheName {
			continue
		}
		c.mlock.Lock()
		c.mcache = make([]*cacheEntry, 0)
		c.mlock.Unlock()
		res = append(res, c)
	}
	return res, nil
}

// create a new cache. "name" must be a prometheus metric compatible name and unique throughout
// good practice: prefix it with servicepackagename. for example:
//  servicename: "lbproxy.LBProxyService"
// -> cachename: "lbproxy_tokencache"
func New(name string, lifetime time.Duration, maxSizeInMB int) *Cache {
	res := &Cache{name: name, MaxLifetime: lifetime}
	res.setCacheGauge(0)
	go res.setCacheGaugeLoop()
	cacheLock.Lock()
	caches = append(caches, res)
	cacheLock.Unlock()
	return res
}
func (c *Cache) Clear() {
	c.mlock.Lock()
	c.mcache = make([]*cacheEntry, 0)
	c.mlock.Unlock()
	c.setCacheGauge(0)
}

func (c *Cache) PutWithExpiry(key string, value interface{}, expiry *time.Time) {
	c.putRaw(key, value, expiry)
}
func (c *Cache) Put(key string, value interface{}) {
	c.putRaw(key, value, nil)
}

func (c *Cache) putRaw(key string, value interface{}, expiry *time.Time) {
	c.mlock.Lock()
	defer c.mlock.Unlock()
	now := time.Now()
	cutOff := time.Now().Add(0 - c.MaxLifetime)
	for _, x := range c.mcache {
		if x.key == key {
			x.created = time.Now()
			x.value = value
			x.accessed = x.created
			x.expiry = expiry
			x.free = false
			return
		}
		if (!x.free) && x.created.Before(cutOff) {
			x.free = true
			continue
		}
		if (x.expiry != nil) && (x.expiry.After(now)) {
			x.free = true
			continue
		}

	}
	for _, x := range c.mcache {
		if x.free {
			x.key = key
			x.created = time.Now()
			x.value = value
			x.expiry = expiry
			x.free = false
			return
		}
	}
	mc := &cacheEntry{free: false, created: time.Now(), expiry: expiry, key: key, value: value}
	mc.accessed = mc.created
	c.mcache = append(c.mcache, mc)
}

func (c *Cache) Get(key string) interface{} {
	label := prometheus.Labels{"cachename": c.name}
	usage.With(label).Inc()
	now := time.Now()
	c.mlock.Lock()
	defer c.mlock.Unlock()
	cutOff := now.Add(0 - c.MaxLifetime)
	for _, x := range c.mcache {
		if (!x.free) && (x.key == key) {
			if x.created.Before(cutOff) {
				x.free = true
				continue
			}
			if (x.expiry != nil) && (x.expiry.After(now)) {
				x.free = true
				continue
			}
			x.accessed = time.Now()
			performance.With(label).Observe(time.Since(now).Seconds())
			efficiency.With(prometheus.Labels{"cachename": c.name, "result": "hit"}).Inc()
			return x.value
		}
	}
	performance.With(label).Observe(time.Since(now).Seconds())
	efficiency.With(prometheus.Labels{"cachename": c.name, "result": "miss"}).Inc()
	return nil
}
func (c *Cache) Keys() []string {
	var res []string
	c.mlock.Lock()
	defer c.mlock.Unlock()
	for _, x := range c.mcache {
		res = append(res, x.key)
	}
	return res
}

func (c *Cache) setCacheGaugeLoop() {
	for {
		c.mlock.Lock()
		i := 0
		for _, x := range c.mcache {
			if x.free {
				continue
			}
			i++
		}
		c.mlock.Unlock()
		c.setCacheGauge(i)
		utils.RandomStall(1)
	}
}
func (c *Cache) setCacheGauge(used int) {
	size.With(prometheus.Labels{"cachename": c.name, "et": "allocated"}).Set(float64(len(c.mcache)))
	size.With(prometheus.Labels{"cachename": c.name, "et": "used"}).Set(float64(used))

}

func (c *Cache) Name() string {
	return c.name
}

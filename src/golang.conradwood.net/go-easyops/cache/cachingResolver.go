package cache

import (
	"context"
	"golang.conradwood.net/go-easyops/prometheus"
	"sync"
	"time"
)

var (
	asyncLookups = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goeasyops_cache_async_lookups",
			Help: "V=1 UNIT=ops DESC=number of looks executed asynchronously",
		},
		[]string{"cachename"},
	)
)

func init() {
	prometheus.MustRegister(asyncLookups)
}

type CachingResolver interface {
	Retrieve(key string, fr func(string) (interface{}, error)) (interface{}, error)
	RetrieveContext(ctx context.Context, key string, fr func(context.Context, string) (interface{}, error)) (interface{}, error)
	SetRefreshAfter(time.Duration)
	SetAsyncRetriever(fr func(string) (interface{}, error))
}

type cachingResolver struct {
	/*
		if true, this will cache errors occured during retrieval and reassert them when the lookup
		is repeated
	*/
	CacheErrors bool
	/*
		if true, this will cache nil results
	*/
	CacheNil bool
	/*
				refresh entries. There is a "hard eviction" noce the lifetime of an object expires. Normally, an object is retrieved
				only when it is requested but no longer exists in the cache.
				However an object might exist in the cache but is "close to" its lifetime. it makes sense to then serve from cache
				and refresh in the background (assuming it will be requested soon again).
		                if this is set, an existing object, older than this will be refreshed in the background after retrieval
	*/
	refreshAfter time.Duration
	/*
		if non-nil this refreshes Errors at a different speed than non-errors
	*/
	refreshErrAfter time.Duration
	// the asynchronous retrieval requires a different codepath, for example, a context might need to be created etc.
	// if this is nil no async retrieval is done, otherwise this function will be used
	asyncRetriever func(string) (interface{}, error)

	// cache
	gccache      *Cache // full of cacheEntry Objects
	retrieveLock sync.Mutex
}

// never exposed outside package. instead "object" is the thing of interest to the user
type cacheEntry2 struct {
	object  interface{}
	err     error
	created time.Time
}

func NewResolvingCache(name string, lifetime time.Duration, maxLimitEntries int) CachingResolver {
	res := &cachingResolver{
		CacheErrors: true,
		CacheNil:    true,
	}
	res.gccache = New(name, lifetime, maxLimitEntries)
	res.refreshAfter = lifetime - (lifetime / 3)
	res.refreshErrAfter = time.Duration(45) * time.Second
	return res
}

/*
retrieves object by key from cache or via retrieval function.
caches result for next time the same key is being looked up
This function _does not_ take a context as parameter.
the resolver might retrieve a value asynchronously whilst synchronously serving from cache
if so, the context might get cancelled before the async retrieval has completed.
** Example: **
o,err :=c.Retrieve("foo", func(k string) (interface{}, error) {
		return "bar",nil
	})
** Example With Async Retrieval: **
c.SetAsyncRetriever(get_by_key)
o,err :=c.Retrieve("foo",get_by_key)
func get_by_key(key string) (interface{}, error) {
 return "bar",nil
"


*/
func (cr *cachingResolver) Retrieve(key string, fr func(string) (interface{}, error)) (interface{}, error) {
	ctx := context.Background()
	return cr.RetrieveContext(ctx, key, func(context.Context, string) (interface{}, error) {
		return fr(key)
	})
}
func (cr *cachingResolver) RetrieveContext(ctx context.Context, key string, fr func(context.Context, string) (interface{}, error)) (interface{}, error) {
	cname := cr.gccache.name
	label := prometheus.Labels{"cachename": cname}
	usage.With(label).Inc()
	started := time.Now()
	var ce *cacheEntry2
	o := cr.gccache.Get(key)
	if o != nil {
		ce = o.(*cacheEntry2)
		if cr.asyncRetriever != nil {
			// we got it in cache. do we need to refresh async?
			if cr.refreshAfter != 0 && time.Since(ce.created) > cr.refreshAfter {
				go cr.refresh(key)
			} else if cr.refreshErrAfter != 0 && ce.err != nil && time.Since(ce.created) > cr.refreshErrAfter {
				go cr.refresh(key)
			}
		}
	}
	if ce == nil {
		// TODO: make this WAY more granular. Goal: do not lookup same key simultanously, but allow different keys to be retrieved simultaneously
		cr.retrieveLock.Lock()
		defer cr.retrieveLock.Unlock()
		o := cr.gccache.Get(key)
		if o != nil {
			ce = o.(*cacheEntry2)
		} else {
			efficiency.With(prometheus.Labels{"cachename": cname, "result": "miss"}).Inc()
			o, err := fr(ctx, key)
			ce = &cacheEntry2{object: o, err: err, created: time.Now()}
			cr.gccache.Put(key, ce)
		}
	} else {
		efficiency.With(prometheus.Labels{"cachename": cname, "result": "hit"}).Inc()
	}
	if ce.err != nil {
		return nil, ce.err
	}
	performance.With(label).Observe(time.Since(started).Seconds())
	return ce.object, nil
}

func (cr *cachingResolver) refresh(key string) {
	cr.retrieveLock.Lock()
	defer cr.retrieveLock.Unlock()
	fr := cr.asyncRetriever
	if fr == nil {
		return
	}
	o, err := fr(key)

	// if we have an error retrieving, put a good object in cache, do not overwrite with error
	if err != nil {
		o := cr.gccache.Get(key)
		if o != nil {
			return
		}
	}
	ce := &cacheEntry2{object: o, err: err, created: time.Now()}
	cr.gccache.Put(key, ce)
	cname := cr.gccache.name
	label := prometheus.Labels{"cachename": cname}
	asyncLookups.With(label).Inc()
}
func (cr *cachingResolver) SetRefreshAfter(d time.Duration) {
	cr.refreshAfter = d
}
func (cr *cachingResolver) SetAsyncRetriever(fr func(string) (interface{}, error)) {
	cr.asyncRetriever = fr
}

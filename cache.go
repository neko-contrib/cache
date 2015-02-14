package cache

import (
	"github.com/rocwong/neko"
	"time"
)

const (
	DefaultExpiration time.Duration = 30 * time.Second

	MemoryStore string = "cache_memory"
	RedisStore  string = "cache_redis"
)

type (
	// Cache is the interface that operates the cache data.
	Cache interface {
		// Set adds an item to the cache, replacing any existing item.
		// If the expire is 0, the cache's default expiration time is used (30 min).
		// If it is nil, the item never expires.
		Set(key string, obj interface{}, expire ...time.Duration)
		// Get gets an item from the cache. returns the item or nil, and a bool indicating
		Get(key string) (interface{}, bool)
		// Delete deletes an item from the cache.
		Delete(key string)
		// Increment increases cached int-type value by given key as a counter.
		Increment(key string, n ...int64) error
		// Decrement decreases cached int-type value by given key as a counter.
		Decrement(key string, n ...int64) error
		// IsExist returns true if cached value exists.
		IsExist(key string) bool
		// Flush deletes all cached data.
		Flush()
		// StartGC starts GC routine based on config string settings.
		StartGC(opt Options)
	}

	Options struct {
		// Store cache store. Default is 'MemoryStore'
		Store string
		// Config stores configuration.
		Config string
		// Interval GC interval time in seconds. Default is 60.
		Interval int
	}
)

var stores = make(map[string]Cache)

func prepareOptions(options []Options) (opt Options) {
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Store) == 0 {
		opt.Store = MemoryStore
	}
	if opt.Interval == 0 {
		opt.Interval = 60
	}
	return
}

func Generate(options ...Options) neko.HandlerFunc {
	opt := prepareOptions(options)
	cache := stores[opt.Store]
	cache.StartGC(opt)
	return func(ctx *neko.Context) {
		ctx.Set(MemoryStore, cache)
	}
}

// Register registers a store.
func Register(name string, store Cache) {
	if store != nil && stores[name] == nil {
		stores[name] = store
	}
}

package cache

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	object     interface{}
	expiration *time.Time
}

type cacher struct {
	sync.RWMutex
	items    map[string]*item
	interval int
}

// newCacher creates and returns a new memory cacher.
func newCacher() *cacher {
	return &cacher{items: make(map[string]*item)}
}

// Expired returns true if the item has expired.
func (c *item) Expired() bool {
	if c.expiration == nil {
		return false
	}
	return c.expiration.Before(time.Now())
}

func (c *cacher) Set(key string, obj interface{}, expire ...time.Duration) {
	c.Lock()
	defer c.Unlock()

	var e *time.Time
	if len(expire) > 0 {
		if expire[0] == 0 {
			expire[0] = DefaultExpiration
		}
		t := time.Now().Add(expire[0])
		e = &t
	}

	c.items[key] = &item{
		object:     obj,
		expiration: e,
	}
}

func (c *cacher) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if item.Expired() {
		go c.Delete(key)
		return nil, false
	}
	return item.object, true
}

func (c *cacher) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *cacher) Increment(key string, n ...int64) error {
	c.RLock()
	defer c.RUnlock()
	var err error
	item, found := c.items[key]
	if !found || item.Expired() {
		return fmt.Errorf("Item %s not found", key)
	}
	if len(n) == 0 {
		n = []int64{1}
	}
	item.object, err = Increment(item.object, n[0])
	return err
}

func (c *cacher) Decrement(key string, n ...int64) error {
	c.RLock()
	defer c.RUnlock()
	var err error
	item, found := c.items[key]
	if !found || item.Expired() {
		return fmt.Errorf("Item %s not found", key)
	}
	if len(n) == 0 {
		n = []int64{1}
	}
	item.object, err = Decrement(item.object, n[0])
	return err
}

func (c *cacher) IsExist(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, found := c.items[key]
	return found
}

func (c *cacher) Flush() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[string]*item)
}

func (c *cacher) StartGC(opt Options) {
	c.interval = opt.Interval
	go c.startGC()
}

func (c *cacher) startGC() {
	if c.interval < 1 {
		return
	}
	if c.items != nil {
		for key, _ := range c.items {
			c.checkExpiration(key)
		}
	}
	time.AfterFunc(time.Duration(c.interval)*time.Second, func() { c.startGC() })
}

func (c *cacher) checkExpiration(key string) {
	c.Lock()
	defer c.Unlock()

	item, found := c.items[key]
	if !found {
		return
	}
	if item.Expired() {
		delete(c.items, key)
	}
}

func init() {
	Register(MemoryStore, newCacher())
}

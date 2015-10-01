package cash

import (
	"sync"
	"time"
)

// Cash configuration.
//
// DefaultExpire: The default expiration time for all items.
// CleanInterval: How often garbage collection should run.
type Conf struct {
	DefaultExpire time.Duration
	CleanInterval time.Duration
}

// Cash contains a map of all the items and holds the
// default expiration and interval times.
type Cash struct {
	sync.RWMutex
	expire   time.Duration
	interval time.Duration
	items    map[string]*item
}

type item struct {
	data   interface{}
	expire time.Time
}

const (
	// Forever means the item will never expire.
	Forever time.Duration = 0
	// Use the default expiration time specified.
	Default time.Duration = -1
)

// New returns a new Cash instance and will also
// garbage collect if the clean interval is bigger
// than -1.
func New(f Conf) *Cash {
	c := &Cash{}
	c.items = map[string]*item{}
	c.interval = f.CleanInterval
	c.expire = f.DefaultExpire
	// Garbage collect
	if c.interval != -1 {
		go c.tick()
	}
	return c
}

// Set will set a value to the cache map.

// To use the default expiration, either pass in 0 or cash.Default,
// to cache it forever pass in -1 or cash.Forever.
func (c *Cash) Set(k string, v interface{}, d time.Duration) {
	c.Lock()
	c.items[k] = c.set(k, v, d)
	c.Unlock()
}

func (c *Cash) set(k string, v interface{}, d time.Duration) *item {
	item := &item{}
	item.data = v

	if d == Default {
		d = c.expire
	}

	if d != Forever {
		item.expire = time.Now().Add(d)
	}

	return item
}

// Get will fetch an item from the cache map and return
// the item along with a boolean, which will be true or
// false depending on whether the item exists.
func (c *Cash) Get(k string) (interface{}, bool) {
	c.RLock()
	item, ok := c.get(k)
	c.RUnlock()
	return item, ok
}

func (c *Cash) get(k string) (interface{}, bool) {
	item, ok := c.items[k]

	if ok && !item.expired() {
		return item.data, true
	}

	return nil, false
}

// Has will return a boolean, which will be true or
// false depending on whether the item exists.
func (c *Cash) Has(k string) bool {
	c.RLock()
	ok := c.has(k)
	c.RUnlock()
	return ok
}

func (c *Cash) has(k string) bool {
	_, ok := c.get(k)
	return ok
}

// Del will attempt to delete an item from
// the cache map. Nothing more, nothing less.
func (c *Cash) Del(k string) {
	c.Lock()
	c.del(k)
	c.Unlock()
}

func (c *Cash) del(k string) {
	delete(c.items, k)
}

// Clean will check all items in the cache map
// and check to see if they have expired, if so
// they will be deleted.
func (c *Cash) Clean() {
	c.Lock()
	c.clean()
	c.Unlock()
}

func (c *Cash) clean() {
	for k, v := range c.items {
		if v.expired() {
			c.del(k)
		}
	}
}

// Flush will delete all items from the cache map.
func (c *Cash) Flush() {
	c.Lock()
	c.flush()
	c.Unlock()
}

func (c *Cash) flush() {
	c.items = map[string]*item{}
}

func (c *Cash) tick() {
	for _ = range time.Tick(c.interval) {
		c.Clean()
	}
}

func (i *item) expired() bool {
	if !i.expire.IsZero() {
		return time.Now().After(i.expire)
	}
	return false
}

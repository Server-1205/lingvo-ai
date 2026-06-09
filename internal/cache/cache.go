package cache

import (
	"sync"
	"time"
)

type entry struct {
	data      interface{}
	expiresAt time.Time
}

type Cache struct {
	mu       sync.RWMutex
	items    map[string]*entry
	defaultTTL time.Duration
}

func New(defaultTTL time.Duration) *Cache {
	c := &Cache{
		items:      make(map[string]*entry),
		defaultTTL: defaultTTL,
	}
	go c.evictLoop()
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}
	return e.data, true
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = &entry{data: data, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *Cache) SetDefault(key string, data interface{}) {
	c.Set(key, data, c.defaultTTL)
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]*entry)
	c.mu.Unlock()
}

func (c *Cache) DeletePrefix(prefix string) {
	c.mu.Lock()
	for k := range c.items {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

func (c *Cache) evictLoop() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.items {
			if now.After(e.expiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]CacheEntry
	mutex    sync.RWMutex
	duration time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		entries:  make(map[string]CacheEntry),
		duration: interval,
	}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	for range ticker.C {
		c.mutex.Lock()
		for k, v := range c.entries {
			if time.Since(v.createdAt) > c.duration {
				delete(c.entries, k)
			}
		}
		c.mutex.Unlock()
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

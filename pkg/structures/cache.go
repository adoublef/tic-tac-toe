package structures

import (
	"log"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	// Get returns the value for the given key and a boolean indicating if the key was found.
	Get(key K) (V, bool)
	// Put adds the given key and value to the cache. If the key already exists, it will be overwritten.
	// If expires is 0, the entry will never expire.
	Put(key K, value V, expires time.Duration)
	// Delete removes the given key from the cache.
	Delete(key K)
}

type cacheEntry[V any] struct {
	// exp is the expiration time in unix nano.
	exp int64
	// v is the value.
	v V
}

type cache[K comparable, V any] struct {
	// mu is a mutex for the cache.
	mu sync.Mutex
	// m is the map of cache entries.
	m map[K]cacheEntry[V]
}

func NewCache[K comparable, V any]() Cache[K, V] {
	c := &cache[K, V]{m: make(map[K]cacheEntry[V])}
	go c.listen()
	return c
}

func (c *cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.m[key]; ok {
		//
		if v.exp == 0 || time.Now().UnixNano() < v.exp {
			return v.v, true
		}
	}
	var zero V
	return zero, false
}

func (c *cache[K, V]) Put(key K, value V, expires time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if expires == 0 {
		c.m[key] = cacheEntry[V]{exp: 0, v: value}
		return
	}

	c.m[key] = cacheEntry[V]{exp: time.Now().Add(expires).UnixNano(), v: value}
}

func (c *cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

func (c *cache[K, V]) listen() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for range t.C {
		c.mu.Lock()
		for k, v := range c.m {
			log.Printf("entry=(%v, %v); time=%v", k, v, time.Now().UnixNano())
			if v.exp != 0 && time.Now().UnixNano() > v.exp {
				log.Printf("deleting entry=(%v, %v)", k, v)
				delete(c.m, k)
			}
		}
		c.mu.Unlock()
	}
}

package pokecache

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mutex   sync.Mutex
}

func GetPokeapi(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	results, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func NewCache(interval time.Duration) *Cache {
	entries := make(map[string]cacheEntry)
	cache := Cache{
		entries: entries,
		mutex:   sync.Mutex{},
	}
	go cache.ReapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.entries[key]
	if !ok {
		val, err := GetPokeapi(key)
		if err != nil {
			fmt.Println("failed the api call")
		}
		c.entries[key] = cacheEntry{
			createdAt: time.Now(),
			val:       val,
		}
	}
	return c.entries[key].val, ok
}

func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(interval)
		done <- true
	}()
	for range ticker.C {
		c.mutex.Lock()
		for k, v := range c.entries {
			if v.createdAt.Add(interval).Before(time.Now()) {
				delete(c.entries, k)
			}
		}
		c.mutex.Unlock()
	}
}

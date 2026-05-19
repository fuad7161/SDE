package mutex

import "sync"

type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]string),
	}
}

package cloudinary_gem

import "sync"

type BreakpointCache interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Flush()
}

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryCache() *MemoryCache { return &MemoryCache{data: map[string]string{}} }
func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}
func (c *MemoryCache) Set(key, value string) { c.mu.Lock(); defer c.mu.Unlock(); c.data[key] = value }
func (c *MemoryCache) Flush()                { c.mu.Lock(); defer c.mu.Unlock(); c.data = map[string]string{} }
func BreakpointCacheKey(publicID, resourceType, deliveryType, transformation, format string) string {
	return publicID + ":" + resourceType + ":" + deliveryType + ":" + transformation + ":" + format
}

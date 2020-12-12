package cache

import "time"

const maxAge = "15m"

var CacheMaxAge, _ = time.ParseDuration(maxAge)

type Cache interface {
	InCache(string) bool
	Get(string) (*CacheEntry, error)
	Put(string, *CacheEntry) error
	StartEngine()
}

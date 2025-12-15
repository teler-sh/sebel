// Package cache provides a generic persistent cache using Otter.
package cache

import (
	"os"
	"sync"
	"time"

	"github.com/maypok86/otter/v2"
)

// Options configures a Cache instance.
// It embeds otter.Options and adds persistence support.
type Options[K comparable, V any] struct {
	otter.Options[K, V]

	// Filepath is the path for persisting the cache.
	// If empty, persistence is disabled.
	Filepath string
}

// Cache is a generic persistent cache that embeds [otter.Cache].
// It adds persistence and background refresh capabilities.
type Cache[K comparable, V any] struct {
	*otter.Cache[K, V]

	cacheFilepath string

	refreshStop chan struct{}
	refreshOnce sync.Once
	refreshMu   sync.Mutex
}

// New creates a new [Cache] instance with the given options.
func New[K comparable, V any](opts *Options[K, V]) (*Cache[K, V], error) {
	otterCache, err := otter.New(&opts.Options)
	if err != nil {
		return nil, err
	}

	c := &Cache[K, V]{
		Cache:         otterCache,
		cacheFilepath: opts.Filepath,
	}

	if c.cacheFilepath != "" {
		if _, err := os.Stat(c.cacheFilepath); err == nil {
			_ = otter.LoadCacheFromFile(c.Cache, c.cacheFilepath)
		}
	}

	return c, nil
}

// Must is like [New] but panics on error.
func Must[K comparable, V any](opts *Options[K, V]) *Cache[K, V] {
	c, err := New(opts)
	if err != nil {
		panic(err)
	}

	return c
}

// Save persists the cache to disk.
// Returns nil if persistence is not configured.
func (c *Cache[K, V]) Save() error {
	if c.cacheFilepath == "" {
		return nil
	}

	return otter.SaveCacheToFile(c.Cache, c.cacheFilepath)
}

// Load loads the cache from disk.
// Returns nil if persistence is not configured.
func (c *Cache[K, V]) Load() error {
	if c.cacheFilepath == "" {
		return nil
	}

	return otter.LoadCacheFromFile(c.Cache, c.cacheFilepath)
}

// FetchFunc is a function that fetches fresh data and populates the cache.
type FetchFunc func() error

// StartBackgroundRefresh starts a goroutine that periodically calls fetchFn.
// Call StopBackgroundRefresh to stop.
func (c *Cache[K, V]) StartBackgroundRefresh(interval time.Duration, fetchFn FetchFunc) {
	if interval <= 0 || fetchFn == nil {
		return
	}

	c.refreshOnce.Do(func() {
		c.refreshStop = make(chan struct{})
		go c.backgroundRefresh(interval, fetchFn)
	})
}

// StopBackgroundRefresh stops the background refresh goroutine.
func (c *Cache[K, V]) StopBackgroundRefresh() {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	if c.refreshStop != nil {
		close(c.refreshStop)
		c.refreshStop = nil
		c.refreshOnce = sync.Once{}
	}
}

// backgroundRefresh periodically calls fetchFn.
func (c *Cache[K, V]) backgroundRefresh(interval time.Duration, fetchFn FetchFunc) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = fetchFn()
		case <-c.refreshStop:
			return
		}
	}
}

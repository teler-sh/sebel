package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/maypok86/otter/v2"
)

func TestCacheSetAndGet(t *testing.T) {
	c := newTestCache[string, string](t)

	c.Set("key1", "value1")

	got, ok := c.GetIfPresent("key1")
	if !ok {
		t.Fatal("key not found in cache")
	}

	if got != "value1" {
		t.Errorf("value mismatch: got %s, want %s", got, "value1")
	}
}

func TestCacheSize(t *testing.T) {
	c := newTestCache[string, int](t)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	if c.EstimatedSize() != 3 {
		t.Errorf("expected size 3, got %d", c.EstimatedSize())
	}
}

func TestCacheInvalidate(t *testing.T) {
	c := newTestCache[string, string](t)

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	if c.EstimatedSize() == 0 {
		t.Fatal("cache should have entries before invalidate")
	}

	c.InvalidateAll()

	if c.EstimatedSize() != 0 {
		t.Errorf("cache should be empty after invalidate, got size %d", c.EstimatedSize())
	}
}

func TestCachePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.gob")

	// Create and populate cache
	c1, err := New(&Options[string, string]{Filepath: tmpFile})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	c1.Set("persist1", "value1")
	c1.Set("persist2", "value2")

	if err := c1.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("cache file was not created")
	}

	// Create new cache and load
	c2, err := New(&Options[string, string]{Filepath: tmpFile})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if err := c2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify data was loaded
	if got, ok := c2.GetIfPresent("persist1"); !ok || got != "value1" {
		t.Errorf("persist1: got %q, want %q", got, "value1")
	}
	if got, ok := c2.GetIfPresent("persist2"); !ok || got != "value2" {
		t.Errorf("persist2: got %q, want %q", got, "value2")
	}
}

func TestCacheAll(t *testing.T) {
	c := newTestCache[string, int](t)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	count := 0
	sum := 0
	for _, v := range c.All() {
		count++
		sum += v
	}

	if count != 3 {
		t.Errorf("expected 3 iterations, got %d", count)
	}
	if sum != 6 {
		t.Errorf("expected sum 6, got %d", sum)
	}
}

func TestCacheWithOtterOptions(t *testing.T) {
	tmpDir := t.TempDir()

	c, err := New(&Options[string, string]{
		Options: otter.Options[string, string]{
			InitialCapacity: 1000,
		},
		Filepath: filepath.Join(tmpDir, "test.gob"),
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	c.Set("key", "value")

	got, ok := c.GetIfPresent("key")
	if !ok || got != "value" {
		t.Errorf("got %q, want %q", got, "value")
	}
}

// newTestCache creates a cache with a temp directory for testing.
func newTestCache[K comparable, V any](t *testing.T) *Cache[K, V] {
	t.Helper()

	tmpDir := t.TempDir()
	c, err := New(&Options[K, V]{
		Filepath: filepath.Join(tmpDir, "test.gob"),
	})
	if err != nil {
		t.Fatalf("failed to create test cache: %v", err)
	}

	return c
}

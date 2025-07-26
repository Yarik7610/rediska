package memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringStorageSet(t *testing.T) {
	storage := NewStringStorage()

	t.Run("simple set", func(t *testing.T) {
		storage.Set("key1", "value1")
		item, ok := storage.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", item.Value)
	})

	t.Run("overwrite value", func(t *testing.T) {
		storage.Set("key1", "value1")
		storage.Set("key1", "value2")
		item, ok := storage.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value2", item.Value)
	})
}

func TestStringStorageSetWithExpiry(t *testing.T) {
	storage := NewStringStorage()

	t.Run("set with expiry", func(t *testing.T) {
		storage.SetWithExpiry("key1", "value1", time.Now().Add(100*time.Millisecond))
		item, ok := storage.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", item.Value)
	})

	t.Run("expired key", func(t *testing.T) {
		storage.SetWithExpiry("key2", "value2", time.Now().Add(-time.Second))
		item, ok := storage.Get("key2")
		assert.False(t, ok)
		assert.Nil(t, item)
	})

	t.Run("zero expiry", func(t *testing.T) {
		storage.SetWithExpiry("key3", "value3", time.Time{})
		item, ok := storage.Get("key3")
		assert.False(t, ok)
		assert.Nil(t, item)
	})
}

func TestStringStorageGet(t *testing.T) {
	storage := NewStringStorage()

	t.Run("non-existent key", func(t *testing.T) {
		item, ok := storage.Get("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, item)
	})

	t.Run("get after set", func(t *testing.T) {
		storage.Set("key1", "value1")
		item, ok := storage.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", item.Value)
	})
}

func TestStringStorageGetKeys(t *testing.T) {
	storage := NewStringStorage()

	t.Run("empty storage", func(t *testing.T) {
		keys := storage.GetKeys()
		assert.Empty(t, keys)
	})

	t.Run("with keys", func(t *testing.T) {
		storage.Set("key1", "value1")
		storage.SetWithExpiry("key2", "value2", time.Now().Add(100*time.Millisecond))
		keys := storage.GetKeys()
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
	})

	t.Run("after expiry", func(t *testing.T) {
		storage.SetWithExpiry("key3", "value3", time.Now().Add(50*time.Millisecond))
		time.Sleep(100 * time.Millisecond)
		keys := storage.GetKeys()
		assert.NotContains(t, keys, "key3")
	})
}

func TestStringStorageCleanExpiredKeys(t *testing.T) {
	storage := NewStringStorage()

	t.Run("clean expired", func(t *testing.T) {
		storage.SetWithExpiry("key1", "value1", time.Now().Add(50*time.Millisecond))
		storage.Set("key2", "value2")
		time.Sleep(100 * time.Millisecond)
		storage.CleanExpiredKeys()

		_, ok := storage.Get("key1")
		assert.False(t, ok)
		item, ok := storage.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", item.Value)
	})
}

func TestStringStorageDel(t *testing.T) {
	storage := NewStringStorage()

	t.Run("delete non-existent", func(t *testing.T) {
		storage.Del("nonexistent")
		_, ok := storage.Get("nonexistent")
		assert.False(t, ok)
	})

	t.Run("delete existing", func(t *testing.T) {
		storage.Set("key1", "value1")
		storage.Del("key1")
		_, ok := storage.Get("key1")
		assert.False(t, ok)
	})
}

func TestStringStorageConcurrent(t *testing.T) {
	storage := NewStringStorage()
	const count = 100

	t.Run("concurrent set and get", func(t *testing.T) {
		var wg sync.WaitGroup
		keys := make([]string, count)
		values := make([]string, count)

		for i := range count {
			keys[i] = fmt.Sprintf("key%d", i)
			values[i] = fmt.Sprintf("value%d", i)
			wg.Add(2)
		}

		for i := range count {
			go func(idx int) {
				defer wg.Done()
				storage.Set(keys[idx], values[idx])
			}(i)

			go func(idx int) {
				defer wg.Done()
				item, ok := storage.Get(keys[idx])
				if ok {
					assert.Equal(t, values[idx], item.Value)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("concurrent overwrite", func(t *testing.T) {
		var wg sync.WaitGroup
		key := "concurrent_key"
		expectedValues := make([]string, count)

		for i := range count {
			expectedValues[i] = fmt.Sprintf("value%d", i)
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				storage.Set(key, expectedValues[idx])
			}(i)
		}
		wg.Wait()

		item, ok := storage.Get(key)
		assert.True(t, ok)
		assert.Contains(t, expectedValues, item.Value)
	})
}

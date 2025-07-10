package memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	t.Run("Concurrent Set and Get mixed", func(t *testing.T) {
		storage := NewStorage()
		var wg sync.WaitGroup

		count := 100
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

	t.Run("Concurrent SetWithExpiry and Get mixed", func(t *testing.T) {
		storage := NewStorage()
		var wg sync.WaitGroup

		count := 100
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
				storage.SetWithExpiry(keys[idx], values[idx], 100*time.Millisecond)
			}(i)

			go func(idx int) {
				defer wg.Done()
				item, ok := storage.Get(keys[idx])
				if ok {
					assert.Equal(t, values[idx], item.Value)
					assert.False(t, item.Expires.IsZero())
				}
			}(i)
		}
		wg.Wait()

		time.Sleep(150 * time.Millisecond)
		for i := range count {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				item, ok := storage.Get(keys[idx])
				assert.False(t, ok)
				assert.Nil(t, item)
			}(i)
		}
		wg.Wait()
	})

	t.Run("GetKeys with mixed keys", func(t *testing.T) {
		storage := NewStorage()
		storage.Set("key1", "value1")
		storage.SetWithExpiry("key2", "value2", 50*time.Millisecond)
		storage.SetWithExpiry("key3", "value3", -1)

		time.Sleep(100 * time.Millisecond)

		keys := storage.GetKeys()
		assert.Equal(t, []string{"key1"}, keys)

		_, ok := storage.Get("key2")
		assert.False(t, ok)
		_, ok = storage.Get("key3")
		assert.False(t, ok)
	})

	t.Run("CleanExpiredKeys", func(t *testing.T) {
		storage := NewStorage()

		storage.SetWithExpiry("key1", "value1", 50*time.Millisecond)
		storage.Set("key2", "value2")

		time.Sleep(100 * time.Millisecond)

		storage.CleanExpiredKeys()

		_, ok := storage.Get("key1")
		assert.False(t, ok)

		item, ok := storage.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", item.Value)
	})

	t.Run("Non-existent key", func(t *testing.T) {
		storage := NewStorage()

		item, ok := storage.Get("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, item)
	})

	t.Run("Concurrent overwrite", func(t *testing.T) {
		storage := NewStorage()

		var wg sync.WaitGroup
		key := "key"
		count := 100

		expectedValues := make([]string, count)

		for i := range count {
			expectedValues[i] = fmt.Sprintf("value%d", i)
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				value := fmt.Sprintf("value%d", idx)
				storage.Set(key, value)
			}(i)
		}
		wg.Wait()

		item, ok := storage.Get(key)
		assert.True(t, ok)
		assert.Contains(t, expectedValues, item.Value)
	})

	t.Run("SetWithExpiry with zero duration", func(t *testing.T) {
		storage := NewStorage()
		storage.SetWithExpiry("key", "value", 0)

		item, ok := storage.Get("key")
		assert.False(t, ok)
		assert.Nil(t, item)
	})
}

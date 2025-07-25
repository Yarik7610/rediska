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
		stringStorage := NewStringStorage()
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
				stringStorage.Set(keys[idx], values[idx])
			}(i)

			go func(idx int) {
				defer wg.Done()
				item, ok := stringStorage.Get(keys[idx])
				if ok {
					assert.Equal(t, values[idx], item.Value)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("Concurrent SetWithExpiry and Get mixed", func(t *testing.T) {
		stringStorage := NewStringStorage()
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
				stringStorage.SetWithExpiry(keys[idx], values[idx], time.Now().Add(100*time.Millisecond))
			}(i)

			go func(idx int) {
				defer wg.Done()
				item, ok := stringStorage.Get(keys[idx])
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
				item, ok := stringStorage.Get(keys[idx])
				assert.False(t, ok)
				assert.Nil(t, item)
			}(i)
		}
		wg.Wait()
	})

	t.Run("GetKeys with mixed keys", func(t *testing.T) {
		stringStorage := NewStringStorage()
		stringStorage.Set("key1", "value1")
		stringStorage.SetWithExpiry("key2", "value2", time.Now().Add(50*time.Millisecond))
		stringStorage.SetWithExpiry("key3", "value3", time.Now().Add(-1))

		time.Sleep(100 * time.Millisecond)

		keys := stringStorage.GetKeys()
		assert.Equal(t, []string{"key1"}, keys)

		_, ok := stringStorage.Get("key2")
		assert.False(t, ok)
		_, ok = stringStorage.Get("key3")
		assert.False(t, ok)
	})

	t.Run("CleanExpiredKeys", func(t *testing.T) {
		stringStorage := NewStringStorage()

		stringStorage.SetWithExpiry("key1", "value1", time.Now().Add(50*time.Millisecond))
		stringStorage.Set("key2", "value2")

		time.Sleep(100 * time.Millisecond)

		stringStorage.CleanExpiredKeys()

		_, ok := stringStorage.Get("key1")
		assert.False(t, ok)

		item, ok := stringStorage.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", item.Value)
	})

	t.Run("Non-existent key", func(t *testing.T) {
		stringStorage := NewStringStorage()

		item, ok := stringStorage.Get("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, item)
	})

	t.Run("Concurrent overwrite", func(t *testing.T) {
		stringStorage := NewStringStorage()

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
				stringStorage.Set(key, value)
			}(i)
		}
		wg.Wait()

		item, ok := stringStorage.Get(key)
		assert.True(t, ok)
		assert.Contains(t, expectedValues, item.Value)
	})

	t.Run("SetWithExpiry with zero expiry", func(t *testing.T) {
		stringStorage := NewStringStorage()
		stringStorage.SetWithExpiry("key", "value", time.Now().Add(0))

		item, ok := stringStorage.Get("key")
		assert.False(t, ok)
		assert.Nil(t, item)
	})
}

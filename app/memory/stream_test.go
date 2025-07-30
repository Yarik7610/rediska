package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamStorageXadd(t *testing.T) {
	ss := NewStreamStorage()

	t.Run("basic add", func(t *testing.T) {
		id, err := ss.Xadd("mystream", "1-0", map[string]string{"field1": "value1"})
		assert.NoError(t, err)
		assert.True(t, ss.Has("mystream"))
		assert.Equal(t, "1-0", id)
	})

	t.Run("duplicate ID", func(t *testing.T) {
		_, err := ss.Xadd("mystream", "1-0", map[string]string{"field2": "value2"})
		assert.Error(t, err)
		assert.Equal(t, "entry with such stream ID already exists", err.Error())
	})

	t.Run("empty fields", func(t *testing.T) {
		_, err := ss.Xadd("mystream", "2-0", map[string]string{})
		assert.Error(t, err)
	})
}

func TestStreamStorageXaddConcurrent(t *testing.T) {
	ss := NewStreamStorage()
	const workers = 10
	var wg sync.WaitGroup

	t.Run("concurrent stream creation", func(t *testing.T) {
		wg.Add(workers)
		for i := range workers {
			go func(idx int) {
				defer wg.Done()
				streamID := fmt.Sprintf("1-%d", idx)
				streamName := fmt.Sprintf("stream%d", idx)
				_, err := ss.Xadd(streamName, streamID, map[string]string{"a": "b"})
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		assert.Equal(t, workers, len(ss.Keys()))
	})
}

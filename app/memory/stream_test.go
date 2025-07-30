package memory

import (
	"fmt"
	"strings"
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
		assert.Equal(t, "The ID specified in XADD is equal or smaller than the target stream top item", err.Error())
	})

	t.Run("empty fields", func(t *testing.T) {
		_, err := ss.Xadd("mystream", "2-0", map[string]string{})
		assert.Error(t, err)
	})

	t.Run("auto generate ID with *", func(t *testing.T) {
		id, err := ss.Xadd("autostream", "*", map[string]string{"auto": "generated"})
		assert.NoError(t, err)
		assert.True(t, strings.Contains(id, "-"))
	})

	t.Run("invalid ID format", func(t *testing.T) {
		_, err := ss.Xadd("badstream", "invalid", map[string]string{"a": "b"})
		assert.Error(t, err)
	})

	t.Run("invalid milliseconds time", func(t *testing.T) {
		_, err := ss.Xadd("badstream", "abc-0", map[string]string{"a": "b"})
		assert.Error(t, err)
	})

	t.Run("invalid sequence number", func(t *testing.T) {
		_, err := ss.Xadd("badstream", "1-abc", map[string]string{"a": "b"})
		assert.Error(t, err)
	})

	t.Run("zero ID", func(t *testing.T) {
		_, err := ss.Xadd("zerostream", "0-0", map[string]string{"a": "b"})
		assert.Error(t, err)
		assert.Equal(t, "The ID specified in XADD must be greater than 0-0", err.Error())
	})

	t.Run("auto sequence number with *", func(t *testing.T) {
		_, err := ss.Xadd("seqstream", "1000-0", map[string]string{"a": "b"})
		assert.NoError(t, err)

		id, err := ss.Xadd("seqstream", "1000-*", map[string]string{"b": "c"})
		assert.NoError(t, err)
		assert.Equal(t, "1000-1", id)

		id, err = ss.Xadd("seqstream", "1000-*", map[string]string{"c": "d"})
		assert.NoError(t, err)
		assert.Equal(t, "1000-2", id)

		id, err = ss.Xadd("seqstream", "1001-*", map[string]string{"d": "e"})
		assert.NoError(t, err)
		assert.Equal(t, "1001-0", id)
	})

	t.Run("smaller than top entry", func(t *testing.T) {
		_, err := ss.Xadd("orderstream", "2000-0", map[string]string{"a": "b"})
		assert.NoError(t, err)

		_, err = ss.Xadd("orderstream", "1999-0", map[string]string{"b": "c"})
		assert.Error(t, err)
		assert.Equal(t, "The ID specified in XADD is equal or smaller than the target stream top item", err.Error())

		_, err = ss.Xadd("orderstream", "2000-0", map[string]string{"b": "c"})
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

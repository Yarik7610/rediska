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
func TestStreamStorageXrange(t *testing.T) {
	ss := NewStreamStorage()

	_, err := ss.Xadd("rangestream", "1000-0", map[string]string{"field1": "value1"})
	assert.NoError(t, err)
	_, err = ss.Xadd("rangestream", "1000-1", map[string]string{"field2": "value2"})
	assert.NoError(t, err)
	_, err = ss.Xadd("rangestream", "1001-0", map[string]string{"field3": "value3"})
	assert.NoError(t, err)
	_, err = ss.Xadd("rangestream", "1002-0", map[string]string{"field4": "value4"})
	assert.NoError(t, err)

	t.Run("basic range", func(t *testing.T) {
		entries, err := ss.Xrange("rangestream", "1000-0", "1001-0")
		assert.NoError(t, err)
		assert.Len(t, entries, 3)
		assert.Equal(t, "1000-0", entries[0].StreamID)
		assert.Equal(t, "value1", entries[0].Entry["field1"])
		assert.Equal(t, "1000-1", entries[1].StreamID)
		assert.Equal(t, "1001-0", entries[2].StreamID)
	})

	t.Run("single entry range", func(t *testing.T) {
		entries, err := ss.Xrange("rangestream", "1000-1", "1000-1")
		assert.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Equal(t, "1000-1", entries[0].StreamID)
	})

	t.Run("full range with special IDs", func(t *testing.T) {
		entries, err := ss.Xrange("rangestream", "-", "+")
		assert.NoError(t, err)
		assert.Len(t, entries, 4)
		assert.Equal(t, "1000-0", entries[0].StreamID)
		assert.Equal(t, "1002-0", entries[3].StreamID)
	})

	t.Run("non-existing entries in range", func(t *testing.T) {
		entries, err := ss.Xrange("rangestream", "1005-0", "1010-0")
		assert.NoError(t, err)
		assert.Len(t, entries, 0)
	})

	t.Run("invalid start ID format", func(t *testing.T) {
		_, err := ss.Xrange("rangestream", "invalid", "1001-0")
		assert.Error(t, err)
	})

	t.Run("invalid end ID format", func(t *testing.T) {
		_, err := ss.Xrange("rangestream", "1000-0", "invalid")
		assert.Error(t, err)
	})

	t.Run("start after end", func(t *testing.T) {
		entries, err := ss.Xrange("rangestream", "1001-0", "1000-0")
		assert.NoError(t, err)
		assert.Len(t, entries, 0)
	})
}

func TestStreamStorageXrangeConcurrent(t *testing.T) {
	ss := NewStreamStorage()
	const workers = 10
	var wg sync.WaitGroup

	for i := range workers {
		streamID := fmt.Sprintf("1000-%d", i)
		_, err := ss.Xadd("concurrentstream", streamID, map[string]string{"value": fmt.Sprintf("%d", i)})
		assert.NoError(t, err)
	}

	t.Run("concurrent reads", func(t *testing.T) {
		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				entries, err := ss.Xrange("concurrentstream", "1000-0", "1000-9")
				assert.NoError(t, err)
				assert.Len(t, entries, 10)
			}()
		}
		wg.Wait()
	})
}

func TestStreamStorageXread(t *testing.T) {
	ss := NewStreamStorage()

	_, err := ss.Xadd("stream1", "1000-0", map[string]string{"field1": "value1"})
	assert.NoError(t, err)
	_, err = ss.Xadd("stream1", "1000-1", map[string]string{"field2": "value2"})
	assert.NoError(t, err)
	_, err = ss.Xadd("stream2", "2000-0", map[string]string{"field3": "value3"})
	assert.NoError(t, err)

	t.Run("read all from one stream", func(t *testing.T) {
		result, err := ss.Xread([]string{"stream1"}, []string{"0-0"})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "stream1", result[0].StreamKey)
		assert.Len(t, result[0].EntriesWithStreamID, 2)
		assert.Equal(t, "1000-0", result[0].EntriesWithStreamID[0].StreamID)
		assert.Equal(t, "1000-1", result[0].EntriesWithStreamID[1].StreamID)
	})

	t.Run("read from multiple streams", func(t *testing.T) {
		result, err := ss.Xread([]string{"stream1", "stream2"}, []string{"1000-1", "0-0"})
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		assert.Equal(t, "stream1", result[0].StreamKey)
		assert.Len(t, result[0].EntriesWithStreamID, 0)

		assert.Equal(t, "stream2", result[1].StreamKey)
		assert.Len(t, result[1].EntriesWithStreamID, 1)
		assert.Equal(t, "2000-0", result[1].EntriesWithStreamID[0].StreamID)
	})

	t.Run("read with non-existing start ID", func(t *testing.T) {
		result, err := ss.Xread([]string{"stream1"}, []string{"9999-0"})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Len(t, result[0].EntriesWithStreamID, 0)
	})

	t.Run("read from non-existing stream", func(t *testing.T) {
		result, err := ss.Xread([]string{"nonexistent"}, []string{"0-0"})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Len(t, result[0].EntriesWithStreamID, 0)
	})
}

func TestStreamStorageXreadConcurrent(t *testing.T) {
	ss := NewStreamStorage()
	const workers = 10
	var wg sync.WaitGroup

	for i := range workers {
		streamID := fmt.Sprintf("1000-%d", i)
		_, err := ss.Xadd("concurrentstream", streamID, map[string]string{"value": fmt.Sprintf("%d", i)})
		assert.NoError(t, err)
	}

	t.Run("concurrent reads", func(t *testing.T) {
		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				result, err := ss.Xread([]string{"concurrentstream"}, []string{"1000-0"})
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Len(t, result[0].EntriesWithStreamID, workers-1)
			}()
		}
		wg.Wait()
	})
}

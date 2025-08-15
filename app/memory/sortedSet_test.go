package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSetStorageZaddAndZcard(t *testing.T) {
	ss := NewSortedSetStorage()

	t.Run("add new elements", func(t *testing.T) {
		added := ss.Zadd("myzset", []float64{1, 2, 3}, []string{"a", "b", "c"})
		assert.Equal(t, 3, added)
		assert.Equal(t, 3, ss.Zcard("myzset"))
	})

	t.Run("update existing element", func(t *testing.T) {
		added := ss.Zadd("myzset", []float64{5}, []string{"b"})
		assert.Equal(t, 0, added)
		score := ss.Zscore("myzset", "b")
		assert.NotNil(t, score)
		assert.Equal(t, 5.0, *score)
	})
}

func TestSortedSetStorageZrangeAndZrank(t *testing.T) {
	ss := NewSortedSetStorage()
	ss.Zadd("myzset", []float64{1, 2, 3}, []string{"a", "b", "c"})

	t.Run("zrange full", func(t *testing.T) {
		values := ss.Zrange("myzset", 0, -1)
		assert.Equal(t, []string{"a", "b", "c"}, values)
	})

	t.Run("zrange partial", func(t *testing.T) {
		values := ss.Zrange("myzset", 1, 2)
		assert.Equal(t, []string{"b", "c"}, values)
	})

	t.Run("zrank existing", func(t *testing.T) {
		rank := ss.Zrank("myzset", "b")
		assert.Equal(t, 1, rank)
	})

	t.Run("zrank non-existing", func(t *testing.T) {
		rank := ss.Zrank("myzset", "non")
		assert.Equal(t, -1, rank)
	})
}

func TestSortedSetStorageZrem(t *testing.T) {
	ss := NewSortedSetStorage()
	ss.Zadd("myzset", []float64{1, 2, 3}, []string{"a", "b", "c"})

	t.Run("remove existing elements", func(t *testing.T) {
		removed := ss.Zrem("myzset", []string{"b", "c"})
		assert.Equal(t, 2, removed)
		assert.Equal(t, 1, ss.Zcard("myzset"))
	})

	t.Run("remove non-existing element", func(t *testing.T) {
		removed := ss.Zrem("myzset", []string{"x"})
		assert.Equal(t, 0, removed)
	})
}

func TestSortedSetStorageKeysHasDel(t *testing.T) {
	ss := NewSortedSetStorage()

	t.Run("keys and has", func(t *testing.T) {
		ss.Zadd("myzset1", []float64{1}, []string{"a"})
		ss.Zadd("myzset2", []float64{2}, []string{"b"})
		keys := ss.Keys()
		assert.Contains(t, keys, "myzset1")
		assert.Contains(t, keys, "myzset2")
		assert.True(t, ss.Has("myzset1"))
		assert.False(t, ss.Has("nonexistent"))
	})

	t.Run("del key", func(t *testing.T) {
		ss.Del("myzset1")
		assert.False(t, ss.Has("myzset1"))
	})
}

func TestSortedSetStorageConcurrentAccess(t *testing.T) {
	ss := NewSortedSetStorage()
	const workers = 10
	var wg sync.WaitGroup

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			key := "concurrent_zset"
			member := fmt.Sprintf("member%d", idx)
			ss.Zadd(key, []float64{float64(idx)}, []string{member})
			rank := ss.Zrank(key, member)
			assert.GreaterOrEqual(t, rank, 0)
			score := ss.Zscore(key, member)
			assert.NotNil(t, score)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, workers, ss.Zcard("concurrent_zset"))
	values := ss.Zrange("concurrent_zset", 0, -1)
	assert.Len(t, values, workers)
}

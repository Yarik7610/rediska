package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipList(t *testing.T) {
	t.Run("InsertSingle", func(t *testing.T) {
		list := New()
		insertedCount := list.Insert(10, "a")
		assert.Equal(t, 1, insertedCount)
		assert.Equal(t, 1, list.Len)
		assert.NotNil(t, list.Head.Tower[0])
		assert.Equal(t, "a", list.Head.Tower[0].Member)
		assert.Equal(t, 10.0, list.Head.Tower[0].Score)

		assert.Equal(t, 1, list.Head.Span[0])
	})

	t.Run("InsertMultipleSortedAndSpanCheck", func(t *testing.T) {
		list := New()
		list.Insert(10, "a")
		list.Insert(5, "b")
		list.Insert(15, "c")
		assert.Equal(t, 3, list.Len)

		cur := list.Head.Tower[0]
		assert.Equal(t, "b", cur.Member)
		assert.Equal(t, 1, list.Head.Span[0])
		cur = cur.Tower[0]
		assert.Equal(t, "a", cur.Member)
		assert.Equal(t, 1, list.Head.Tower[0].Span[0])
		cur = cur.Tower[0]
		assert.Equal(t, "c", cur.Member)
		assert.Equal(t, 1, list.Head.Tower[0].Tower[0].Span[0])
		assert.Nil(t, cur.Tower[0])

		_, _, rank := list.Search(15, "c")
		assert.Equal(t, 2, rank[0])
	})

	t.Run("DeleteExistingAndSpanCheck", func(t *testing.T) {
		list := New()
		list.Insert(10, "a")
		list.Insert(5, "b")
		list.Insert(15, "c")

		deletedCount := list.Delete(10, "a")
		assert.Equal(t, 1, deletedCount)
		assert.Equal(t, 2, list.Len)

		cur := list.Head.Tower[0]
		assert.Equal(t, "b", cur.Member)
		assert.Equal(t, 1, list.Head.Span[0])
		cur = cur.Tower[0]
		assert.Equal(t, "c", cur.Member)
		assert.Equal(t, 1, list.Head.Tower[0].Span[0])
		assert.Nil(t, cur.Tower[0])

		_, _, rank := list.Search(15, "c")
		assert.Equal(t, 1, rank[0])
	})

	t.Run("DeleteNonExisting", func(t *testing.T) {
		list := New()
		list.Insert(10, "a")
		deletedCount := list.Delete(5, "b")
		assert.Equal(t, 0, deletedCount)
		assert.Equal(t, 1, list.Len)
	})

	t.Run("InsertDuplicate", func(t *testing.T) {
		list := New()
		list.Insert(10, "a")
		insertedCount := list.Insert(10, "a")
		assert.Equal(t, 0, insertedCount)
		assert.Equal(t, 1, list.Len)
		node := list.Head.Tower[0]
		assert.Equal(t, "a", node.Member)
		assert.Equal(t, 10.0, node.Score)
	})

	t.Run("HeightIncreases", func(t *testing.T) {
		list := New()
		for i := range 100 {
			list.Insert(float64(i), fmt.Sprintf("member %d", i))
		}
		assert.True(t, list.Height > 1)
		assert.True(t, list.Height < MAX_HEIGHT)
	})
}

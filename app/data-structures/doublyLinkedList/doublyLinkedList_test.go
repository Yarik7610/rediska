package doublylinkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListStorageDoubleLinkedList(t *testing.T) {
	t.Run("InsertInTheStart", func(t *testing.T) {
		list := &List{}
		InsertInTheStart(list, &Node{Val: "a"})
		assert.Equal(t, 1, list.Len)
		assert.Equal(t, "a", list.Head.Val)
		assert.Equal(t, "a", list.Tail.Val)

		InsertInTheStart(list, &Node{Val: "b"})
		assert.Equal(t, 2, list.Len)
		assert.Equal(t, "b", list.Head.Val)
		assert.Equal(t, "a", list.Tail.Val)
		assert.Equal(t, list.Head.Next, list.Tail)
		assert.Equal(t, list.Tail.Prev, list.Head)
	})

	t.Run("InsertInTheEnd", func(t *testing.T) {
		list := &List{}
		InsertInTheEnd(list, &Node{Val: "a"})
		assert.Equal(t, 1, list.Len)
		assert.Equal(t, "a", list.Head.Val)
		assert.Equal(t, "a", list.Tail.Val)

		InsertInTheEnd(list, &Node{Val: "b"})
		assert.Equal(t, 2, list.Len)
		assert.Equal(t, "a", list.Head.Val)
		assert.Equal(t, "b", list.Tail.Val)
		assert.Equal(t, list.Head.Next, list.Tail)
		assert.Equal(t, list.Tail.Prev, list.Head)
	})

	t.Run("DeleteFromStart", func(t *testing.T) {
		list := &List{}
		assert.Nil(t, DeleteFromStart(list))

		InsertInTheStart(list, &Node{Val: "a"})
		deleted := DeleteFromStart(list)
		assert.Equal(t, "a", deleted.Val)
		assert.Equal(t, 0, list.Len)
		assert.Nil(t, list.Head)
		assert.Nil(t, list.Tail)

		InsertInTheStart(list, &Node{Val: "a"})
		InsertInTheStart(list, &Node{Val: "b"})
		deleted = DeleteFromStart(list)
		assert.Equal(t, "b", deleted.Val)
		assert.Equal(t, 1, list.Len)
		assert.Equal(t, "a", list.Head.Val)
		assert.Equal(t, list.Head, list.Tail)
	})

	t.Run("DeleteFromEnd", func(t *testing.T) {
		list := &List{}
		assert.Nil(t, DeleteFromEnd(list))

		InsertInTheEnd(list, &Node{Val: "a"})
		deleted := DeleteFromEnd(list)
		assert.Equal(t, "a", deleted.Val)
		assert.Equal(t, 0, list.Len)
		assert.Nil(t, list.Head)
		assert.Nil(t, list.Tail)

		InsertInTheEnd(list, &Node{Val: "a"})
		InsertInTheEnd(list, &Node{Val: "b"})
		deleted = DeleteFromEnd(list)
		assert.Equal(t, "b", deleted.Val)
		assert.Equal(t, 1, list.Len)
		assert.Equal(t, "a", list.Tail.Val)
		assert.Equal(t, list.Head, list.Tail)
	})
}

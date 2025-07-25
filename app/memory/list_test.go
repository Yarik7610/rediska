package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleLinkedList(t *testing.T) {
	t.Run("insertInTheStart", func(t *testing.T) {
		list := DoubleLinkedList{}
		list.insertInTheStart(&Node{val: "a"})
		assert.Equal(t, 1, list.len)
		assert.Equal(t, "a", list.head.val)
		assert.Equal(t, "a", list.tail.val)

		list.insertInTheStart(&Node{val: "b"})
		assert.Equal(t, 2, list.len)
		assert.Equal(t, "b", list.head.val)
		assert.Equal(t, "a", list.tail.val)
		assert.Equal(t, list.head.next, list.tail)
		assert.Equal(t, list.tail.prev, list.head)
	})

	t.Run("insertInTheEnd", func(t *testing.T) {
		list := DoubleLinkedList{}
		list.insertInTheEnd(&Node{val: "a"})
		assert.Equal(t, 1, list.len)
		assert.Equal(t, "a", list.head.val)
		assert.Equal(t, "a", list.tail.val)

		list.insertInTheEnd(&Node{val: "b"})
		assert.Equal(t, 2, list.len)
		assert.Equal(t, "a", list.head.val)
		assert.Equal(t, "b", list.tail.val)
		assert.Equal(t, list.head.next, list.tail)
		assert.Equal(t, list.tail.prev, list.head)
	})

	t.Run("deleteFromStart", func(t *testing.T) {
		list := DoubleLinkedList{}
		assert.Nil(t, list.deleteFromStart())

		list.insertInTheStart(&Node{val: "a"})
		deleted := list.deleteFromStart()
		assert.Equal(t, "a", deleted.val)
		assert.Equal(t, 0, list.len)
		assert.Nil(t, list.head)
		assert.Nil(t, list.tail)

		list.insertInTheStart(&Node{val: "a"})
		list.insertInTheStart(&Node{val: "b"})
		deleted = list.deleteFromStart()
		assert.Equal(t, "b", deleted.val)
		assert.Equal(t, 1, list.len)
		assert.Equal(t, "a", list.head.val)
		assert.Equal(t, list.head, list.tail)
	})

	t.Run("deleteFromEnd", func(t *testing.T) {
		list := DoubleLinkedList{}
		assert.Nil(t, list.deleteFromEnd())

		list.insertInTheEnd(&Node{val: "a"})
		deleted := list.deleteFromEnd()
		assert.Equal(t, "a", deleted.val)
		assert.Equal(t, 0, list.len)
		assert.Nil(t, list.head)
		assert.Nil(t, list.tail)

		list.insertInTheEnd(&Node{val: "a"})
		list.insertInTheEnd(&Node{val: "b"})
		deleted = list.deleteFromEnd()
		assert.Equal(t, "b", deleted.val)
		assert.Equal(t, 1, list.len)
		assert.Equal(t, "a", list.tail.val)
		assert.Equal(t, list.head, list.tail)
	})
}

func TestLpop(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("empty list", func(t *testing.T) {
		assert.Nil(t, ls.Lpop("nonexistent", 1))
	})

	t.Run("pop single element", func(t *testing.T) {
		ls.Lpush("list1", "a")
		popped := ls.Lpop("list1", 1)
		assert.Equal(t, []string{"a"}, popped)
		assert.Equal(t, 0, ls.data["list1"].len)
	})

	t.Run("pop more than available", func(t *testing.T) {
		ls.Lpush("list2", "a", "b")
		popped := ls.Lpop("list2", 5)
		assert.Equal(t, []string{"b", "a"}, popped)
		assert.Equal(t, 0, ls.data["list2"].len)
	})

	t.Run("pop with count=0", func(t *testing.T) {
		ls.Lpush("list3", "a")
		popped := ls.Lpop("list3", 0)
		assert.Empty(t, popped)
		assert.Equal(t, 1, ls.data["list3"].len)
	})
}

func TestRpop(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("empty list", func(t *testing.T) {
		assert.Nil(t, ls.Rpop("nonexistent", 1))
	})

	t.Run("pop single element", func(t *testing.T) {
		ls.Rpush("list1", "a")
		popped := ls.Rpop("list1", 1)
		assert.Equal(t, []string{"a"}, popped)
		assert.Equal(t, 0, ls.data["list1"].len)
	})

	t.Run("pop more than available", func(t *testing.T) {
		ls.Rpush("list2", "a", "b")
		popped := ls.Rpop("list2", 5)
		assert.Equal(t, []string{"b", "a"}, popped)
		assert.Equal(t, 0, ls.data["list2"].len)
	})

	t.Run("pop with count=0", func(t *testing.T) {
		ls.Rpush("list3", "a")
		popped := ls.Rpop("list3", 0)
		assert.Empty(t, popped)
		assert.Equal(t, 1, ls.data["list3"].len)
	})
}

func TestGetKeys(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("no keys", func(t *testing.T) {
		assert.Empty(t, ls.GetKeys())
	})

	t.Run("with keys", func(t *testing.T) {
		ls.Lpush("key1", "a")
		ls.Lpush("key2", "b")
		keys := ls.GetKeys()
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
	})
}

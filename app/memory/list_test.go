package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListStorageDoubleLinkedList(t *testing.T) {
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

func TestListStorageLpop(t *testing.T) {
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

	t.Run("concurrent pops", func(t *testing.T) {
		const workers = 10
		var wg sync.WaitGroup
		key := "concurrent_list"

		ls.Lpush(key, "a", "b", "c", "d", "e")

		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				popped := ls.Lpop(key, 1)
				if len(popped) > 0 {
					assert.Contains(t, []string{"a", "b", "c", "d", "e"}, popped[0])
				}
			}()
		}
		wg.Wait()

		assert.True(t, ls.data[key].len >= 0)
	})
}

func TestListStorageRpop(t *testing.T) {
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

	t.Run("concurrent pops", func(t *testing.T) {
		const workers = 10
		var wg sync.WaitGroup
		key := "concurrent_rlist"

		ls.Rpush(key, "a", "b", "c", "d", "e")

		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				popped := ls.Rpop(key, 1)
				if len(popped) > 0 {
					assert.Contains(t, []string{"a", "b", "c", "d", "e"}, popped[0])
				}
			}()
		}
		wg.Wait()

		assert.True(t, ls.data[key].len >= 0)
	})
}

func TestListStorageGetKeys(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("no keys", func(t *testing.T) {
		assert.Empty(t, ls.GetKeys())
	})

	t.Run("concurrent get keys", func(t *testing.T) {
		const workers = 5
		var wg sync.WaitGroup

		wg.Add(workers)
		for i := range workers {
			go func(idx int) {
				defer wg.Done()
				key := fmt.Sprintf("concurrent_key%d", idx)
				ls.Lpush(key, "val")
				keys := ls.GetKeys()
				assert.Contains(t, keys, key)
			}(i)
		}
		wg.Wait()
	})
}

func TestListStorageLlen(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("empty list", func(t *testing.T) {
		assert.Equal(t, 0, ls.Llen("list"))
	})

	t.Run("concurrent list len", func(t *testing.T) {
		const workers = 5
		var wg sync.WaitGroup

		wg.Add(workers)
		for i := range workers {
			go func(idx int) {
				defer wg.Done()
				ls.Lpush("list", "val")
			}(i)
		}
		wg.Wait()
		assert.Equal(t, workers, ls.Llen("list"))
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

	t.Run("concurrent pops", func(t *testing.T) {
		const workers = 10
		var wg sync.WaitGroup
		key := "concurrent_rlist"

		ls.Rpush(key, "a", "b", "c", "d", "e")

		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				popped := ls.Rpop(key, 1)
				if len(popped) > 0 {
					assert.Contains(t, []string{"a", "b", "c", "d", "e"}, popped[0])
				}
			}()
		}
		wg.Wait()

		assert.True(t, ls.data[key].len >= 0)
	})
}

func TestListStorageGetKey(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("no key", func(t *testing.T) {
		list, ok := ls.Get("key1")
		assert.Empty(t, list)
		assert.False(t, ok)
	})

	t.Run("concurrent get", func(t *testing.T) {
		const workers = 10
		var wg sync.WaitGroup
		key := "concurrent_get_key"
		ls.Lpush(key, "a", "b", "c")

		wg.Add(workers)
		for range workers {
			go func() {
				defer wg.Done()
				list, ok := ls.Get(key)
				assert.True(t, ok)
				assert.Equal(t, 3, list.len)
			}()
		}
		wg.Wait()
	})
}

func TestListStorageDelKey(t *testing.T) {
	ls := NewListStorage()
	ls.data = make(map[string]*DoubleLinkedList)

	t.Run("no key", func(t *testing.T) {
		assert.Empty(t, ls.data)
		ls.Del("key1")
		assert.Empty(t, ls.data)
	})

	t.Run("concurrent delete", func(t *testing.T) {
		const workers = 5
		var wg sync.WaitGroup

		wg.Add(workers)
		for i := range workers {
			go func(idx int) {
				defer wg.Done()
				key := fmt.Sprintf("del_key%d", idx)
				ls.Lpush(key, "val")
				ls.Del(key)
				list, ok := ls.Get(key)
				assert.Empty(t, list)
				assert.False(t, ok)
			}(i)
		}
		wg.Wait()
	})
}

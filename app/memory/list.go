package memory

import (
	"sync"
	"time"
)

type Node struct {
	val  string
	next *Node
	prev *Node
}

type DoubleLinkedList struct {
	head *Node
	tail *Node
	len  int
}

type ListStorage struct {
	data  map[string]*DoubleLinkedList
	rwMut sync.RWMutex
	cond  *sync.Cond
}

func NewListStorage() *ListStorage {
	ls := &ListStorage{data: make(map[string]*DoubleLinkedList)}
	ls.cond = sync.NewCond(&ls.rwMut)
	return ls
}

func (ls *ListStorage) Llen(key string) int {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	list, ok := ls.data[key]
	if !ok {
		return 0
	}
	return list.len
}

func (ls *ListStorage) Lrange(key string, startIdx, stopIdx int) []string {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	values := make([]string, 0)

	list, ok := ls.data[key]
	if !ok {
		return values
	}

	len := list.len

	if startIdx < 0 {
		startIdx += len
		if startIdx < 0 {
			startIdx = 0
		}
	}
	if stopIdx < 0 {
		stopIdx += len
		if stopIdx >= len {
			stopIdx = len - 1
		}
	}

	if stopIdx < startIdx {
		return values
	}

	if startIdx >= len {
		return values
	}
	if stopIdx >= len {
		stopIdx = len - 1
	}

	curNode := list.head

	for range startIdx {
		curNode = curNode.next
	}
	for range stopIdx - startIdx + 1 {
		if curNode == nil {
			break
		}
		values = append(values, curNode.val)
		curNode = curNode.next
	}
	return values
}

func (ls *ListStorage) GetKeys() []string {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	keys := make([]string, 0)
	for key := range ls.data {
		keys = append(keys, key)
	}

	return keys
}

func (ls *ListStorage) Get(key string) (*DoubleLinkedList, bool) {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	list, ok := ls.data[key]
	return list, ok
}

func (ls *ListStorage) Del(key string) {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	delete(ls.data, key)
}

func (ls *ListStorage) Rpop(key string, count int) []string {
	return ls.pop(key, count, deleteFromEnd)
}

func (ls *ListStorage) Lpop(key string, count int) []string {
	return ls.pop(key, count, deleteFromStart)
}

func (ls *ListStorage) Brpop(key string, timeoutS float64) *string {
	return ls.bpop(key, timeoutS, deleteFromEnd)
}

func (ls *ListStorage) Blpop(key string, timeoutS float64) *string {
	return ls.bpop(key, timeoutS, deleteFromStart)
}

func (ls *ListStorage) Lpush(key string, values ...string) int {
	return ls.push(insertInTheStart, key, values...)
}

func (ls *ListStorage) Rpush(key string, values ...string) int {
	return ls.push(insertInTheEnd, key, values...)
}

func (ls *ListStorage) pop(key string, count int, popFn func(list *DoubleLinkedList) *Node) []string {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	list, ok := ls.data[key]
	if !ok {
		return nil
	}

	if count > list.len {
		count = list.len
	}

	popped := make([]string, 0)
	for range count {
		deleted := popFn(list)
		if deleted == nil {
			break
		}
		popped = append(popped, deleted.val)
	}

	return popped
}

func (ls *ListStorage) bpop(key string, timeoutS float64, popFn func(list *DoubleLinkedList) *Node) *string {
	ls.rwMut.Lock()

	if list, ok := ls.data[key]; ok && list.len > 0 {
		popped := popFn(list)
		ls.rwMut.Unlock()
		return &popped.val
	}

	if timeoutS < 0 {
		ls.rwMut.Unlock()
		return nil
	}

	if timeoutS == 0 {
		for {
			ls.cond.Wait()
			if list, ok := ls.data[key]; ok && list.len > 0 {
				popped := popFn(list)
				ls.rwMut.Unlock()
				return &popped.val
			}
		}
	}

	timer := time.After(time.Duration(timeoutS * float64(time.Second)))
	for {
		ls.rwMut.Unlock()

		select {
		case <-timer:
			return nil
		default:
			time.Sleep(25 * time.Millisecond)
		}

		ls.rwMut.Lock()
		if list, ok := ls.data[key]; ok && list.len > 0 {
			popped := popFn(list)
			ls.rwMut.Unlock()
			return &popped.val
		}
	}
}

func (ls *ListStorage) push(pushFn func(list *DoubleLinkedList, n *Node), key string, values ...string) int {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	if _, ok := ls.data[key]; !ok {
		ls.data[key] = &DoubleLinkedList{}
	}

	list := ls.data[key]

	for _, val := range values {
		n := &Node{val: val}
		pushFn(list, n)
	}
	ls.cond.Signal()
	return list.len
}

func insertInTheStart(list *DoubleLinkedList, n *Node) {
	if list.head == nil {
		list.head = n
		list.tail = n
	} else {
		n.next = list.head
		list.head.prev = n
		list.head = n
	}
	list.len++
}

func insertInTheEnd(list *DoubleLinkedList, n *Node) {
	if list.tail == nil {
		list.tail = n
		list.head = n
	} else {
		list.tail.next = n
		n.prev = list.tail
		list.tail = n
	}
	list.len++
}

func deleteFromStart(list *DoubleLinkedList) *Node {
	if list.head == nil {
		return nil
	}
	deleted := list.head
	next := list.head.next
	if next == nil {
		list.head = nil
		list.tail = nil
	} else {
		list.head.next = nil
		next.prev = nil
		list.head = next
	}
	list.len--
	return deleted
}

func deleteFromEnd(list *DoubleLinkedList) *Node {
	if list.tail == nil {
		return nil
	}
	deleted := list.tail
	prev := list.tail.prev
	if prev == nil {
		list.head = nil
		list.tail = nil
	} else {
		list.tail.prev = nil
		prev.next = nil
		list.tail = prev
	}
	list.len--
	return deleted
}

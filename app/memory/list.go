package memory

import "sync"

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
}

func NewListStorage() *ListStorage {
	return &ListStorage{}
}

func (ls *ListStorage) Lpush(key string, values ...string) int {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	if _, ok := ls.data[key]; !ok {
		ls.data[key] = &DoubleLinkedList{}
	}

	list := ls.data[key]

	for _, val := range values {
		n := &Node{val: val}
		list.insertInTheStart(n)
	}
	return list.len
}

func (ls *ListStorage) Rpush(key string, values ...string) int {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	if _, ok := ls.data[key]; !ok {
		ls.data[key] = &DoubleLinkedList{}
	}

	list := ls.data[key]

	for _, val := range values {
		n := &Node{val: val}
		list.insertInTheEnd(n)
	}
	return list.len
}

func (ls *ListStorage) Rpop(key string, count int) []string {
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
		deleted := list.deleteFromEnd()
		if deleted == nil {
			break
		}
		popped = append(popped, deleted.val)
	}

	return popped
}

func (ls *ListStorage) Lpop(key string, count int) []string {
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
		deleted := list.deleteFromStart()
		if deleted == nil {
			break
		}
		popped = append(popped, deleted.val)
	}

	return popped
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

func (list *DoubleLinkedList) insertInTheStart(n *Node) {
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

func (list *DoubleLinkedList) insertInTheEnd(n *Node) {
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

func (list *DoubleLinkedList) deleteFromStart() *Node {
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

func (list *DoubleLinkedList) deleteFromEnd() *Node {
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

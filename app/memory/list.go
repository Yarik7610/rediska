package memory

import (
	"sync"
	"time"

	doublylinkedlist "github.com/codecrafters-io/redis-starter-go/app/data-structures/doubly-linked-list"
)

type ListStorage interface {
	baseStorage
	Llen(key string) int
	Lrange(key string, startIdx, stopIdx int) []string
	Rpop(key string, count int) []string
	Lpop(key string, count int) []string
	Brpop(key string, timeoutS float64) *string
	Blpop(key string, timeoutS float64) *string
	Lpush(key string, values ...string) int
	Rpush(key string, values ...string) int
}

type listStorage struct {
	data  map[string]*doublylinkedlist.List
	rwMut sync.RWMutex
	cond  *sync.Cond
}

func NewListStorage() ListStorage {
	ls := &listStorage{data: make(map[string]*doublylinkedlist.List)}
	ls.cond = sync.NewCond(&ls.rwMut)
	return ls
}

func (ls *listStorage) Keys() []string {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	keys := make([]string, 0)
	for key := range ls.data {
		keys = append(keys, key)
	}

	return keys
}

func (ls *listStorage) Has(key string) bool {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()
	_, ok := ls.data[key]
	return ok
}

func (ls *listStorage) Del(key string) {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()
	delete(ls.data, key)
}

func (ls *listStorage) Llen(key string) int {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	list, ok := ls.data[key]
	if !ok {
		return 0
	}
	return list.Len
}

func (ls *listStorage) Lrange(key string, startIdx, stopIdx int) []string {
	ls.rwMut.RLock()
	defer ls.rwMut.RUnlock()

	values := make([]string, 0)

	list, ok := ls.data[key]
	if !ok {
		return values
	}

	var err error
	startIdx, stopIdx, err = handleRangeIndexes(startIdx, stopIdx, list.Len)
	if err != nil {
		return values
	}

	cur := list.Head
	for range startIdx {
		cur = cur.Next
	}
	for range stopIdx - startIdx + 1 {
		if cur == nil {
			break
		}
		values = append(values, cur.Val)
		cur = cur.Next
	}
	return values
}

func (ls *listStorage) Rpop(key string, count int) []string {
	return ls.pop(key, count, doublylinkedlist.DeleteFromEnd)
}

func (ls *listStorage) Lpop(key string, count int) []string {
	return ls.pop(key, count, doublylinkedlist.DeleteFromStart)
}

func (ls *listStorage) Brpop(key string, timeoutS float64) *string {
	return ls.bpop(key, timeoutS, doublylinkedlist.DeleteFromEnd)
}

func (ls *listStorage) Blpop(key string, timeoutS float64) *string {
	return ls.bpop(key, timeoutS, doublylinkedlist.DeleteFromStart)
}

func (ls *listStorage) Lpush(key string, values ...string) int {
	return ls.push(doublylinkedlist.InsertInTheStart, key, values...)
}

func (ls *listStorage) Rpush(key string, values ...string) int {
	return ls.push(doublylinkedlist.InsertInTheEnd, key, values...)
}

func (ls *listStorage) pop(key string, count int, popFn func(list *doublylinkedlist.List) *doublylinkedlist.Node) []string {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	list, ok := ls.data[key]
	if !ok {
		return nil
	}

	if count > list.Len {
		count = list.Len
	}

	popped := make([]string, 0)
	for range count {
		deleted := popFn(list)
		if deleted == nil {
			break
		}
		popped = append(popped, deleted.Val)
	}

	return popped
}

func (ls *listStorage) bpop(key string, timeoutS float64, popFn func(list *doublylinkedlist.List) *doublylinkedlist.Node) *string {
	ls.rwMut.Lock()

	if list, ok := ls.data[key]; ok && list.Len > 0 {
		popped := popFn(list)
		ls.rwMut.Unlock()
		return &popped.Val
	}

	if timeoutS < 0 {
		ls.rwMut.Unlock()
		return nil
	}

	if timeoutS == 0 {
		for {
			ls.cond.Wait()
			if list, ok := ls.data[key]; ok && list.Len > 0 {
				popped := popFn(list)
				ls.rwMut.Unlock()
				return &popped.Val
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
		if list, ok := ls.data[key]; ok && list.Len > 0 {
			popped := popFn(list)
			ls.rwMut.Unlock()
			return &popped.Val
		}
	}
}

func (ls *listStorage) push(pushFn func(list *doublylinkedlist.List, n *doublylinkedlist.Node), key string, values ...string) int {
	ls.rwMut.Lock()
	defer ls.rwMut.Unlock()

	if _, ok := ls.data[key]; !ok {
		ls.data[key] = &doublylinkedlist.List{}
	}

	list := ls.data[key]

	for _, val := range values {
		n := &doublylinkedlist.Node{Val: val}
		pushFn(list, n)
	}
	ls.cond.Signal()
	return list.Len
}

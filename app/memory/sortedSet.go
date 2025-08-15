package memory

import (
	"sync"

	skiplist "github.com/codecrafters-io/redis-starter-go/app/data-structures/skip-list"
)

type sortedSet struct {
	dict     map[string]float64
	skipList *skiplist.List
}

type SortedSetStorage interface {
	baseStorage
	Zadd(key string, scores []float64, members []string) int
	Zrem(key string, members []string) int
	Zrank(key string, member string) int
	Zrange(key string, startIdx, stopIdx int) []string
	Zcard(key string) int
	Zscore(key string, member string) *float64
}

type sortedSetStorage struct {
	data  map[string]*sortedSet
	rwMut sync.RWMutex
}

func NewSortedSetStorage() SortedSetStorage {
	return &sortedSetStorage{data: make(map[string]*sortedSet)}
}

func (s *sortedSetStorage) Zadd(key string, scores []float64, members []string) int {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	if _, ok := s.data[key]; !ok {
		s.data[key] = &sortedSet{
			dict:     make(map[string]float64),
			skipList: skiplist.New(),
		}
	}

	sortedSet := s.data[key]

	insertedCount := 0
	for i, member := range members {
		if oldScore, ok := sortedSet.dict[member]; ok {
			sortedSet.skipList.Delete(oldScore, member)
			sortedSet.skipList.Insert(scores[i], member)
		} else {
			insertedCount += sortedSet.skipList.Insert(scores[i], member)
		}
		sortedSet.dict[member] = scores[i]
	}

	return insertedCount
}

func (s *sortedSetStorage) Zrem(key string, members []string) int {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	sortedSet, ok := s.data[key]
	if !ok {
		return 0
	}

	deletedCount := 0
	for _, member := range members {
		score, ok := sortedSet.dict[member]
		if !ok {
			continue
		}
		deletedCount += sortedSet.skipList.Delete(score, member)
		delete(sortedSet.dict, member)
	}
	return deletedCount
}

func (s *sortedSetStorage) Zrank(key string, member string) int {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()

	sortedSet, ok := s.data[key]
	if !ok {
		return -1
	}

	score, ok := sortedSet.dict[member]
	if !ok {
		return -1
	}

	_, _, rank := sortedSet.skipList.Search(score, member)
	return rank[0]
}

func (s *sortedSetStorage) Zrange(key string, startIdx, stopIdx int) []string {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()

	values := make([]string, 0)

	sortedSet, ok := s.data[key]
	if !ok {
		return values
	}

	var err error
	startIdx, stopIdx, err = handleRangeIndexes(startIdx, stopIdx, sortedSet.skipList.Len)
	if err != nil {
		return values
	}

	cur := sortedSet.skipList.Head

	for range startIdx + 1 {
		cur = cur.Tower[0]
	}
	for range stopIdx - startIdx + 1 {
		if cur == nil {
			break
		}
		values = append(values, cur.Member)
		cur = cur.Tower[0]
	}
	return values
}

func (s *sortedSetStorage) Zcard(key string) int {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()

	sortedSet, ok := s.data[key]
	if !ok {
		return 0
	}

	return sortedSet.skipList.Len
}

func (s *sortedSetStorage) Zscore(key string, member string) *float64 {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()

	sortedSet, ok := s.data[key]
	if !ok {
		return nil
	}

	score, ok := sortedSet.dict[member]
	if !ok {
		return nil
	}
	return &score
}

func (s *sortedSetStorage) Keys() []string {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()

	keys := make([]string, 0)
	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

func (s *sortedSetStorage) Has(key string) bool {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()
	_, ok := s.data[key]
	return ok
}

func (s *sortedSetStorage) Del(key string) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()
	delete(s.data, key)
}

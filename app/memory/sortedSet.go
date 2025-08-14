package memory

import (
	"sync"

	skiplist "github.com/codecrafters-io/redis-starter-go/app/data-structures/skipList"
)

type sortedSet struct {
	dict     map[string]float64
	skipList *skiplist.List
}

type SortedSetStorage interface {
	baseStorage
	Zadd(key string, scores []float64, members []string) int
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
		sortedSet.dict[member] = scores[i]
		insertedCount = sortedSet.skipList.Insert(scores[i], member)
	}

	return insertedCount
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

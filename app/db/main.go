package db

import (
	"sync"
	"time"
)

type Item struct {
	Value   string
	Expires time.Time
}

type Storage struct {
	data  map[string]Item
	rwMut sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]Item, 0),
	}
}

func (s *Storage) Get(key string) (Item, bool) {
	s.rwMut.RLock()
	defer s.rwMut.RUnlock()
	item, ok := s.data[key]
	return item, ok
}

func (s *Storage) Set(key, value string) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()
	s.data[key] = Item{Value: value}
}

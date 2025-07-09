package memory

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

func (s *Storage) Get(key string) (*Item, bool) {
	s.rwMut.RLock()
	item, ok := s.data[key]
	if !ok {
		s.rwMut.RUnlock()
		return nil, false
	}

	if !item.Expires.IsZero() && item.Expires.Before(time.Now()) {
		s.rwMut.RUnlock()
		s.rwMut.Lock()
		defer s.rwMut.Unlock()

		// Repeat checking because of small non-blocking window between RUnlock() and Lock()
		item, ok = s.data[key]
		if !ok {
			return nil, false
		}
		if !item.Expires.IsZero() && item.Expires.Before(time.Now()) {
			delete(s.data, key)
			return nil, false
		}

		return &item, true
	}

	s.rwMut.RUnlock()
	return &item, ok
}

func (s *Storage) Set(key, value string) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()
	s.data[key] = Item{Value: value}
}

func (s *Storage) SetWithExpiry(key, value string, expiry time.Duration) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	if expiry <= 0 {
		delete(s.data, key)
		return
	}

	s.data[key] = Item{Value: value, Expires: time.Now().Add(expiry)}
}

func (s *Storage) CleanExpiredKeys() {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	for key, item := range s.data {
		if !item.Expires.IsZero() && item.Expires.Before(time.Now()) {
			delete(s.data, key)
		}
	}
}

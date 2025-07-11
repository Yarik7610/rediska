package memory

import (
	"sync"
	"time"
)

type Item struct {
	Value   string
	Expires time.Time
}

// No division on 2 maps: expired and unexpired
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

	if ItemExpired(&item) {
		s.rwMut.RUnlock()
		s.rwMut.Lock()
		defer s.rwMut.Unlock()

		// Repeat checking because of small non-blocking window between RUnlock() and Lock()
		item, ok = s.data[key]
		if !ok {
			return nil, false
		}
		if ItemExpired(&item) {
			delete(s.data, key)
			return nil, false
		}

		return &item, true
	}

	s.rwMut.RUnlock()
	return &item, ok
}

func (s *Storage) GetKeys() []string {
	s.rwMut.RLock()
	var keys []string
	var expiredKeys []string

	for key, item := range s.data {
		if ItemExpired(&item) {
			expiredKeys = append(expiredKeys, key)
		} else {
			keys = append(keys, key)
		}
	}
	s.rwMut.RUnlock()

	if len(expiredKeys) > 0 {
		s.rwMut.Lock()
		for _, key := range expiredKeys {
			// Repeat checking because of small non-blocking window between RUnlock() and Lock()
			if item, ok := s.data[key]; ok && ItemExpired(&item) {
				delete(s.data, key)
			}
		}
		s.rwMut.Unlock()
	}

	return keys
}

func (s *Storage) Set(key, value string) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()
	s.data[key] = Item{Value: value}
}

func (s *Storage) SetWithExpiry(key, value string, expires time.Time) {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	if !expires.After(time.Now()) || expires.IsZero() {
		delete(s.data, key)
		return
	}

	s.data[key] = Item{Value: value, Expires: expires}
}

func (s *Storage) CleanExpiredKeys() {
	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	for key, item := range s.data {
		if ItemExpired(&item) {
			delete(s.data, key)
		}
	}
}

func ItemExpired(item *Item) bool {
	return ItemHasExpiration(item) && item.Expires.Before(time.Now())
}

func ItemHasExpiration(item *Item) bool {
	return !item.Expires.IsZero()
}

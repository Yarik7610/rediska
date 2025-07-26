package memory

import (
	"sync"
	"time"
)

type String struct {
	Value   string
	Expires time.Time
}

type StringStorage struct {
	data  map[string]String
	rwMut sync.RWMutex
}

func NewStringStorage() *StringStorage {
	return &StringStorage{
		data: make(map[string]String, 0),
	}
}

func (ss *StringStorage) Get(key string) (*String, bool) {
	ss.rwMut.RLock()
	item, ok := ss.data[key]
	if !ok {
		ss.rwMut.RUnlock()
		return nil, false
	}

	if ss.ItemExpired(&item) {
		ss.rwMut.RUnlock()
		ss.rwMut.Lock()
		defer ss.rwMut.Unlock()

		// Repeat checking because of small non-blocking window between RUnlock() and Lock()
		item, ok = ss.data[key]
		if !ok {
			return nil, false
		}
		if ss.ItemExpired(&item) {
			delete(ss.data, key)
			return nil, false
		}

		return &item, true
	}

	ss.rwMut.RUnlock()
	return &item, ok
}

func (ss *StringStorage) GetKeys() []string {
	ss.rwMut.RLock()
	var keys []string
	var expiredKeys []string

	for key, item := range ss.data {
		if ss.ItemExpired(&item) {
			expiredKeys = append(expiredKeys, key)
		} else {
			keys = append(keys, key)
		}
	}
	ss.rwMut.RUnlock()

	if len(expiredKeys) > 0 {
		ss.rwMut.Lock()
		for _, key := range expiredKeys {
			// Repeat checking because of small non-blocking window between RUnlock() and Lock()
			if item, ok := ss.data[key]; ok && ss.ItemExpired(&item) {
				delete(ss.data, key)
			}
		}
		ss.rwMut.Unlock()
	}

	return keys
}

func (ss *StringStorage) Set(key, value string) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()
	ss.data[key] = String{Value: value}
}

func (ss *StringStorage) SetWithExpiry(key, value string, expires time.Time) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	if !expires.After(time.Now()) || expires.IsZero() {
		delete(ss.data, key)
		return
	}

	ss.data[key] = String{Value: value, Expires: expires}
}

func (ss *StringStorage) Del(key string) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	delete(ss.data, key)
}

func (ss *StringStorage) CleanExpiredKeys() {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	for key, item := range ss.data {
		if ss.ItemExpired(&item) {
			delete(ss.data, key)
		}
	}
}

func (ss *StringStorage) ItemExpired(item *String) bool {
	return ss.ItemHasExpiration(item) && item.Expires.Before(time.Now())
}

func (ss *StringStorage) ItemHasExpiration(item *String) bool {
	return !item.Expires.IsZero()
}

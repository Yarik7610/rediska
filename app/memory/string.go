package memory

import (
	"sync"
	"time"
)

type String struct {
	Value   string
	Expires time.Time
}

type StringStorage interface {
	baseStorage
	Get(key string) (*String, bool)
	Set(key, value string)
	SetWithExpiry(key, value string, expires time.Time)
	CleanExpiredKeys()
	ItemExpired(item *String) bool
	ItemHasExpiration(item *String) bool
}

type stringStorage struct {
	data  map[string]String
	rwMut sync.RWMutex
}

var _ StringStorage = (*stringStorage)(nil)

func NewStringStorage() *stringStorage {
	return &stringStorage{
		data: make(map[string]String, 0),
	}
}

func (ss *stringStorage) Keys() []string {
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

func (ss *stringStorage) Has(key string) bool {
	ss.rwMut.RLock()
	defer ss.rwMut.RUnlock()
	_, ok := ss.data[key]
	return ok
}

func (ss *stringStorage) Del(key string) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()
	delete(ss.data, key)
}

func (ss *stringStorage) Get(key string) (*String, bool) {
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

func (ss *stringStorage) Set(key, value string) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()
	ss.data[key] = String{Value: value}
}

func (ss *stringStorage) SetWithExpiry(key, value string, expires time.Time) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	if !expires.After(time.Now()) || expires.IsZero() {
		delete(ss.data, key)
		return
	}

	ss.data[key] = String{Value: value, Expires: expires}
}

func (ss *stringStorage) CleanExpiredKeys() {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	for key, item := range ss.data {
		if ss.ItemExpired(&item) {
			delete(ss.data, key)
		}
	}
}

func (ss *stringStorage) ItemExpired(item *String) bool {
	return ss.ItemHasExpiration(item) && item.Expires.Before(time.Now())
}

func (ss *stringStorage) ItemHasExpiration(item *String) bool {
	return !item.Expires.IsZero()
}

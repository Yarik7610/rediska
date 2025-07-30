package memory

import (
	"fmt"
	"maps"
	"sync"
)

type Entry map[string]string

type Stream struct {
	data  map[string]Entry
	rwMut sync.RWMutex
}

type StreamStorage interface {
	baseStorage
	Xadd(streamKey string, requestedStreamID string, entryFields map[string]string) (string, error)
}

type streamStorage struct {
	data  map[string]*Stream
	rwMut sync.RWMutex
}

var _ StreamStorage = (*streamStorage)(nil)

func NewStreamStorage() *streamStorage {
	return &streamStorage{data: make(map[string]*Stream)}
}

func (ss *streamStorage) Xadd(streamKey string, requestedStreamID string, entryFields map[string]string) (string, error) {
	stream := ss.getOrCreateStream(streamKey)
	stream.rwMut.Lock()
	defer stream.rwMut.Unlock()

	if _, ok := stream.data[requestedStreamID]; ok {
		return "", fmt.Errorf("entry with such stream ID already exists")
	}

	entry := make(Entry)
	maps.Copy(entry, entryFields)
	stream.data[requestedStreamID] = entry

	return requestedStreamID, nil
}

func (ss *streamStorage) Keys() []string {
	ss.rwMut.RLock()
	defer ss.rwMut.RUnlock()

	keys := make([]string, 0)
	for key := range ss.data {
		keys = append(keys, key)
	}

	return keys
}

func (ss *streamStorage) Has(key string) bool {
	ss.rwMut.RLock()
	defer ss.rwMut.RUnlock()
	_, ok := ss.data[key]
	return ok
}

func (ss *streamStorage) Del(key string) {
	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()
	delete(ss.data, key)
}

func (ss *streamStorage) getOrCreateStream(streamKey string) *Stream {
	ss.rwMut.RLock()
	if stream, ok := ss.data[streamKey]; ok {
		ss.rwMut.RUnlock()
		return stream
	}
	ss.rwMut.RUnlock()

	ss.rwMut.Lock()
	defer ss.rwMut.Unlock()

	// Repeat checking because of small non-blocking window between RUnlock() and Lock()
	if stream, ok := ss.data[streamKey]; ok {
		return stream
	}

	stream := &Stream{data: make(map[string]Entry)}
	ss.data[streamKey] = stream
	return stream
}

func (s *Stream) getOrCreateEntry(streamID string) map[string]string {
	s.rwMut.RLock()
	if entry, ok := s.data[streamID]; ok {
		s.rwMut.RUnlock()
		return entry
	}
	s.rwMut.RUnlock()

	s.rwMut.Lock()
	defer s.rwMut.Unlock()

	// Repeat checking because of small non-blocking window between RUnlock() and Lock()
	if entry, ok := s.data[streamID]; ok {
		return entry
	}

	entry := make(Entry)
	s.data[streamID] = entry
	return entry
}

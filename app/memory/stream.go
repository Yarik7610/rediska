package memory

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Entry map[string]string

type Stream struct {
	data         map[string]Entry
	lastStreamID string
	lastMSTime   int64
	lastSeqNum   int
	rwMut        sync.RWMutex
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
	if len(entryFields) == 0 {
		return "", fmt.Errorf("entry with empty fields isn't allowed")
	}

	stream := ss.getOrCreateStream(streamKey)

	var msTime int64
	var seqNum int
	var streamID string

	stream.rwMut.Lock()
	defer stream.rwMut.Unlock()

	if requestedStreamID == "*" {
		streamID, msTime, seqNum = generateStreamID()
	} else {
		if requestedStreamID == "0-0" {
			return "", fmt.Errorf("The ID specified in XADD must be greater than 0-0")
		}
		splitted := strings.Split(requestedStreamID, "-")
		if len(splitted) != 2 {
			return "", fmt.Errorf("detected wrong stream id format, need <millisecondsTime>-<sequenceNumber, got: %s", streamID)
		}
		rawSeqNum := splitted[1]
		newMSTime, err := strconv.ParseInt(splitted[0], 10, 64)
		if err != nil {
			return "", fmt.Errorf("milliseconds time format int error: %v", err)
		}
		if newMSTime < stream.lastMSTime {
			return "", fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
		}
		msTime = newMSTime

		if rawSeqNum == "*" {
			if msTime == stream.lastMSTime {
				seqNum = stream.lastSeqNum + 1
			} else {
				seqNum = 0
			}
		} else {
			newSeqNum, err := strconv.Atoi(rawSeqNum)
			if err != nil {
				return "", fmt.Errorf("sequence number atoi error: %v", err)
			}
			if msTime == stream.lastMSTime && newSeqNum <= stream.lastSeqNum {
				return "", fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
			}
			seqNum = newSeqNum
		}

		streamID = fmt.Sprintf("%d-%d", newMSTime, seqNum)
	}

	if _, ok := stream.data[streamID]; ok {
		return "", fmt.Errorf("entry with such stream ID already exists")
	}

	entry := make(Entry)
	maps.Copy(entry, entryFields)

	stream.data[streamID] = entry
	stream.lastMSTime = msTime
	stream.lastSeqNum = seqNum
	stream.lastStreamID = streamID

	return streamID, nil
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

func generateStreamID() (streamID string, msTime int64, seqNum int) {
	seqNum = 0
	msTime = time.Now().Local().UnixMilli()
	streamID = fmt.Sprintf("%d-%d", msTime, seqNum)
	return
}

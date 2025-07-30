package memory

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entry map[string]string

type topEntry struct {
	streamID string
	msTime   int64
	seqNum   int
}

type stream struct {
	data     map[string]entry
	topEntry topEntry
	rwMut    sync.RWMutex
}

type StreamStorage interface {
	baseStorage
	Xadd(streamKey string, requestedStreamID string, entryFields map[string]string) (string, error)
}

type streamStorage struct {
	data  map[string]*stream
	rwMut sync.RWMutex
}

var _ StreamStorage = (*streamStorage)(nil)

func NewStreamStorage() *streamStorage {
	return &streamStorage{data: make(map[string]*stream)}
}

func (ss *streamStorage) Xadd(streamKey string, requestedStreamID string, entryFields map[string]string) (string, error) {
	if len(entryFields) == 0 {
		return "", fmt.Errorf("entry with empty fields isn't allowed")
	}

	stream := ss.getOrCreateStream(streamKey)
	stream.rwMut.Lock()
	defer stream.rwMut.Unlock()

	streamID, msTime, seqNum, err := validateAndGenerateStreamID(requestedStreamID, stream)
	if err != nil {
		return "", err
	}

	entry := make(entry)
	maps.Copy(entry, entryFields)

	topEntry := topEntry{streamID: streamID, msTime: msTime, seqNum: seqNum}
	stream.data[streamID] = entry
	stream.topEntry = topEntry

	return streamID, nil
}

func validateAndGenerateStreamID(requestedStreamID string, stream *stream) (string, int64, int, error) {
	if requestedStreamID == "*" {
		streamID, msTime, seqNum := generateStreamID()
		return streamID, msTime, seqNum, nil
	}

	if requestedStreamID == "0-0" {
		return "", 0, 0, fmt.Errorf("The ID specified in XADD must be greater than 0-0")
	}

	splitted := strings.Split(requestedStreamID, "-")
	if len(splitted) != 2 {
		return "", 0, 0, fmt.Errorf("detected wrong stream id format, need <millisecondsTime>-<sequenceNumber, got: %s", requestedStreamID)
	}

	rawMSTime := splitted[0]
	rawSeqNum := splitted[1]

	msTime, err := strconv.ParseInt(rawMSTime, 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("milliseconds time format int error: %v", err)
	}

	if msTime < stream.topEntry.msTime {
		return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	var seqNum int
	if rawSeqNum == "*" {
		if msTime == stream.topEntry.msTime {
			seqNum = stream.topEntry.seqNum + 1
		} else {
			seqNum = 0
		}
	} else {
		seqNum, err = strconv.Atoi(rawSeqNum)
		if err != nil {
			return "", 0, 0, fmt.Errorf("sequence number atoi error: %v", err)
		}
		if msTime == stream.topEntry.msTime && seqNum <= stream.topEntry.seqNum {
			return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

	streamID := fmt.Sprintf("%d-%d", msTime, seqNum)
	return streamID, msTime, seqNum, nil
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

func (ss *streamStorage) getOrCreateStream(streamKey string) *stream {
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

	stream := &stream{data: make(map[string]entry)}
	ss.data[streamKey] = stream
	return stream
}

func generateStreamID() (streamID string, msTime int64, seqNum int) {
	seqNum = 0
	msTime = time.Now().Local().UnixMilli()
	streamID = fmt.Sprintf("%d-%d", msTime, seqNum)
	return
}

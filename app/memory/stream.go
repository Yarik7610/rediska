package memory

import (
	"fmt"
	"maps"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entry map[string]string

type topEntry struct {
	streamID string
	timeMS   int64
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
	Xrange(streamKey string, startID string, endID string) ([]XrangeEntry, error)
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

	streamID, timeMS, seqNum, err := validateAndGenerateStreamID(requestedStreamID, stream)
	if err != nil {
		return "", err
	}

	entry := make(entry)
	maps.Copy(entry, entryFields)

	topEntry := topEntry{streamID: streamID, timeMS: timeMS, seqNum: seqNum}
	stream.data[streamID] = entry
	stream.topEntry = topEntry

	return streamID, nil
}

type XrangeEntry struct {
	StreamID string
	Entry    *entry
}

func (ss *streamStorage) Xrange(streamKey string, startID string, endID string) ([]XrangeEntry, error) {
	stream := ss.getOrCreateStream(streamKey)
	stream.rwMut.RLock()
	defer stream.rwMut.RUnlock()

	startTimeMS, startSeqNum, err := validateAndParseRangeStreamID(startID, true, stream)
	if err != nil {
		return nil, fmt.Errorf("start ID validation failed: %v", err)
	}
	endTimeMS, endSeqNum, err := validateAndParseRangeStreamID(endID, false, stream)
	if err != nil {
		return nil, fmt.Errorf("start ID validation failed: %v", err)
	}

	entries := make([]XrangeEntry, 0)
	timeMS := startTimeMS
	seqNum := startSeqNum

	for timeMS <= endTimeMS {
		if timeMS == endTimeMS && seqNum > endSeqNum {
			break
		}
		streamID := fmt.Sprintf("%d-%d", timeMS, seqNum)
		entry, ok := stream.data[streamID]
		if !ok {
			timeMS += 1
			seqNum = 0
			continue
		}
		entries = append(entries, XrangeEntry{StreamID: streamID, Entry: &entry})
		seqNum += 1
	}

	return entries, nil
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

func generateStreamID() (streamID string, timeMS int64, seqNum int) {
	seqNum = 0
	timeMS = time.Now().Local().UnixMilli()
	streamID = fmt.Sprintf("%d-%d", timeMS, seqNum)
	return
}

func validateAndGenerateStreamID(requestedStreamID string, stream *stream) (string, int64, int, error) {
	if requestedStreamID == "*" {
		streamID, timeMS, seqNum := generateStreamID()
		return streamID, timeMS, seqNum, nil
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

	timeMS, err := strconv.ParseInt(rawMSTime, 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("milliseconds time parse int error: %v", err)
	}

	if timeMS < stream.topEntry.timeMS {
		return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	var seqNum int
	if rawSeqNum == "*" {
		if timeMS == stream.topEntry.timeMS {
			seqNum = stream.topEntry.seqNum + 1
		} else {
			seqNum = 0
		}
	} else {
		seqNum, err = strconv.Atoi(rawSeqNum)
		if err != nil {
			return "", 0, 0, fmt.Errorf("sequence number atoi error: %v", err)
		}
		if timeMS == stream.topEntry.timeMS && seqNum <= stream.topEntry.seqNum {
			return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

	streamID := fmt.Sprintf("%d-%d", timeMS, seqNum)
	return streamID, timeMS, seqNum, nil
}

func validateAndParseRangeStreamID(rawStreamID string, isStart bool, stream *stream) (int64, int, error) {
	if rawStreamID == "-" {
		return int64(0), 1, nil
	}
	if rawStreamID == "+" {
		return stream.topEntry.timeMS, stream.topEntry.seqNum, nil
	}

	splitted := strings.Split(rawStreamID, "-")
	switch len(splitted) {
	case 1:
		timeMS, err := strconv.ParseInt(rawStreamID, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("milliseconds parse int error: %v", err)
		}
		if isStart {
			return timeMS, 0, nil
		} else {
			return timeMS, math.MaxInt, nil
		}
	case 2:
		timeMS, err := strconv.ParseInt(splitted[0], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("milliseconds parse int error: %v", err)
		}
		seqNum, err := strconv.Atoi(splitted[1])
		if err != nil {
			return 0, 0, fmt.Errorf("sequence number atoi error: %v", err)
		}
		return timeMS, seqNum, nil
	default:
		return 0, 0, fmt.Errorf("wrong stream ID format provided")
	}
}

package memory

import (
	"context"
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
	cond     *sync.Cond
}

type StreamStorage interface {
	baseStorage
	Xadd(streamKey string, requestedStreamID string, entryFields map[string]string) (string, error)
	Xrange(streamKey string, startID string, endID string) ([]EntryWithStreamID, error)
	Xread(streamKeys []string, startIDs []string, timeoutMS int) ([]StreamWithEntries, error)
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

	streamID, timeMS, seqNum, err := stream.validateAndGenerateStreamID(requestedStreamID)
	if err != nil {
		return "", err
	}

	entry := make(entry)
	maps.Copy(entry, entryFields)

	topEntry := topEntry{streamID: streamID, timeMS: timeMS, seqNum: seqNum}
	stream.data[streamID] = entry
	stream.topEntry = topEntry

	stream.cond.Broadcast()

	return streamID, nil
}

type EntryWithStreamID struct {
	StreamID string
	Entry    entry
}

func (ss *streamStorage) Xrange(streamKey string, startID string, endID string) ([]EntryWithStreamID, error) {
	stream := ss.getOrCreateStream(streamKey)
	stream.rwMut.RLock()
	defer stream.rwMut.RUnlock()

	startTimeMS, startSeqNum, err := stream.validateAndParseIntervalStreamID(startID, true)
	if err != nil {
		return nil, fmt.Errorf("start ID validation failed: %v", err)
	}
	endTimeMS, endSeqNum, err := stream.validateAndParseIntervalStreamID(endID, false)
	if err != nil {
		return nil, fmt.Errorf("end ID validation failed: %v", err)
	}

	return stream.traverseEntries(startTimeMS, endTimeMS, startSeqNum, endSeqNum, false), nil
}

type StreamWithEntries struct {
	StreamKey           string
	EntriesWithStreamID []EntryWithStreamID
}

func (ss *streamStorage) Xread(streamKeys []string, startIDs []string, timeoutMS int) ([]StreamWithEntries, error) {
	if timeoutMS < -1 {
		return nil, fmt.Errorf("negative timeout provided")
	}

	streamsWithEntries, err := ss.readStreamsSync(streamKeys, startIDs)
	if err != nil {
		return nil, fmt.Errorf("readStreamsSync error: %s", err)
	}
	if len(streamsWithEntries) > 0 || timeoutMS == -1 {
		return streamsWithEntries, nil
	}

	var timer <-chan time.Time
	if timeoutMS > 0 {
		timer = time.After(time.Millisecond * time.Duration(timeoutMS))
	}

	receivedXadd := make(chan struct{})
	streamError := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i, streamKey := range streamKeys {
		go func(streamKey, startID string) {
			stream := ss.getOrCreateStream(streamKey)
			// Use write locker for sync.Cond
			stream.rwMut.Lock()
			defer stream.rwMut.Unlock()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				entriesWithStreamID, err := stream.read(startID)
				if err != nil {
					streamError <- fmt.Errorf("stream %s read error: %w", streamKey, err)
				}
				if len(entriesWithStreamID) > 0 {
					receivedXadd <- struct{}{}
				}

				stream.cond.Wait()
			}
		}(streamKey, startIDs[i])
	}

	select {
	case err := <-streamError:
		return nil, err
	case <-receivedXadd:
		streamsWithEntries, err := ss.readStreamsSync(streamKeys, startIDs)
		if err != nil {
			return nil, fmt.Errorf("readStreamsSync error: %s", err)
		}
		return streamsWithEntries, nil
	case <-timer:
		return nil, nil
	}
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
	stream.cond = sync.NewCond(&stream.rwMut)
	ss.data[streamKey] = stream
	return stream
}

func (ss *streamStorage) readStreamsSync(streamKeys []string, startIDs []string) ([]StreamWithEntries, error) {
	streamsWithEntries := make([]StreamWithEntries, 0)

	for i, streamKey := range streamKeys {
		stream := ss.getOrCreateStream(streamKey)
		stream.rwMut.RLock()

		entriesWithStreamID, err := stream.read(startIDs[i])
		if err != nil {
			stream.rwMut.RUnlock()
			return nil, fmt.Errorf("stream read err: %s", err)
		}

		if len(entriesWithStreamID) > 0 {
			streamsWithEntries = append(streamsWithEntries, StreamWithEntries{
				StreamKey:           streamKey,
				EntriesWithStreamID: entriesWithStreamID,
			})
		}
		stream.rwMut.RUnlock()
	}

	return streamsWithEntries, nil
}

func (s *stream) read(startID string) ([]EntryWithStreamID, error) {
	startTimeMS, startSeqNum, err := s.validateAndParseIntervalStreamID(startID, true)
	if err != nil {
		return nil, fmt.Errorf("start ID validation failed: %v", err)
	}
	endTimeMS, endSeqNum, _ := s.validateAndParseIntervalStreamID("+", false)

	return s.traverseEntries(startTimeMS, endTimeMS, startSeqNum, endSeqNum, true), nil
}

func (s *stream) traverseEntries(startTimeMS, endTimeMS int64, startSeqNum, endSeqNum int, isExclusive bool) []EntryWithStreamID {
	entriesWithStreamID := make([]EntryWithStreamID, 0)
	timeMS := startTimeMS
	seqNum := startSeqNum
	if isExclusive {
		seqNum += 1
	}

	for timeMS <= endTimeMS {
		if timeMS == endTimeMS && seqNum > endSeqNum {
			break
		}
		streamID := fmt.Sprintf("%d-%d", timeMS, seqNum)
		entry, ok := s.data[streamID]
		if !ok {
			timeMS += 1
			seqNum = 0
			continue
		}
		entriesWithStreamID = append(entriesWithStreamID, EntryWithStreamID{StreamID: streamID, Entry: entry})
		seqNum += 1
	}

	return entriesWithStreamID
}

func (s *stream) validateAndGenerateStreamID(requestedStreamID string) (string, int64, int, error) {
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

	if timeMS < s.topEntry.timeMS {
		return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
	}

	var seqNum int
	if rawSeqNum == "*" {
		if timeMS == s.topEntry.timeMS {
			seqNum = s.topEntry.seqNum + 1
		} else {
			seqNum = 0
		}
	} else {
		seqNum, err = strconv.Atoi(rawSeqNum)
		if err != nil {
			return "", 0, 0, fmt.Errorf("sequence number atoi error: %v", err)
		}
		if timeMS == s.topEntry.timeMS && seqNum <= s.topEntry.seqNum {
			return "", 0, 0, fmt.Errorf("The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

	streamID := fmt.Sprintf("%d-%d", timeMS, seqNum)
	return streamID, timeMS, seqNum, nil
}

func (s *stream) validateAndParseIntervalStreamID(rawStreamID string, isStart bool) (int64, int, error) {
	if rawStreamID == "-" {
		return int64(0), 1, nil
	}
	if rawStreamID == "+" || rawStreamID == "$" {
		return s.topEntry.timeMS, s.topEntry.seqNum, nil
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

func generateStreamID() (streamID string, timeMS int64, seqNum int) {
	seqNum = 0
	timeMS = time.Now().Local().UnixMilli()
	streamID = fmt.Sprintf("%d-%d", timeMS, seqNum)
	return
}

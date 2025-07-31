package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) xread(args []string) resp.Value {
	if len(args) < 3 {
		return resp.SimpleError{Value: "XREAD command must have at least 3 args"}
	}

	streamsStr := args[0]
	if streamsStr != "streams" {
		return resp.SimpleError{Value: fmt.Sprintf("XREAD command undefined word instead of 'streams': %s", streamsStr)}
	}

	keysAndStartIDs := args[1:]
	keysAndStartIDsLen := len(keysAndStartIDs)
	if keysAndStartIDsLen%2 != 0 {
		return resp.SimpleError{Value: fmt.Sprintf("XREAD command must have even count of stream keys and theirs start IDs in summary, got: %d", keysAndStartIDsLen)}
	}

	streamKeys := keysAndStartIDs[:keysAndStartIDsLen/2]
	streamStartIDs := keysAndStartIDs[keysAndStartIDsLen/2:]

	for _, streamKey := range streamKeys {
		if c.storage.KeyExistsWithOtherType(streamKey, memory.TYPE_STREAM) {
			return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
	}

	gotEntries, err := c.storage.StreamStorage().Xread(streamKeys, streamStartIDs)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	respStreamsWithEntries := make([]resp.Value, 0)
	for _, streamWithEntry := range gotEntries {
		respStreamWithEntries := make([]resp.Value, 2)
		respStreamWithEntries[0] = resp.BulkString{Value: &streamWithEntry.StreamKey}
		respEntriesWithStreamID := getRESPEntriesWithStreamID(streamWithEntry.EntriesWithStreamID)
		respStreamWithEntries[1] = resp.Array{Value: respEntriesWithStreamID}

		respStreamsWithEntries = append(respStreamsWithEntries, resp.Array{Value: respStreamWithEntries})
	}
	return resp.Array{Value: respStreamsWithEntries}
}

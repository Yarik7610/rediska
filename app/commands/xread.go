package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) xread(args []string) resp.Value {
	if len(args) < 3 {
		return resp.SimpleError{Value: "XREAD command must have at least 3 args"}
	}

	firstArg := args[0]
	blockArgsOffset := 0
	timeoutMS := -1
	switch firstArg {
	case "streams":
	case "block":
		var err error
		timeoutMS, err = strconv.Atoi(args[1])
		if err != nil {
			return resp.SimpleError{Value: fmt.Sprintf("XREAD timeout (MS) argument atoi error: %v", err)}
		}
		blockArgsOffset = 2
	default:
		return resp.SimpleError{Value: fmt.Sprintf("XREAD command undefined arg: %s", firstArg)}
	}

	keysAndStartIDs := args[blockArgsOffset+1:]
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

	gotEntries, err := c.storage.StreamStorage().Xread(streamKeys, streamStartIDs, timeoutMS)
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

	if len(respStreamsWithEntries) == 0 {
		return resp.BulkString{Value: nil}
	}
	return resp.Array{Value: respStreamsWithEntries}
}

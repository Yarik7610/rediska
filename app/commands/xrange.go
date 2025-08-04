package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) xrange(args []string) resp.Value {
	if len(args) != 3 {
		return resp.SimpleError{Value: "XRANGE command must have 3 args"}
	}

	streamKey := args[0]
	startStreamID := args[1]
	endStreamID := args[2]
	if c.storage.KeyExistsWithOtherType(streamKey, memory.TYPE_STREAM) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	gotEntries, err := c.storage.StreamStorage().Xrange(streamKey, startStreamID, endStreamID)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	return resp.Array{Value: getRESPEntriesWithStreamID(gotEntries)}
}

func getRESPEntriesWithStreamID(entriesWithStreamID []memory.EntryWithStreamID) []resp.Value {
	respEntriesWithStreamID := make([]resp.Value, 0)
	for _, entryWithStreamID := range entriesWithStreamID {
		keyValues := extractKeyValuesToStringSlice(entryWithStreamID.Entry)
		entry := resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &entryWithStreamID.StreamID},
			resp.CreateBulkStringArray(keyValues...),
		}}
		respEntriesWithStreamID = append(respEntriesWithStreamID, entry)
	}
	return respEntriesWithStreamID
}

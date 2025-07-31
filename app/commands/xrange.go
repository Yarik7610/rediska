package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) xrange(args []string) resp.Value {
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

	xrangeEntries := make([]resp.Value, 0)
	for _, gotEntry := range gotEntries {
		keyValues := extractKeyValuesToStringSlice(&gotEntry)
		xrangeEntry := resp.Array{Value: []resp.Value{
			resp.BulkString{Value: &gotEntry.StreamID},
			resp.CreateBulkStringArray(keyValues...),
		}}
		xrangeEntries = append(xrangeEntries, xrangeEntry)
	}
	return resp.Array{Value: xrangeEntries}
}

func extractKeyValuesToStringSlice(gotEntry *memory.XrangeEntry) []string {
	keyValues := make([]string, 0)
	for key, value := range *gotEntry.Entry {
		keyValues = append(keyValues, key, value)
	}
	return keyValues
}

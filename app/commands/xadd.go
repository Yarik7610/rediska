package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *Controller) xadd(args, commandAndArgs []string) resp.Value {
	if len(args) < 4 {
		return resp.SimpleError{Value: "XADD command must have at least 4 args"}
	}

	streamKey := args[0]
	requestedStreamID := args[1]
	if c.storage.KeyExistsWithOtherType(streamKey, memory.TYPE_STREAM) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	entryFields := make(map[string]string)
	rawEntryFields := args[2:]
	rawEntryFieldsLen := len(rawEntryFields)
	if rawEntryFieldsLen%2 != 0 {
		return resp.SimpleError{Value: fmt.Sprintf("XADD wrong entry fields count, need even count, detected count: %d", rawEntryFieldsLen)}
	}
	for i := 0; i < len(rawEntryFields)-1; i++ {
		entryFields[rawEntryFields[i]] = rawEntryFields[i+1]
	}

	gotStreamID, err := c.storage.StreamStorage().Xadd(streamKey, requestedStreamID, entryFields)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("XADD error: %s", err)}
	}

	c.propagateWriteCommand(commandAndArgs)
	return resp.BulkString{Value: &gotStreamID}
}

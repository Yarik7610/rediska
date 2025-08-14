package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) xadd(args, commandAndArgs []string) resp.Value {
	if len(args) < 4 {
		return resp.SimpleError{Value: "XADD command must have at least 4 args"}
	}

	streamKey := args[0]
	requestedStreamID := args[1]
	if c.storage.KeyExistsWithOtherType(streamKey, memory.TYPE_STREAM) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	entryFields, err := parseEntryFields(args[2:])
	if err != nil {
		return resp.SimpleError{Value: err.Error()}
	}

	gotStreamID, err := c.storage.StreamStorage().Xadd(streamKey, requestedStreamID, entryFields)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	c.propagateWriteCommand(commandAndArgs)
	return resp.BulkString{Value: &gotStreamID}
}

func parseEntryFields(rawFields []string) (map[string]string, error) {
	entryFields := make(map[string]string)

	rawEntryFieldsLen := len(rawFields)
	if rawEntryFieldsLen%2 != 0 {
		return nil, fmt.Errorf("XADD wrong entry fields count, need even count, detected count: %d", rawEntryFieldsLen)
	}

	for i := 0; i < len(rawFields)-1; i++ {
		entryFields[rawFields[i]] = rawFields[i+1]
	}
	return entryFields, nil
}

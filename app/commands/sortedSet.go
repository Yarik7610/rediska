package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) zadd(args, commandAndArgs []string) resp.Value {
	if len(args) < 3 {
		return resp.SimpleError{Value: "ZADD command must have at least 3 args"}
	}

	sortedSetKey := args[0]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	members, scores, err := parseMembersAndScores(args[1:])
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	insertedCount := c.storage.SortedSetStorage().Zadd(sortedSetKey, scores, members)

	c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: insertedCount}
}

func parseMembersAndScores(rawFields []string) ([]string, []float64, error) {
	members := make([]string, 0)
	scores := make([]float64, 0)

	rawEntryFieldsLen := len(rawFields)
	if rawEntryFieldsLen%2 != 0 {
		return nil, nil, fmt.Errorf("ZADD wrong member and scores count, need even count, detected count: %d", rawEntryFieldsLen)
	}

	for i := 0; i <= len(rawFields)-2; i += 2 {
		score, err := strconv.ParseFloat(rawFields[i], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("ZADD wrong score format: %v", err)
		}
		scores = append(scores, score)
		members = append(members, rawFields[i+1])
	}
	return members, scores, nil
}

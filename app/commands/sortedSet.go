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

func (c *controller) zrank(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "ZRANK command must have 2 args"}
	}

	sortedSetKey := args[0]
	sortedSetMember := args[1]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	rank := c.storage.SortedSetStorage().Zrank(sortedSetKey, sortedSetMember)

	if rank == -1 {
		return resp.BulkString{Value: nil}
	}
	return resp.Integer{Value: rank}
}

func (c *controller) zrange(args []string) resp.Value {
	if len(args) != 3 {
		return resp.SimpleError{Value: "ZRANGE command must have 3 args"}
	}

	key := args[0]
	if c.storage.KeyExistsWithOtherType(key, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	startIdx := args[1]
	stopIdx := args[2]

	startIdxAtoi, err := strconv.Atoi(startIdx)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ZRANGE command start atoi error: %v", err)}
	}
	stopIdxAtoi, err := strconv.Atoi(stopIdx)
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ZRANGE command stop atoi error: %v", err)}
	}

	values := c.storage.SortedSetStorage().Zrange(key, startIdxAtoi, stopIdxAtoi)
	return resp.CreateBulkStringArray(values...)
}

func (c *controller) zcard(args []string) resp.Value {
	if len(args) != 1 {
		return resp.SimpleError{Value: "ZCARD command must have 1 arg"}
	}

	sortedSetKey := args[0]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	card := c.storage.SortedSetStorage().Zcard(sortedSetKey)
	return resp.Integer{Value: card}
}

func (c *controller) zscore(args []string) resp.Value {
	if len(args) != 2 {
		return resp.SimpleError{Value: "ZSCORE command must have 2 args"}
	}

	sortedSetKey := args[0]
	sortedSetMember := args[1]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	score := c.storage.SortedSetStorage().Zscore(sortedSetKey, sortedSetMember)
	if score == nil {
		return resp.BulkString{Value: nil}
	}

	floatString := strconv.FormatFloat(*score, 'e', -1, 64)
	return resp.BulkString{Value: &floatString}
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

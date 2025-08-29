package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/geo"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (c *controller) geoadd(args, commandAndArgs []string) resp.Value {
	if len(args) < 4 {
		return resp.SimpleError{Value: "GEOADD command must have at least 4 args"}
	}

	sortedSetKey := args[0]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	locations, err := parseLocations(args[1:])
	if err != nil {
		return resp.SimpleError{Value: fmt.Sprintf("ERR %s", err)}
	}

	scores, members := geo.ConvertToScoresAndMembersSlices(locations)
	insertedCount := c.storage.SortedSetStorage().Zadd(sortedSetKey, scores, members)

	c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: insertedCount}
}

func parseLocations(rawFields []string) ([]geo.Location, error) {
	rawEntryFieldsLen := len(rawFields)
	if rawEntryFieldsLen%3 != 0 {
		return nil, fmt.Errorf("GEOADD wrong locations argument count, need count multiple of 3, detected count: %d", rawEntryFieldsLen)
	}

	locations := make([]geo.Location, 0)
	for i := 0; i < rawEntryFieldsLen; i += 3 {
		longitude, err := strconv.ParseFloat(rawFields[i], 64)
		if err != nil {
			return nil, fmt.Errorf("GEOADD wrong location longitude format: %v", err)
		}
		latitude, err := strconv.ParseFloat(rawFields[i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("GEOADD wrong location latitude format: %v", err)
		}

		if !geo.ValidLongitude(longitude) || !geo.ValidLatitude(latitude) {
			return nil, fmt.Errorf("invalid longitude,latitude pair %f,%f", longitude, latitude)
		}

		locations = append(locations, geo.Location{
			Longitude: longitude,
			Latitude:  latitude,
			Member:    rawFields[i+2],
		})
	}

	return locations, nil
}

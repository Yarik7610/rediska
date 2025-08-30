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

	scores, members := convertToScoresAndMembersSlices(c.geoController, locations)
	insertedCount := c.storage.SortedSetStorage().Zadd(sortedSetKey, scores, members)

	c.propagateWriteCommand(commandAndArgs)
	return resp.Integer{Value: insertedCount}
}

func (c *controller) geopos(args []string) resp.Value {
	if len(args) < 2 {
		return resp.SimpleError{Value: "GEOPOS command must have at least 2 args"}
	}

	sortedSetKey := args[0]
	if c.storage.KeyExistsWithOtherType(sortedSetKey, memory.TYPE_SORTED_SET) {
		return resp.SimpleError{Value: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	members := args[1:]
	multipleRESPResponses := make([]resp.Value, 0)
	for _, member := range members {
		score := c.storage.SortedSetStorage().Zscore(sortedSetKey, member)
		if score == nil {
			multipleRESPResponses = append(multipleRESPResponses, resp.Array{Value: nil})
			continue
		}

		location := c.geoController.Decode(uint64(*score))
		longitudeString := strconv.FormatFloat(location.Longitude, 'f', -1, 64)
		latitudeString := strconv.FormatFloat(location.Latitude, 'f', -1, 64)
		multipleRESPResponses = append(multipleRESPResponses, resp.CreateBulkStringArray(longitudeString, latitudeString))
	}

	return resp.Array{Value: multipleRESPResponses}
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

func convertToScoresAndMembersSlices(geoController geo.Controller, locations []geo.Location) ([]float64, []string) {
	scores := make([]float64, 0)
	members := make([]string, 0)

	for _, location := range locations {
		scores = append(scores, float64(geoController.Encode(&location)))
		members = append(members, location.Member)
	}

	return scores, members
}

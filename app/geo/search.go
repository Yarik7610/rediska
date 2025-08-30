package geo

import (
	"fmt"
)

type SearchOptions struct {
	FromLonLatOption *FromLonLatOption
	ByRadiusOption   *ByRadiusOption
	Shape            Shape
}

func NewSearchOptions() *SearchOptions {
	return &SearchOptions{Shape: -1}
}

type FromLonLatOption struct {
	Longitude float64
	Latitude  float64
}

type ByRadiusOption struct {
	Value float64
	Unit  string
}

type Shape int

const RADIUS_SHAPE Shape = iota

func (c controller) InArea(searchOptions *SearchOptions, startLocation, location *Location) (bool, error) {
	switch searchOptions.Shape {
	case RADIUS_SHAPE:
		distance := c.Dist(startLocation, location)
		radiusMeters, err := ConvertRangeToMeters(searchOptions.ByRadiusOption.Value, searchOptions.ByRadiusOption.Unit)
		if err != nil {
			return false, err
		}
		return distance <= radiusMeters, nil
	default:
		return false, fmt.Errorf("unknown area shape to search in")
	}
}

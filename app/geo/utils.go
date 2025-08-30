package geo

import (
	"fmt"
	"slices"
	"strings"
)

func ValidLongitude(longitude float64) bool {
	return longitude <= MAX_LONGITUDE && longitude >= MIN_LONGITUDE
}

func ValidLatitude(latitude float64) bool {
	return latitude <= MAX_LATITUDE && latitude >= MIN_LATITUDE
}

func ValidRangeUnit(unit string) bool {
	return slices.Contains(RANGE_UNITS, strings.ToLower(unit))
}

func ConvertRangeToMeters(value float64, unit string) (float64, error) {
	u := strings.ToLower(unit)
	switch u {
	case "m":
		return value, nil
	case "km":
		return value * 1000, nil
	case "mi":
		return value * METERS_PER_MILE, nil
	case "ft":
		return value * METERS_PER_FEET, nil
	default:
		return 0, fmt.Errorf("invalid range unit type: %s", unit)
	}
}

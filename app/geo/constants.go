package geo

const (
	MIN_LONGITUDE = -180
	MAX_LONGITUDE = 180
	MIN_LATITUDE  = -85.05112878
	MAX_LATITUDE  = 85.05112878

	LATITUDE_RANGE  = MAX_LATITUDE - MIN_LATITUDE
	LONGITUDE_RANGE = MAX_LONGITUDE - MIN_LONGITUDE
)

const (
	// Point has (X, Y) coordinates
	// Mantissa of float64 is 52 bits, so we can safely store 52 / 2 = 26 bits
	BITS_PER_COORD        = 26
	MAX_DECODE_DIFFERENCE = 0.000001
	EARTH_RADIUS_METERS   = 6372797.560856
)

var RANGE_UNITS = []string{"m", "km", "mi", "ft"}

const (
	METERS_PER_MILE = 1609.344
	METERS_PER_FEET = 0.3048
)

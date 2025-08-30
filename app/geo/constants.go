package geo

const (
	MIN_LONGITUDE = -180
	MAX_LONGITUDE = 180
	MIN_LATITUDE  = -85.05112878
	MAX_LATITUDE  = 85.05112878

	// Point has (X, Y) coordinates
	// Mantissa of float64 is 52 bits, so we can safely store 52 / 2 = 26 bits
	BITS_PER_COORD = 26

	LATITUDE_RANGE  = MAX_LATITUDE - MIN_LATITUDE
	LONGITUDE_RANGE = MAX_LONGITUDE - MIN_LONGITUDE
)

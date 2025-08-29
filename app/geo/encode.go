package geo

func (c controller) Encode(latitude, longitude float64) uint64 {
	// Normalize to the range [0, 2^26)
	normalizedLongitude := 1 << BITS_PER_COORD * (longitude - MIN_LONGITUDE) / LONGITUDE_RANGE
	normalizedLatitude := 1 << BITS_PER_COORD * (latitude - MIN_LATITUDE) / LATITUDE_RANGE

	// Traverse to uint32 (not int, because normalized to [0..) to perform bit operations
	normalizedLongitudeUint := uint32(normalizedLongitude)
	normalizedLatitudeUint := uint32(normalizedLatitude)

	return c.interleave(normalizedLatitudeUint, normalizedLongitudeUint)
}

func (c controller) interleave(x uint32, y uint32) uint64 {
	xSpreaded := c.spreadUint32ToUint64(x)
	ySpreaded := c.spreadUint32ToUint64(y)
	yShifted := ySpreaded << 1
	return yShifted | xSpreaded
}

// We expand our number to uint64 and also alternate it with zeroes
// For instance:
// 5 (32 bits) = ...101,
// 5 (64 bits, spreaded with zeroes) = ...010001,
func (c controller) spreadUint32ToUint64(val uint32) uint64 {
	result := uint64(val)
	result = (result | (result << 16)) & 0x0000FFFF0000FFFF
	result = (result | (result << 8)) & 0x00FF00FF00FF00FF
	result = (result | (result << 4)) & 0x0F0F0F0F0F0F0F0F
	result = (result | (result << 2)) & 0x3333333333333333
	result = (result | (result << 1)) & 0x5555555555555555
	return result
}

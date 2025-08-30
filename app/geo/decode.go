package geo

func (c controller) Decode(code uint64) *Location {
	y := code >> 1
	x := code

	yNormalized := compactUint64ToUint32(y)
	xNormalized := compactUint64ToUint32(x)

	return c.convertNormalizedNumsToCoordinates(xNormalized, yNormalized)
}

func (controller) convertNormalizedNumsToCoordinates(xNormalized, yNormalized uint32) *Location {
	// We were storing uint [0, 2^26). Each value (grid) here isn't float
	// So we can't give back a precise result, i.g 54.231312851
	// To make it +- balanced and accurate,
	// we need the center of our uint value (grid) in representation of latitude / longitude range
	upperEncodeBoundary := 1 << BITS_PER_COORD
	gridLatitudeMin := MIN_LATITUDE + LATITUDE_RANGE*(float64(xNormalized)/float64(upperEncodeBoundary))
	gridLatitudeMax := MIN_LATITUDE + LATITUDE_RANGE*(float64(xNormalized+1)/float64(upperEncodeBoundary))
	gridLongitudeMin := MIN_LONGITUDE + LONGITUDE_RANGE*(float64(yNormalized)/float64(upperEncodeBoundary))
	gridLongitudeMax := MIN_LONGITUDE + LONGITUDE_RANGE*(float64(yNormalized+1)/float64(upperEncodeBoundary))

	latitude := (gridLatitudeMin + gridLatitudeMax) / 2
	longitude := (gridLongitudeMin + gridLongitudeMax) / 2

	return &Location{Longitude: longitude, Latitude: latitude}
}

// Collapse our interleaved number to uint32 removing alternated zeroes
func compactUint64ToUint32(v uint64) uint32 {
	result := v & 0x5555555555555555
	result = (result | (result >> 1)) & 0x3333333333333333
	result = (result | (result >> 2)) & 0x0F0F0F0F0F0F0F0F
	result = (result | (result >> 4)) & 0x00FF00FF00FF00FF
	result = (result | (result >> 8)) & 0x0000FFFF0000FFFF
	result = (result | (result >> 16)) & 0x00000000FFFFFFFF
	return uint32(result)
}

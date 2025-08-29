package geo

func ValidLongitude(longitude float64) bool {
	return longitude <= MAX_LONGITUDE && longitude >= MIN_LONGITUDE
}

func ValidLatitude(latitude float64) bool {
	return latitude <= MAX_LATITUDE && latitude >= MIN_LATITUDE
}

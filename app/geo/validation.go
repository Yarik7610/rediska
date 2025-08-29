package geo

func ValidLongitude(longitude float64) bool {
	return longitude <= 180 && longitude >= -180
}

func ValidLatitude(latitude float64) bool {
	return latitude <= 85.05112878 && latitude >= -85.05112878
}

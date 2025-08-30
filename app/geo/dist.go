package geo

import (
	"math"
)

type radianPoint struct {
	latitude  float64
	longitude float64
}

func (controller) Dist(location1, location2 *Location) float64 {
	radianPoint1 := toRadianPoint(location1)
	radianPoint2 := toRadianPoint(location2)

	radianDeltaLatitude := radianPoint2.latitude - radianPoint1.latitude
	radianDeltaLongitude := radianPoint2.longitude - radianPoint1.longitude

	x := haversine(radianDeltaLatitude) + math.Cos(radianPoint1.latitude)*math.Cos(radianPoint2.latitude)*haversine(radianDeltaLongitude)
	return 2 * EARTH_RADIUS_METERS * math.Asin(math.Sqrt(x))
}

func toRadianPoint(location *Location) *radianPoint {
	return &radianPoint{
		latitude:  location.Latitude * math.Pi / 180,
		longitude: location.Longitude * math.Pi / 180,
	}
}

func haversine(theta float64) float64 {
	return 0.5 * (1 - math.Cos(theta))
}

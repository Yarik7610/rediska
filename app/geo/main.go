package geo

type Controller interface {
	Encode(location *Location) uint64
	Decode(code uint64) *Location
	Dist(location1, location2 *Location) float64
}

type controller struct{}

func NewController() Controller {
	return controller{}
}

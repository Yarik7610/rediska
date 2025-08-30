package geo

type Controller interface {
	Encode(location *Location) uint64
	Decode(code uint64) *Location
}

type controller struct{}

func NewController() Controller {
	return controller{}
}

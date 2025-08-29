package geo

type Controller interface {
	Encode(longitude, latitude float64) uint64
	Decode(code uint64) (float64, float64)
}

type controller struct{}

func NewController() Controller {
	return controller{}
}

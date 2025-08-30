package geo

type FromLonLatOption struct {
	Longitude float64
	Latitude  float64
}
type ByRadiusOption struct {
	Value float64
	Unit  string
}

type SearchOptions struct {
	FromLonLatOption *FromLonLatOption
	ByRadiusOption   *ByRadiusOption
}

package geo

func ConvertToScoresAndMembersSlices(locations []Location) ([]float64, []string) {
	scores := make([]float64, 0)
	members := make([]string, 0)

	for _, location := range locations {
		scores = append(scores, countScore(location.Longitude, location.Latitude))
		members = append(members, location.Member)
	}

	return scores, members
}

func countScore(longitude, latitude float64) float64 {
	return 0
}

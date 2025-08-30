package geo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDist(t *testing.T) {
	c := controller{}

	tests := []struct {
		name           string
		loc1           *Location
		loc2           *Location
		expectedMeters float64
	}{
		{
			"Bangkok - Beijing",
			&Location{Latitude: 13.722000686932997, Longitude: 100.52520006895065},
			&Location{Latitude: 39.9075003315814, Longitude: 116.39719873666763},
			3299195.1357588563,
		},
		{
			"Berlin - Copenhagen",
			&Location{Latitude: 52.52439934649943, Longitude: 13.410500586032867},
			&Location{Latitude: 55.67589927498264, Longitude: 12.56549745798111},
			354828.2423250971,
		},
		{
			"New York - London",
			&Location{Latitude: 40.712798986951505, Longitude: -74.00600105524063},
			&Location{Latitude: 51.50740077990134, Longitude: -0.12779921293258667},
			5571793.967617199,
		},
		{
			"Tokyo - Sydney",
			&Location{Latitude: 35.68950126697936, Longitude: 139.691701233387},
			&Location{Latitude: -33.86880091934156, Longitude: 151.2092998623848},
			7828823.534197976,
		},
		{
			"Paris - Vienna",
			&Location{Latitude: 48.85340071224621, Longitude: 2.348802387714386},
			&Location{Latitude: 48.20640046271915, Longitude: 16.370699107646942},
			1033845.7464546118,
		},
	}

	for _, test := range tests {
		actual := fmt.Sprintf("%0.4f", c.Dist(test.loc1, test.loc2))
		expected := fmt.Sprintf("%0.4f", test.expectedMeters)
		assert.Equal(t, expected, actual, test.name)
	}
}

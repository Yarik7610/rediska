package geo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	c := controller{}

	tests := []struct {
		name          string
		latitude      float64
		longitude     float64
		expectedScore uint64
	}{
		{"Bangkok", 13.722000686932997, 100.52520006895065, 3962257306574459},
		{"Beijing", 39.9075003315814, 116.39719873666763, 4069885364908765},
		{"Berlin", 52.52439934649943, 13.410500586032867, 3673983964876493},
		{"Copenhagen", 55.67589927498264, 12.56549745798111, 3685973395504349},
		{"New Delhi", 28.666698899347338, 77.21670180559158, 3631527070936756},
		{"Kathmandu", 27.701700137333084, 85.3205993771553, 3639507404773204},
		{"London", 51.50740077990134, -0.12779921293258667, 2163557714755072},
		{"New York", 40.712798986951505, -74.00600105524063, 1791873974549446},
		{"Paris", 48.85340071224621, 2.348802387714386, 3663832752681684},
		{"Sydney", -33.86880091934156, 151.2092998623848, 3252046221964352},
		{"Tokyo", 35.68950126697936, 139.691701233387, 4171231230197045},
		{"Vienna", 48.20640046271915, 16.370699107646942, 3673109836391743},
	}

	for _, test := range tests {
		location := &Location{Longitude: test.longitude, Latitude: test.latitude}
		actualScore := c.Encode(location)
		assert.Equal(t, actualScore, test.expectedScore)
	}
}

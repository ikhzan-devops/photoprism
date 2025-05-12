package tz

import (
	"github.com/ugjka/go-tz/v2"
)

// Position returns the time zone for the given latitude and longitude.
func Position(lat, lng float64) (name string) {
	if lat == 0.0 || lng == 0.0 {
		return ""
	}

	zones, err := tz.GetZone(tz.Point{
		Lat: lat,
		Lon: lng,
	})

	if err != nil || len(zones) == 0 {
		return ""
	}

	name = Name(zones[0])

	if name == UTC || name == Local {
		return ""
	}

	return name
}

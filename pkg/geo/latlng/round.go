package latlng

import "math"

var RoundDecimals = float64(10000000)

// Round rounds the given coordinate to six decimal places.
func Round(c float64) float64 {
	return math.Round(c*RoundDecimals) / RoundDecimals
}

// RoundCoords rounds the given latitude and longitude to six decimal places.
func RoundCoords(lat, lng float64) (float64, float64) {
	return Round(lat), Round(lng)
}

package txt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("Berlin", func(t *testing.T) {
		lat, lng, err := Position("48.4565496,35.0719864")
		assert.NoError(t, err)
		assert.InDelta(t, lat, 48.45, 0.01)
		assert.InDelta(t, lng, 35.07, 0.01)
	})
	t.Run("China", func(t *testing.T) {
		lat, lng, err := Position("39.892215944444445, 116.31463963888889")
		assert.NoError(t, err)
		assert.InDelta(t, lat, 39.89, 0.01)
		assert.InDelta(t, lng, 116.31, 0.01)
	})
	t.Run("California", func(t *testing.T) {
		lat, lng, err := Position("+37.75326666666667, -122.42250833333334")
		assert.NoError(t, err)
		assert.InDelta(t, lat, 37.75, 0.01)
		assert.InDelta(t, lng, -122.42, 0.01)
	})
	t.Run("Invalid", func(t *testing.T) {
		_, _, err := Position("+91.75326666666667, -122.42250833333334")
		assert.Error(t, err)
		_, _, err = Position("48.4565496,190.0719864")
		assert.Error(t, err)
	})
}

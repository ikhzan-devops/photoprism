package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstantQuality(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		assert.Equal(t, "100", QvQuality(BestQuality))
		assert.Equal(t, "50", QvQuality(DefaultQuality))
		assert.Equal(t, "1", QvQuality(WorstQuality))
	})
}

func TestGlobalQuality(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		assert.Equal(t, "1", GlobalQuality(BestQuality))
		assert.Equal(t, "25", GlobalQuality(DefaultQuality))
		assert.Equal(t, "49", GlobalQuality(WorstQuality))
	})
}

func TestCrfQuality(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		assert.Equal(t, "0", CrfQuality(BestQuality))
		assert.Equal(t, "25", CrfQuality(DefaultQuality))
		assert.Equal(t, "49", CrfQuality(WorstQuality))
	})
}

func TestQpQuality(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		assert.Equal(t, "0", QpQuality(BestQuality))
		assert.Equal(t, "25", QpQuality(DefaultQuality))
		assert.Equal(t, "49", QpQuality(WorstQuality))
	})
}

func TestCqQuality(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		assert.Equal(t, "1", CqQuality(BestQuality))
		assert.Equal(t, "25", CqQuality(DefaultQuality))
		assert.Equal(t, "49", CqQuality(WorstQuality))
	})
}

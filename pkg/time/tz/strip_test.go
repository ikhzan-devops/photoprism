package tz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrip(t *testing.T) {
	t.Run("Local", func(t *testing.T) {
		assert.Equal(t, time.Time(time.Date(1990, time.April, 18, 1, 0, 0, 0, time.UTC)), Strip(time.Date(1990, 4, 18, 1, 0, 0, 0, time.Local)))
	})
}

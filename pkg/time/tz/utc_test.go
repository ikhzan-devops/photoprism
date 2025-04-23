package tz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUTC(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		assert.True(t, IsUTC("utc"))
		assert.True(t, IsUTC("UTC"))
		assert.True(t, IsUTC("Z"))
		assert.True(t, IsUTC("Zulu"))
	})
	t.Run("False", func(t *testing.T) {
		assert.False(t, IsUTC(""))
		assert.False(t, IsUTC("GMT"))
		assert.False(t, IsUTC("local"))
		assert.False(t, IsUTC("Local"))
	})
}

func TestTruncateUTC(t *testing.T) {
	now := time.Now().In(TimeLocal)
	ns := now.Nanosecond()

	result := TruncateUTC(now)
	assert.Equal(t, TimeUTC, result.Location())
	timeZone, _ := result.Zone()
	assert.Equal(t, UTC, timeZone)
	assert.Equal(t, 0, result.Nanosecond())
	if ns > 0 {
		assert.NotEqual(t, ns, result.Nanosecond())
	}
}

package tz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsLocal(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		assert.True(t, IsLocal(Unknown))
		assert.True(t, IsLocal(Local))
		assert.True(t, IsLocal("local"))
		assert.True(t, IsLocal("LOCAL"))
	})
	t.Run("False", func(t *testing.T) {
		assert.False(t, IsLocal("utc"))
		assert.False(t, IsLocal(UTC))
		assert.False(t, IsLocal(GMT))
		assert.False(t, IsLocal(Zulu))
	})
}

func TestTruncateLocal(t *testing.T) {
	now := time.Now().In(TimeLocal)
	ns := now.Nanosecond()

	result := TruncateLocal(now)
	timeZone, _ := result.Zone()
	assert.Equal(t, Local, timeZone)
	assert.Equal(t, 0, result.Nanosecond())
	if ns > 0 {
		assert.NotEqual(t, ns, result.Nanosecond())
	}
}

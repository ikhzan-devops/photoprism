package tz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	t.Run("UTC", func(t *testing.T) {
		assert.Equal(t, time.UTC.String(), Find(time.UTC.String()).String())
		assert.Equal(t, time.UTC.String(), Find("Z").String())
		assert.Equal(t, time.UTC.String(), Find("UTC").String())
		assert.Equal(t, time.UTC, Find("UTC"))
		assert.Equal(t, "UTC", Find("0").String())
		assert.Equal(t, "UTC", Find("UTC+0").String())
		assert.Equal(t, "UTC", Find("UTC+00:00").String())
	})
	t.Run("GMT", func(t *testing.T) {
		assert.Equal(t, "GMT", Find("GMT").String())
		assert.Equal(t, "GMT", Find("Etc/GMT").String())
	})
	t.Run("Local", func(t *testing.T) {
		assert.Equal(t, "Local", Find("").String())
		assert.Equal(t, TimeLocal, Find(""))
		assert.Equal(t, "Local", Find("Local").String())
		assert.Equal(t, TimeLocal, Find("Local"))
	})
	t.Run("Berlin", func(t *testing.T) {
		assert.Equal(t, "Europe/Berlin", Find("Europe/Berlin").String())
	})
	t.Run("Offset", func(t *testing.T) {
		local, err := time.Parse("2006-01-02 15:04:05 Z07:00", "2023-10-02 13:20:17 +00:00")

		if err != nil {
			t.Fatal(err)
		}

		utc, err := time.Parse("2006-01-02 15:04:05 Z07:00", "2023-10-02 11:20:17 +00:00")

		if err != nil {
			t.Fatal(err)
		}

		timeZone := UtcOffset(utc, local, "")

		assert.Equal(t, "UTC+2", timeZone)

		loc := Find(timeZone)

		assert.Equal(t, "UTC+2", loc.String())
	})
}

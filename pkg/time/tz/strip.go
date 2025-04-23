package tz

import (
	"time"
)

// Strip removes the time zone from a time.
func Strip(t time.Time) (result time.Time) {
	if t.IsZero() {
		return t
	}

	result, _ = time.ParseInLocation("2006:01:02 15:04:05", t.Format("2006:01:02 15:04:05"), time.UTC)

	return result.Truncate(time.Second)
}

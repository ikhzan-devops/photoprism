package tz

import (
	"strings"
	"time"
)

// IsUTC returns true if the time zone string represents Universal Coordinated Time (UTC).
func IsUTC(s string) bool {
	if s == Unknown || len(s) > 7 {
		return false
	}

	s = strings.ToUpper(s)

	return s == UTC || s == Zulu || s == "ZULU" || s == "ETC/UTC"
}

// TruncateUTC sets the time zone to UTC and changes the precision to seconds.
func TruncateUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return t.UTC().Truncate(time.Second)
}

// LocationUTC returns the time at the locale with the time zone set to UTC.
func LocationUTC(t time.Time, loc *time.Location) (result time.Time) {
	var err error
	if result, err = time.ParseInLocation("2006-01-02T15:04:05", t.In(loc).Format("2006-01-02T15:04:05"), time.UTC); err != nil {
		return result
	} else {
		return t
	}
}

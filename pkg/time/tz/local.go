package tz

import (
	"strings"
	"time"
)

// IsLocal returns true if the time zone string represents Local time.
func IsLocal(s string) bool {
	if s == Unknown {
		return true
	} else if len(s) != len(Local) {
		return false
	}

	return strings.EqualFold(s, Local)
}

// TruncateLocal changes the precision of Local Time to full seconds to avoid jitter.
func TruncateLocal(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	return t.Truncate(time.Second)
}

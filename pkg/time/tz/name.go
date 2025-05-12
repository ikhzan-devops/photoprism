package tz

import (
	"strings"
)

// MaxLen specifies the maximum length of time zone strings.
const MaxLen = 64

// Name normalizes the specified time zone string and returns it.
func Name(s string) string {
	s = strings.TrimSpace(s)

	// Detect and return standard time zones.
	if IsUTC(s) {
		return UTC
	} else if IsLocal(s) {
		return Local
	}

	// Clip to max length.
	if len(s) > MaxLen {
		s = s[:MaxLen]
	}

	// Handle time zone offset strings.
	if zone := NormalizeUtcOffset(s); zone != "" {
		return zone
	}

	return s
}

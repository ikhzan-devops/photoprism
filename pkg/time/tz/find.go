package tz

import (
	"fmt"
	"time"
)

// Find returns the matching time zone location.
func Find(name string) *time.Location {
	if IsUTC(name) {
		return time.UTC
	} else if IsLocal(name) {
		return TimeLocal
	}

	// Normalize zone name.
	name = Name(name)

	if offsetSec, offsetErr := Offset(name); offsetErr != nil {
		// Do nothing.
	} else if h := offsetSec / 3600; h > 0 || h < 0 {
		return time.FixedZone(fmt.Sprintf("UTC%+d", h), offsetSec)
	}

	// Find location by name.
	if loc, err := time.LoadLocation(name); err != nil || loc == nil {
		// Return Local location if not found.
		return TimeLocal
	} else {
		// Return location.
		return loc
	}
}

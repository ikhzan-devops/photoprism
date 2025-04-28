package tz

import (
	"time"
)

var (
	TimeUTC   = time.UTC
	TimeLocal = time.FixedZone(Local, 0)
)

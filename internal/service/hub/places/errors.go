package places

import (
	"errors"
)

var ErrMissingQuery = errors.New("missing query")
var ErrMissingCoordinates = errors.New("missing coordinates")

package hooks

import (
	"encoding/json"
	"time"

	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/clean"
)

// Payload represents a webhook payload.
type Payload struct {
	Type      string     `form:"type" json:"type"`
	Timestamp time.Time  `form:"timestamp" json:"timestamp,omitempty"`
	Data      event.Data `form:"data" json:"data"`
}

// JSON returns the payload data as JSON-encoded bytes.
func (p *Payload) JSON() (b []byte) {
	b, jsonErr := json.Marshal(p)

	if jsonErr != nil {
		log.Warningf("hook: %s (json encode)", clean.Error(jsonErr))
	}

	return b
}

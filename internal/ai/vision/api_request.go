package vision

import (
	"encoding/json"

	"github.com/photoprism/photoprism/pkg/rnd"
)

type Files = []string

// ApiRequest represents a Vision API service request.
type ApiRequest struct {
	Id     string `form:"id" yaml:"Id,omitempty" json:"id,omitempty"`
	Model  string `form:"model" yaml:"Model,omitempty" json:"model,omitempty"`
	Images Files  `form:"images" yaml:"Images,omitempty" json:"images,omitempty"`
}

// GetId returns the request ID string and generates a random ID if none was set.
func (r *ApiRequest) GetId() string {
	if r.Id == "" {
		r.Id = rnd.UUID()
	}

	return r.Id
}

// MarshalJSON returns request as JSON.
func (r *ApiRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}

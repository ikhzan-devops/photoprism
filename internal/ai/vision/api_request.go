package vision

import (
	"github.com/photoprism/photoprism/pkg/rnd"
)

// ApiRequest represents a Vision API service request.
type ApiRequest struct {
	Id     string   `form:"id" yaml:"Id,omitempty" json:"id,omitempty"`
	Model  string   `form:"model" yaml:"Model,omitempty" json:"model,omitempty"`
	Images []string `form:"images" yaml:"Images,omitempty" json:"images,omitempty"`
	Videos []string `form:"videos" yaml:"Videos,omitempty" json:"videos,omitempty"`
}

// GetId returns the request ID string and generates a random ID if none was set.
func (r *ApiRequest) GetId() string {
	if r.Id == "" {
		r.Id = rnd.UUID()
	}

	return r.Id
}

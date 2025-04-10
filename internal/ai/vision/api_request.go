package vision

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/photoprism/photoprism/internal/api/download"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
	"github.com/photoprism/photoprism/pkg/rnd"
)

type Files = []string

// ApiRequest represents a Vision API service request.
type ApiRequest struct {
	Id     string `form:"id" yaml:"Id,omitempty" json:"id,omitempty"`
	Model  string `form:"model" yaml:"Model,omitempty" json:"model,omitempty"`
	Images Files  `form:"images" yaml:"Images,omitempty" json:"images,omitempty"`
}

// NewClientRequest returns a new Vision API request with the specified file payload and scheme.
func NewClientRequest(images Files, fileScheme string) (*ApiRequest, error) {
	imageUrls := make(Files, len(images))

	if fileScheme == scheme.Https && !strings.HasPrefix(DownloadUrl, "https://") {
		log.Tracef("vision: file request scheme changed from https to data because https is not configured")
		fileScheme = scheme.Data
	}

	for i := range images {
		switch fileScheme {
		case scheme.Https:
			if id, err := download.Register(images[i]); err != nil {
				return nil, fmt.Errorf("%s (register download)", err)
			} else {
				imageUrls[i] = fmt.Sprintf("%s/%s", DownloadUrl, id)
			}
		case scheme.Data:
			if file, err := os.Open(images[i]); err != nil {
				return nil, fmt.Errorf("%s (create data url)", err)
			} else {
				imageUrls[i] = media.DataUrl(file)
			}
		default:
			return nil, fmt.Errorf("invalid file scheme %s", clean.Log(fileScheme))
		}
	}

	return &ApiRequest{
		Id:     rnd.UUID(),
		Model:  "",
		Images: imageUrls,
	}, nil
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

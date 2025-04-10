package vision

import (
	"net/http"
)

// Default computer vision model configuration.
var (
	NasnetModel = &Model{
		Type:       ModelTypeLabels,
		Name:       "NASNet",
		Version:    "Mobile",
		Resolution: 224,
		Tags:       []string{"photoprism"},
	}
	NsfwModel = &Model{
		Type:       ModelTypeNsfw,
		Name:       "Nsfw",
		Version:    "",
		Resolution: 224,
		Tags:       []string{"serve"},
	}
	FacenetModel = &Model{
		Type:       ModelTypeFaceEmbeddings,
		Name:       "FaceNet",
		Version:    "",
		Resolution: 160,
		Tags:       []string{"serve"},
	}
	CaptionModel = &Model{
		Type:       ModelTypeCaption,
		Name:       "Caption",
		Uri:        "http://photoprism-vision/api/v1/vision/describe",
		Method:     http.MethodPost,
		Resolution: 720,
	}
	DefaultModels     = Models{NasnetModel, NsfwModel, FacenetModel, CaptionModel}
	DefaultThresholds = Thresholds{Confidence: 10}
)

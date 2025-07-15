package vision

import (
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

// Default computer vision model configuration.
var (
	NasnetModel = &Model{
		Type:       ModelTypeLabels,
		Name:       "NASNet",
		Version:    ModelVersionMobile,
		Resolution: 224,
		Tags:       []string{"photoprism"},
	}
	NsfwModel = &Model{
		Type:       ModelTypeNsfw,
		Name:       "Nsfw",
		Version:    ModelVersionNone,
		Resolution: 224,
		Tags:       []string{"serve"},
	}
	FacenetModel = &Model{
		Type:       ModelTypeFace,
		Name:       "FaceNet",
		Version:    ModelVersionNone,
		Resolution: 160,
		Tags:       []string{"serve"},
	}
	CaptionModel = &Model{
		Type:    ModelTypeCaption,
		Name:    CaptionModelDefault,
		Version: ModelVersionLatest,
		Prompt:  CaptionPromptDefault,
		Service: Service{
			// Uri:            "http://photoprism-vision:5000/api/v1/vision/caption",
			FileScheme:     scheme.Data,
			RequestFormat:  ApiFormatVision,
			ResponseFormat: ApiFormatVision,
		},
	}
	DefaultModels     = Models{NasnetModel, NsfwModel, FacenetModel, CaptionModel}
	DefaultThresholds = Thresholds{Confidence: 10}
)

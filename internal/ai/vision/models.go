package vision

import (
	"github.com/photoprism/photoprism/pkg/media/http/scheme"
)

// Default computer vision model configuration.
var (
	NasnetModel = &Model{
		Type:       ModelTypeLabels,
		Name:       "nasnet",
		Version:    VersionMobile,
		Resolution: 224, // Cropped image tile with 224x224 pixels.
		Tags:       []string{"photoprism"},
	}
	NsfwModel = &Model{
		Type:       ModelTypeNsfw,
		Name:       "nsfw",
		Version:    VersionLatest,
		Resolution: 224, // Cropped image tile with 224x224 pixels.
		Tags:       []string{"serve"},
	}
	FacenetModel = &Model{
		Type:       ModelTypeFace,
		Name:       "facenet",
		Version:    VersionLatest,
		Resolution: 160, // Cropped image tile with 160x160 pixels.
		Tags:       []string{"serve"},
	}
	CaptionModel = &Model{
		Type:       ModelTypeCaption,
		Name:       CaptionModelDefault,
		Version:    VersionLatest,
		Resolution: 720, // Original aspect ratio, with a max size of 720 x 720 pixels.
		Prompt:     CaptionPromptDefault,
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

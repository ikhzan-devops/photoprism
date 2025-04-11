package vision

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
		Resolution: 224,
		Uri:        "http://photoprism-vision:5000/api/v1/vision/caption",
	}
	DefaultModels     = Models{NasnetModel, NsfwModel, FacenetModel, CaptionModel}
	DefaultThresholds = Thresholds{Confidence: 10}
)

package vision

import (
	"github.com/photoprism/photoprism/pkg/fs"
)

var (
	AssetsPath        = fs.Abs("../../../assets")
	FaceNetModelPath  = fs.Abs("../../../assets/facenet")
	NsfwModelPath     = fs.Abs("../../../assets/nsfw")
	CachePath         = fs.Abs("../../../storage/cache")
	ServiceUri        = ""
	ServiceKey        = ""
	DownloadUrl       = ""
	DefaultResolution = 224
)

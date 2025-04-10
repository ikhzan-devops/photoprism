package vision

import (
	"time"

	"github.com/photoprism/photoprism/pkg/fs"
)

var (
	AssetsPath        = fs.Abs("../../../assets")
	FaceNetModelPath  = fs.Abs("../../../assets/facenet")
	NsfwModelPath     = fs.Abs("../../../assets/nsfw")
	CachePath         = fs.Abs("../../../storage/cache")
	ServiceUri        = ""
	ServiceKey        = ""
	ServiceTimeout    = time.Minute
	DownloadUrl       = ""
	DefaultResolution = 224
)

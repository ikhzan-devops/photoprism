package meta

import (
	"github.com/photoprism/photoprism/pkg/media/video"
)

const CodecUnknown = ""
const CodecAv1 = video.CodecAV1
const CodecAvc1 = video.CodecAVC
const CodecHvc1 = video.CodecHEVC
const CodecJpeg = "jpeg"
const CodecHeic = "heic"
const CodecXMP = "xmp"

// CodecAvc returns true if the video codec is AVC.
func (data Data) CodecAvc() bool {
	return data.Codec == CodecAvc1
}

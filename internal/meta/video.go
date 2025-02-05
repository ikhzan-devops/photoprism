package meta

import (
	"github.com/photoprism/photoprism/pkg/media/video"
)

const CodecUnknown = ""
const CodecJpeg = "jpeg"
const CodecHeic = "heic"
const CodecXMP = "xmp"

// CodecAvc returns true if the video codec is AVC.
func (data Data) CodecAvc() bool {
	return data.Codec == video.CodecAvc || data.Codec == video.CodecAvc3
}

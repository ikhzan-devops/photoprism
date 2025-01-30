package video

import (
	"fmt"
	"strings"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/header"
)

// ContentType returns a normalized video content type strings based on the video file type and codec.
func ContentType(mediaType, fileType, videoCodec string) string {
	if mediaType == "" && fileType == "" && videoCodec == "" {
		return header.ContentTypeBinary
	}

	if mediaType == "" {
		var c string

		if videoCodec != "" {
			c = Codecs[videoCodec]
		} else {
			c = fileType
		}

		switch c {
		case "mov", "mp4":
			mediaType = header.ContentTypeMP4
		case CodecAVC, "avc":
			mediaType = header.ContentTypeAVC // Advanced Video Coding (AVC), also known as H.264
		case CodecHEVC, "hvc", "hevc":
			mediaType = header.ContentTypeHEVC // High Efficiency Video Coding (HEVC), also known as H.265
		case CodecHEV1, "hev":
			mediaType = header.ContentTypeHEV1 // High Efficiency Video Coding (HEVC) Bitstream
		case CodecVVC, "vvc":
			mediaType = header.ContentTypeVVC // Versatile Video Coding (VVC), also known as H.266
		case CodecEVC, "evc":
			mediaType = header.ContentTypeEVC // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
		case CodecVP8, "vp08":
			mediaType = header.ContentTypeVP8
		case CodecVP9, "vp9":
			mediaType = header.ContentTypeVP9
		case CodecAV1, "av1":
			mediaType = header.ContentTypeAV1
		case CodecOGV, "ogg":
			mediaType = header.ContentTypeOGV
		case CodecWebM:
			mediaType = header.ContentTypeWebM
		}
	}

	// Add codec parameter, if possible.
	if mediaType != "" && !strings.Contains(mediaType, ";") {
		if codec, found := Codecs[videoCodec]; found && codec != "" {
			mediaType = fmt.Sprintf("%s; codecs=\"%s\"", mediaType, codec)
		}
	}

	return clean.ContentType(mediaType)
}

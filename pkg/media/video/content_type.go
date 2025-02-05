package video

import (
	"fmt"
	"mime"
	"strings"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/media/http/header"
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
			mediaType = header.ContentTypeMp4
		case CodecAvc3:
			mediaType = header.ContentTypeMp4Avc3Main // AVC (H.264) Bitstream, High Profile
		case string(fs.VideoAvc), CodecAvc:
			mediaType = header.ContentTypeMp4AvcMain // Advanced Video Coding (H.264), High Profile
		case string(fs.VideoHvc), CodecHvc:
			mediaType = header.ContentTypeMp4HvcMain // High Efficiency Video Coding (H.265)
		case string(fs.VideoHev), CodecHev:
			mediaType = header.ContentTypeMp4HevMain // High Efficiency Video Coding (HEVC) Bitstream
		case CodecVvc, "vvc":
			mediaType = header.ContentTypeMp4Vvc // Versatile Video Coding (VVC), also known as H.266
		case CodecEvc, "evc":
			mediaType = header.ContentTypeMp4Evc // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
		case CodecVp8, "vp08":
			mediaType = header.ContentTypeWebmVp8
		case CodecVp9, "vp9":
			mediaType = header.ContentTypeWebmVp9
		case CodecAv1, "av1":
			mediaType = header.ContentTypeWebmAv1
		case CodecTheora, "ogg":
			mediaType = header.ContentTypeOggTheora
		case string(fs.VideoWebm):
			mediaType = header.ContentTypeWebm
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

// Compatible tests if the video content types are expected to be compatible,
func Compatible(contentType1, contentType2 string) bool {
	// Content is likely compatible if the content type strings match exactly (case-insensitive).
	if contentType1 == "" || contentType2 == "" {
		return false
	} else if strings.EqualFold(contentType1, contentType2) {
		return true
	}

	// Sanitize and normalize content type strings.
	contentType1 = clean.ContentType(contentType1)
	contentType2 = clean.ContentType(contentType2)

	// Parse content type strings.
	mediaType1, params1, err1 := mime.ParseMediaType(contentType1)
	mediaType2, params2, err2 := mime.ParseMediaType(contentType2)

	// If parsing fails, assume the content is invalid or incompatible.
	if err1 != nil || err2 != nil {
		return false
	} else if len(params1) == 0 && len(params2) == 0 {
		return strings.EqualFold(mediaType1, mediaType2)
	}

	// If the media types don't match, assume the content is incompatible.
	if !strings.EqualFold(mediaType1, mediaType2) {
		return false
	}

	// Compare the media codecs.
	codec1 := params1["codecs"]
	codec2 := params2["codecs"]

	// Content is likely compatible if the full codec details match (case-insensitive).
	if strings.EqualFold(codec1, codec2) {
		return true
	}

	// Compare main codec names.
	codec1, _, _ = strings.Cut(codec1, ",")
	codec2, _, _ = strings.Cut(codec2, ",")

	codecName1, _, _ := strings.Cut(strings.TrimSpace(codec1), ".")
	codecName2, _, _ := strings.Cut(strings.TrimSpace(codec2), ".")

	// Content is likely compatible if the name of the main codec matches (case-insensitive).
	return strings.EqualFold(codecName1, codecName2)
}

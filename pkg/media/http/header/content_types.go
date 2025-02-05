package header

/*
	Standard content types for use in HTTP headers and the web interface.

	Browser support can be tested on one or more of the following sites:
	- https://ott.dolby.com/codec_test/index.html
	- https://dmnsgn.github.io/media-codecs/
	- https://cconcolato.github.io/media-mime-support/
    - https://cconcolato.github.io/media-mime-support/mediacapabilities.html
	- https://thorium.rocks/misc/h265-tester.html
    - https://developers.google.com/cast/docs/media
    - https://privacycheck.sec.lrz.de/active/fp_cpt/fp_can_play_type.html
	- https://chromium.googlesource.com/chromium/src.git/+/62.0.3178.1/content/browser/media/media_canplaytype_browsertest.cc
*/

// Standard ContentType identifiers for audio and video files.
const (
	ContentTypeMp4            = "video/mp4"
	ContentTypeMp4Avc         = ContentTypeMp4 + "; codecs=\"avc1\""             // MPEG-4 AVC (H.264)
	ContentTypeMp4Avc3        = ContentTypeMp4 + "; codecs=\"avc3\""             // MPEG-4 AVC Bitstream
	ContentTypeMp4AvcBaseline = ContentTypeMp4 + "; codecs=\"avc1.420028\""      // MPEG-4 AVC (H.264), Baseline Level 4.0
	ContentTypeMp4AvcMain     = ContentTypeMp4 + "; codecs=\"avc1.4d0028\""      // MPEG-4 AVC (H.264), Main Level 4.0
	ContentTypeMp4AvcHigh     = ContentTypeMp4 + "; codecs=\"avc1.640028\""      // MPEG-4 AVC (H.264), High Level 4.0
	ContentTypeMp4Avc3Main    = ContentTypeMp4 + "; codecs=\"avc3.4d0028\""      // MPEG-4 AVC Bitstream, Main Profile, may not be supported on macOS
	ContentTypeMp4Avc3High    = ContentTypeMp4 + "; codecs=\"avc3.640028\""      // MPEG-4 AVC Bitstream, High Profile, may not be supported on macOS
	ContentTypeMp4Hvc         = ContentTypeMp4 + "; codecs=\"hvc1\""             // MPEG-4 HEVC (H.265)
	ContentTypeMp4HvcMain     = ContentTypeMp4 + "; codecs=\"hvc1.1.6.L93.B0\""  // MPEG-4 HEVC (H.265), Main Profile
	ContentTypeMp4HvcMain10   = ContentTypeMp4 + "; codecs=\"hvc1.2.4.L153.B0\"" // MPEG-4 HEVC (H.265), Main 10 Profile
	ContentTypeMp4Hev         = ContentTypeMp4 + "; codecs=\"hev1\""             // MPEG-4 HEVC Bitstream
	ContentTypeMp4HevMain     = ContentTypeMp4 + "; codecs=\"hev1.1.6.L93.B0\""  // MPEG-4 HEVC Bitstream, Main Profile, not supported on macOS
	ContentTypeMp4HevMain10   = ContentTypeMp4 + "; codecs=\"hev1.2.4.L153.B0\"" // MPEG-4 HEVC Bitstream, Main 10 Profile, not supported on macOS
	ContentTypeMp4Vvc         = ContentTypeMp4 + "; codecs=\"vvc1\""             // Versatile Video Coding (VVC), also known as H.266
	ContentTypeMp4Evc         = ContentTypeMp4 + "; codecs=\"evc1\""             // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
	ContentTypeMov            = "video/quicktime"
	ContentTypeMovAvc         = ContentTypeMov + "; codecs=\"avc1\""        // Apple QuickTime AVC
	ContentTypeMovAvcMain     = ContentTypeMov + "; codecs=\"avc1.4d0028\"" // Apple QuickTime AVC, Main Level 4.0
	ContentTypeMovAvcHigh     = ContentTypeMov + "; codecs=\"avc1.640028\"" // Apple QuickTime AVC, High Level 4.0
	ContentTypeOgg            = "video/ogg"
	ContentTypeOggVorbis      = ContentTypeOgg + "; codecs=\"vorbis\""
	ContentTypeOggTheora      = ContentTypeOgg + "; codecs=\"theora, vorbis\""
	ContentTypeWebm           = "video/webm"
	ContentTypeWebmVp8        = ContentTypeWebm + "; codecs=\"vp8\""
	ContentTypeWebmVp9        = ContentTypeWebm + "; codecs=\"vp09.00.10.08\""
	ContentTypeWebmAv1        = ContentTypeWebm + "; codecs=\"av01.2.10M.10\""
)

// Standard ContentType identifiers for images and vector graphics.
const (
	ContentTypePng    = "image/png"
	ContentTypeAPng   = "image/vnd.mozilla.apng"
	ContentTypeJpeg   = "image/jpeg"
	ContentTypeJpegXL = "image/jxl"
	ContentTypeGif    = "image/gif"
	ContentTypeBmp    = "image/bmp"
	ContentTypeTiff   = "image/tiff"
	ContentTypeDng    = "image/dng"
	ContentTypeAvif   = "image/avif"
	ContentTypeAvifS  = "image/avif-sequence"
	ContentTypeHeic   = "image/heic"
	ContentTypeHeicS  = "image/heic-sequence"
	ContentTypeWebp   = "image/webp"
	ContentTypeAI     = "application/vnd.adobe.illustrator"
	ContentTypePS     = "application/postscript"
	ContentTypeEPS    = "image/eps"
	ContentTypeSVG    = "image/svg+xml"
)

// Standard ContentType identifiers for markup and sidecar files.
const (
	ContentTypeBinary    = "application/octet-stream"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeMultipart = "multipart/form-data"
	ContentTypeJson      = "application/json"
	ContentTypeJsonUtf8  = "application/json; charset=utf-8"
	ContentTypeXml       = "text/xml"
	ContentTypeHtml      = "text/html; charset=utf-8"
	ContentTypeText      = "text/plain; charset=utf-8"
	ContentTypePDF       = "application/pdf"
)

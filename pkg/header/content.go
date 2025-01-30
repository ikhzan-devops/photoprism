package header

// Standard content request and response header names.
const (
	Accept             = "Accept"
	AcceptEncoding     = "Accept-Encoding"
	AcceptLanguage     = "Accept-Language"
	AcceptRanges       = "Accept-Ranges"
	ContentType        = "Content-Type"
	ContentDisposition = "Content-Disposition"
	ContentEncoding    = "Content-Encoding"
	ContentRange       = "Content-Range"
	Location           = "Location"
	Origin             = "Origin"
	Vary               = "Vary"
)

// Standard ContentType header values.
//
// Browser support can be tested by visiting one of the following sites:
// - https://ott.dolby.com/codec_test/index.html
// - https://dmnsgn.github.io/media-codecs/
// - https://cconcolato.github.io/media-mime-support/
// - https://thorium.rocks/misc/h265-tester.html
const (
	ContentTypeBinary    = "application/octet-stream"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeMultipart = "multipart/form-data"
	ContentTypeJson      = "application/json"
	ContentTypeJsonUtf8  = "application/json; charset=utf-8"
	ContentTypeHtml      = "text/html; charset=utf-8"
	ContentTypeText      = "text/plain; charset=utf-8"
	ContentTypePDF       = "application/pdf"
	ContentTypePNG       = "image/png"
	ContentTypeJPEG      = "image/jpeg"
	ContentTypeSVG       = "image/svg+xml"
	ContentTypeQT        = "video/quicktime"
	ContentTypeMP4       = "video/mp4"
	ContentTypeAVC32     = ContentTypeMP4 + "; codecs=\"avc1.640020\""      // MPEG-4 AVC, High Profile Level 3.2
	ContentTypeAVC       = ContentTypeMP4 + "; codecs=\"avc1.640028\""      // MPEG-4 AVC, High Profile Level 4.0
	ContentTypeHEVC      = ContentTypeMP4 + "; codecs=\"hvc1.2.4.L120.B0\"" // HEVC MP4 Main10 Profile, Main Tier, Level 4.0
	ContentTypeHEV1      = ContentTypeMP4 + "; codecs=\"hev1.2.4.L120.B0\"" // HEVC Bitstream, not supported on macOS
	ContentTypeVVC       = ContentTypeMP4 + "; codecs=\"vvc1\""             // Versatile Video Coding (VVC), also known as H.266
	ContentTypeEVC       = ContentTypeMP4 + "; codecs=\"evc1\""             // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
	ContentTypeOGV       = "video/ogg"
	ContentTypeWebM      = "video/webm"
	ContentTypeVP8       = "video/webm; codecs=\"vp8\""
	ContentTypeVP9       = "video/webm; codecs=\"vp09.00.10.08\""
	ContentTypeAV1       = "video/webm; codecs=\"av01.2.10M.10\""
)

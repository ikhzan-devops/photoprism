package video

type Codec = string

// Standard video Codec types.
//
// Browser support can be tested by visiting one of the following sites:
// - https://ott.dolby.com/codec_test/index.html
// - https://dmnsgn.github.io/media-codecs/
// - https://cconcolato.github.io/media-mime-support/
// - https://thorium.rocks/misc/h265-tester.html
const (
	CodecUnknown Codec = ""
	CodecAVC     Codec = "avc1" // Advanced Video Coding (AVC), also known as H.264
	CodecHEVC    Codec = "hvc1" // High Efficiency Video Coding (HEVC), also known as H.265
	CodecHEV1    Codec = "hev1" // HEVC bitstream with the parameter sets stored in the samples, not supported on macOS
	CodecVVC     Codec = "vvc1" // Versatile Video Coding (VVC), also known as H.266
	CodecEVC     Codec = "evc1" // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
	CodecAV1     Codec = "av01" // AOMedia Video 1 (AV1)
	CodecVP8     Codec = "vp8"  // Google VP8
	CodecVP9     Codec = "vp09" // Google VP9
	CodecOGV     Codec = "ogv"  // Ogg Vorbis Video
	CodecWebM    Codec = "webm" // Google WebM
)

// Codecs maps supported string identifiers to standard Codec types.
var Codecs = StandardCodecs{
	"":                CodecUnknown,
	"a_opus":          CodecUnknown,
	"a_vorbis":        CodecUnknown,
	"avc":             CodecAVC,
	"avc1":            CodecAVC,
	"avc3":            CodecAVC,
	"v_avc":           CodecAVC,
	"v_avc1":          CodecAVC,
	"iso/avc":         CodecAVC,
	"v_mpeg4/avc":     CodecAVC,
	"v_mpeg4/iso/avc": CodecAVC,
	"hevc":            CodecHEVC,
	"hevC":            CodecHEVC,
	"hvc":             CodecHEVC,
	"hvc1":            CodecHEVC,
	"v_hvc":           CodecHEVC,
	"v_hvc1":          CodecHEVC,
	"hvcC":            CodecHEVC,
	"hvcc":            CodecHEVC,
	"hev":             CodecHEV1,
	"hev1":            CodecHEV1,
	"evc":             CodecEVC,
	"evc1":            CodecEVC,
	"evcC":            CodecEVC,
	"evcc":            CodecEVC,
	"v_evc":           CodecEVC,
	"v_evc1":          CodecEVC,
	"vvc":             CodecVVC,
	"vvcC":            CodecVVC,
	"vvcc":            CodecVVC,
	"vvc1":            CodecVVC,
	"v_vvc":           CodecVVC,
	"v_vvc1":          CodecVVC,
	"av1f":            CodecAV1,
	"av1m":            CodecAV1,
	"av1M":            CodecAV1,
	"av1s":            CodecAV1,
	"av1c":            CodecAV1,
	"av1C":            CodecAV1,
	"av1":             CodecAV1,
	"av01":            CodecAV1,
	"v_av1":           CodecAV1,
	"v_av01":          CodecAV1,
	"vp8":             CodecVP8,
	"vp08":            CodecVP8,
	"vp80":            CodecVP8,
	"v_vp8":           CodecVP8,
	"vp9":             CodecVP9,
	"vp09":            CodecVP9,
	"vp90":            CodecVP9,
	"v_vp9":           CodecVP9,
	"ogv":             CodecOGV,
	"webm":            CodecWebM,
}

// StandardCodecs maps strings to codec types.
type StandardCodecs map[string]Codec

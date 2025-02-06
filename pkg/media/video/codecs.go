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
	CodecAvc     Codec = "avc1" // Advanced Video Coding (AVC), also known as H.264
	CodecAvc3    Codec = "avc3" // AVC bitstream with the parameter sets stored in the samples, not supported on macOS
	CodecHvc     Codec = "hvc1" // High Efficiency Video Coding (HEVC), also known as H.265
	CodecHev     Codec = "hev1" // HEVC bitstream with the parameter sets stored in the samples, not supported on macOS
	CodecVvc     Codec = "vvc1" // Versatile Video Coding (VVC), also known as H.266
	CodecEvc     Codec = "evc1" // MPEG-5 Essential Video Coding (EVC), also known as ISO/IEC 23094-1
	CodecAv1     Codec = "av01" // AOMedia Video 1 (AV1)
	CodecVp8     Codec = "vp8"  // Google VP8
	CodecVp9     Codec = "vp09" // Google VP9
	CodecTheora  Codec = "ogv"  // Ogg Vorbis Video
	CodecWebm    Codec = "webm" // Google WebM
)

// Codecs maps supported string identifiers to standard Codec types.
var Codecs = StandardCodecs{
	"":                CodecUnknown,
	"a_opus":          CodecUnknown,
	"a_vorbis":        CodecUnknown,
	"avc":             CodecAvc,
	CodecAvc:          CodecAvc,
	"v_avc":           CodecAvc,
	"v_avc1":          CodecAvc,
	"iso/avc":         CodecAvc,
	"v_mpeg4/avc":     CodecAvc,
	"v_mpeg4/iso/avc": CodecAvc,
	CodecAvc3:         CodecAvc3,
	"v_avc3":          CodecAvc3,
	"hevc":            CodecHvc,
	"hevC":            CodecHvc,
	"hvc":             CodecHvc,
	CodecHvc:          CodecHvc,
	"v_hvc":           CodecHvc,
	"v_hvc1":          CodecHvc,
	"hvcC":            CodecHvc,
	"hvcc":            CodecHvc,
	"hev":             CodecHev,
	CodecHev:          CodecHev,
	"evc":             CodecEvc,
	CodecEvc:          CodecEvc,
	"evcC":            CodecEvc,
	"evcc":            CodecEvc,
	"v_evc":           CodecEvc,
	"v_evc1":          CodecEvc,
	"vvc":             CodecVvc,
	"vvcC":            CodecVvc,
	"vvcc":            CodecVvc,
	CodecVvc:          CodecVvc,
	"v_vvc":           CodecVvc,
	"v_vvc1":          CodecVvc,
	"av1f":            CodecAv1,
	"av1m":            CodecAv1,
	"av1M":            CodecAv1,
	"av1s":            CodecAv1,
	"av1c":            CodecAv1,
	"av1C":            CodecAv1,
	"av1":             CodecAv1,
	CodecAv1:          CodecAv1,
	"v_av1":           CodecAv1,
	"v_av01":          CodecAv1,
	CodecVp8:          CodecVp8,
	"vp08":            CodecVp8,
	"vp80":            CodecVp8,
	"v_vp8":           CodecVp8,
	"vp9":             CodecVp9,
	CodecVp9:          CodecVp9,
	"vp90":            CodecVp9,
	"v_vp9":           CodecVp9,
	CodecTheora:       CodecTheora,
	CodecWebm:         CodecWebm,
}

// StandardCodecs maps strings to codec types.
type StandardCodecs map[string]Codec

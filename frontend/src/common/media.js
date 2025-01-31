// Media types.
export const Animated = "animated";
export const Audio = "audio";
export const Document = "document";
export const Image = "image";
export const Raw = "raw";
export const Sidecar = "sidecar";
export const Live = "live";
export const Vector = "vector";
export const Video = "video";

// Video codec names.
//
// Browser support can be tested by visiting one of the following sites:
// - https://ott.dolby.com/codec_test/index.html
// - https://dmnsgn.github.io/media-codecs/
// - https://cconcolato.github.io/media-mime-support/
// - https://thorium.rocks/misc/h265-tester.html
export const CodecAVC = "avc1";
export const CodecAVC3 = "avc3";
export const CodecHEVC = "hvc1";
export const CodecHEV1 = "hev1";
export const CodecVVC = "vvc1";
export const CodecEVC = "evc1";
export const CodecOGV = "ogv";
export const CodecVP8 = "vp8";
export const CodecVP9 = "vp09";
export const CodecAV1 = "av01";
export const CodecAV1C = "av1c";

// Media file formats.
export const FormatMP4 = "mp4";
export const FormatAVC = "avc";
export const FormatHEVC = "hevc";
export const FormatHEV1 = "hev1";
export const FormatVVC = "vvc";
export const FormatEVC = "evc";
export const FormatOGG = "ogg";
export const FormatWebM = "webm";
export const FormatVP8 = "vp8";
export const FormatVP9 = "vp9";
export const FormatAV1 = "av1";
export const FormatWebP = "webp";
export const FormatJPEG = "jpg";
export const FormatPNG = "png";
export const FormatSVG = "svg";
export const FormatGIF = "gif";

// HTTP Content types (MIME).
export const ContentTypeMP4 = "video/mp4";
export const ContentTypeAVC = 'video/mp4; codecs="avc1.640028"'; // AVC High Profile Level 4
export const ContentTypeHEVC = 'video/mp4; codecs="hvc1.2.4.L120.B0"';
export const ContentTypeHEV1 = 'video/mp4; codecs="hev1.2.4.L120.B0"';
export const ContentTypeVVC = 'video/mp4; codecs="vvc1"';
export const ContentTypeOGV = "video/ogg";
export const ContentTypeWebM = "video/webm";
export const ContentTypeVP8 = 'video/webm; codecs="vp8"';
export const ContentTypeVP9 = 'video/webm; codecs="vp09.00.10.08"';
export const ContentTypeAV1 = 'video/webm; codecs="av01.2.10M.10"';

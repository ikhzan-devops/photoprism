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
export const CodecAvc = "avc1";
export const CodecAvc3 = "avc3";
export const CodecHvc = "hvc1";
export const CodecHev = "hev1";
export const CodecVvc = "vvc1";
export const CodecEvc = "evc1";
export const CodecTheora = "ogv";
export const CodecVp8 = "vp8";
export const CodecVp9 = "vp09";
export const CodecAv1 = "av01";
export const CodecAv1C = "av1c";

// Media file formats.
export const FormatMp4 = "mp4";
export const FormatAvc = "avc";
export const FormatHvc = "hvc";
export const FormatHev = "hev";
export const FormatVvc = "vvc";
export const FormatEvc = "evc";
export const FormatTheora = "ogg";
export const FormatWebm = "webm";
export const FormatVp8 = "vp8";
export const FormatVp9 = "vp9";
export const FormatAv1 = "av1";
export const FormatWebp = "webp";
export const FormatJpeg = "jpg";
export const FormatJpegXL = "jxl";
export const FormatPng = "png";
export const FormatGif = "gif";
export const FormatSVG = "svg";

// HTTP Content types (MIME).
export const ContentTypeMp4 = "video/mp4";
export const ContentTypeMp4AvcMain = ContentTypeMp4 + '; codecs="avc1.4d0028"'; // AVC High Profile Level 4
export const ContentTypeMp4HvcMain = ContentTypeMp4 + '; codecs="hvc1.1.6.L93.B0"';
export const ContentTypeMp4HevMain = ContentTypeMp4 + '; codecs="hev1.1.6.L93.B0"';
export const ContentTypeMp4Vvc = ContentTypeMp4 + '; codecs="vvc1"';
export const ContentTypeOgg = "video/ogg";
export const ContentTypeOggTheora = ContentTypeOgg + '; codecs="theora, vorbis"';
export const ContentTypeWebm = "video/webm";
export const ContentTypeWebmVp8 = ContentTypeWebm + '; codecs="vp8"';
export const ContentTypeWebmVp9 = ContentTypeWebm + '; codecs="vp09.00.10.08"';
export const ContentTypeWebmAv1 = ContentTypeWebm + '; codecs="av01.2.10M.10"';

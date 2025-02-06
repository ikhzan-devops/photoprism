package fs

import (
	"path/filepath"
	"strings"
)

// FileExtensions maps file extensions to standard formats
type FileExtensions map[string]Type

// Extensions contains the filename extensions of file formats known to PhotoPrism.
var Extensions = FileExtensions{
	ExtJpeg:     ImageJpeg, // .jpg
	".jpeg":     ImageJpeg,
	".jpe":      ImageJpeg,
	".jif":      ImageJpeg,
	".jfif":     ImageJpeg,
	".jfi":      ImageJpeg,
	".jxl":      ImageJpegXL,
	ExtThm:      ImageThumb,
	".tif":      ImageTiff,
	".tiff":     ImageTiff,
	".psd":      ImagePsd,
	ExtPng:      ImagePng, // .png
	".apng":     ImagePng,
	".pnga":     ImagePng,
	".pn":       ImagePng,
	".gif":      ImageGif,
	".bmp":      ImageBmp,
	ExtDng:      ImageDng, // .dng
	".avif":     ImageAvif,
	".avis":     ImageAvifS,
	".avifs":    ImageAvifS,
	".hif":      ImageHeic,
	".heif":     ImageHeic,
	".heic":     ImageHeic,
	".avci":     ImageHeic,
	".avcs":     ImageHeic,
	".heifs":    ImageHeicS,
	".heics":    ImageHeicS,
	".webp":     ImageWebp,
	".mpo":      ImageMPO,
	".3fr":      ImageRaw,
	".ari":      ImageRaw,
	".arw":      ImageRaw,
	".bay":      ImageRaw,
	".cap":      ImageRaw,
	".crw":      ImageRaw,
	".cr2":      ImageRaw,
	".cr3":      ImageRaw,
	".data":     ImageRaw,
	".dcs":      ImageRaw,
	".dcr":      ImageRaw,
	".drf":      ImageRaw,
	".eip":      ImageRaw,
	".erf":      ImageRaw,
	".fff":      ImageRaw,
	".gpr":      ImageRaw,
	".iiq":      ImageRaw,
	".k25":      ImageRaw,
	".kdc":      ImageRaw,
	".mdc":      ImageRaw,
	".mef":      ImageRaw,
	".mos":      ImageRaw,
	".mrw":      ImageRaw,
	".nef":      ImageRaw,
	".nrw":      ImageRaw,
	".obm":      ImageRaw,
	".orf":      ImageRaw,
	".pef":      ImageRaw,
	".ptx":      ImageRaw,
	".pxn":      ImageRaw,
	".r3d":      ImageRaw,
	".raf":      ImageRaw,
	".raw":      ImageRaw,
	".rwl":      ImageRaw,
	".rwz":      ImageRaw,
	".rw2":      ImageRaw,
	".srf":      ImageRaw,
	".srw":      ImageRaw,
	".sr2":      ImageRaw,
	".x3f":      ImageRaw,
	ExtAvc:      VideoAvc, // .avc
	ExtHevc:     VideoHvc, // .hevc
	ExtHev1:     VideoHev, // .hev1
	ExtVvc:      VideoVvc, // .vvc
	ExtEvc:      VideoEvc, // .evc
	".mov":      VideoMov,
	".qt":       VideoMov,
	".avi":      VideoAVI,
	".av1":      VideoAv1,
	".mpg":      VideoMpeg,
	".mpeg":     VideoMpeg,
	".mjpg":     VideoMjpeg,
	".mjpeg":    VideoMjpeg,
	".mp2":      VideoMp2,
	".mpv":      VideoMp2,
	".mp":       VideoMp4,
	ExtMp4:      VideoMp4, // .mp4
	".m4v":      VideoM4V,
	".mxf":      VideoMXF,
	".3gp":      Video3GP,
	".3g2":      Video3G2,
	".flv":      VideoFlash,
	".f4v":      VideoFlash,
	".mkv":      VideoMkv,
	".mts":      VideoAvcHD,
	".m2ts":     VideoBDAV,
	".ogv":      VideoTheora,
	".ogg":      VideoTheora,
	".ogx":      VideoTheora,
	".webm":     VideoWebm,
	".asf":      VideoASF,
	".wmv":      VideoWMV,
	".dv":       VideoDV,
	".svg":      VectorSVG,
	".ai":       VectorAI,
	".ps":       VectorPS,
	".ps2":      VectorPS,
	".ps3":      VectorPS,
	".eps":      VectorEPS,
	".eps2":     VectorEPS,
	".eps3":     VectorEPS,
	".epi":      VectorEPS,
	".ept":      VectorEPS,
	".epsf":     VectorEPS,
	".epsi":     VectorEPS,
	ExtXMP:      SidecarXMP,
	".aae":      SidecarAppleXml,
	ExtXml:      SidecarXml,
	ExtYaml:     SidecarYaml, // .yml
	".yaml":     SidecarYaml,
	ExtJson:     SidecarJson,
	ExtTxt:      SidecarText,
	".nfo":      SidecarInfo,
	ExtMd:       SidecarMarkdown,
	".markdown": SidecarMarkdown,
}

// Known tests if the file extension is known (supported).
func (m FileExtensions) Known(name string) bool {
	if name == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(name))

	if ext == "" {
		return false
	}

	if _, ok := m[ext]; ok {
		return true
	}

	return false
}

// Types returns known extensions by file type.
func (m FileExtensions) Types(noUppercase bool) TypesExt {
	result := make(TypesExt)

	if noUppercase {
		for ext, t := range m {
			if _, ok := result[t]; ok {
				result[t] = append(result[t], ext)
			} else {
				result[t] = []string{ext}
			}
		}
	} else {
		for ext, t := range m {
			extUpper := strings.ToUpper(ext)
			if _, ok := result[t]; ok {
				result[t] = append(result[t], ext, extUpper)
			} else {
				result[t] = []string{ext, extUpper}
			}
		}
	}

	return result
}

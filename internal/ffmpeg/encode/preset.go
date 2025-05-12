package encode

// FFmpeg encoding preset names from fastest to slowest,
// see https://trac.ffmpeg.org/wiki/Encode/H.264#Preset.
const (
	PresetUltraFast = "ultrafast"
	PresetSuperFast = "superfast"
	PresetVeryFast  = "veryfast"
	PresetFaster    = "faster"
	PresetFast      = "fast"
	PresetMedium    = "medium"
	PresetSlow      = "slow"
	PresetSlower    = "slower"
	PresetVerySlow  = "veryslow"
)

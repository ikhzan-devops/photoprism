package video

// Types maps identifiers to standards.
var Types = Standards{
	"":       Avc,
	"mp4":    Mp4,
	"mpeg4":  Mp4,
	"avc":    Avc,
	"avc1":   Avc,
	"avc3":   Avc, // H.264 with the parameter sets stored in the samples instead of the sample descriptions.
	"hevc":   Hvc,
	"hevC":   Hvc,
	"hvc":    Hvc,
	"hvc1":   Hvc,
	"v_hvc":  Hvc,
	"v_hvc1": Hvc,
	"hev":    Hev, // H.265 with the parameter sets stored in the samples instead of the sample descriptions.
	"hev1":   Hev,
	"evc":    Evc,
	"evc1":   Evc,
	"evcC":   Evc,
	"vvc":    Vvc,
	"vvc1":   Vvc,
	"vvcC":   Vvc,
	"vp8":    Vp8,
	"vp08":   Vp8,
	"vp80":   Vp8,
	"vp9":    Vp9,
	"vp09":   Vp9,
	"vp90":   Vp9,
	"av1":    Av1,
	"av01":   Av1,
	"ogg":    Theora,
	"ogv":    Theora,
	"webm":   Webm,
}

// Standards maps names to standardized formats.
type Standards map[string]Type

package media

// Constants representing standard media types.
const (
	Unknown   Type = ""
	Archive   Type = "archive"
	Sidecar   Type = "sidecar"
	Image     Type = "image"
	Video     Type = "video"
	Animated  Type = "animated"
	Audio     Type = "audio"
	Document  Type = "document"
	Raw       Type = "raw"
	Vector    Type = "vector"
	Live      Type = "live"
	Restoring Type = "restorng"
)

// PriorityMain specifies the minimum priority for main media types,
// like Animated, Audio, Document, Image, Live, Raw, Vector, and Video.
const PriorityMain = 4

// Priorities maps media types to integer values that represent their relative importance.
type Priorities map[Type]int

// Priority assigns a relative priority value to the media type constants defined above.
var Priority = Priorities{
	Unknown:  0,
	Sidecar:  1,
	Archive:  2,
	Image:    PriorityMain,
	Video:    8,
	Animated: 16,
	Audio:    16,
	Document: 16,
	Raw:      32,
	Vector:   32,
	Live:     64,
}

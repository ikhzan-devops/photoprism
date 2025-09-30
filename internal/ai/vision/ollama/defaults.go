package ollama

import "github.com/photoprism/photoprism/internal/ai/vision/schema"

const (
	CaptionPrompt = "Create a caption with exactly one sentence in the active voice that describes the main visual content. Begin with the main subject and clear action. Avoid text formatting, meta-language, and filler words."
	CaptionModel  = "gemma3"
	LabelSystem   = "You are a PhotoPrism vision model. Output concise JSON that matches the schema."
	LabelPrompt   = "Analyze the image and return label objects with name, confidence (0-1), and topicality (0-1)."
	Resolution    = 720
)

func LabelSchema() string {
	return schema.LabelDefaultV1
}

package openai

import "github.com/photoprism/photoprism/internal/ai/vision/schema"

const (
	DefaultModel = "gpt-5-mini"
	Resolution   = 720
)

func LabelSchema() string {
	return schema.LabelDefaultV1
}

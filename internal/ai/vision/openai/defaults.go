package openai

import "github.com/photoprism/photoprism/internal/ai/vision/schema"

var (
	DefaultModel      = "gpt-5-mini"
	DefaultResolution = 720
)

func LabelsSchema() string {
	return schema.LabelsDefaultV1
}

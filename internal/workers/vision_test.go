package workers

import (
	"testing"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/entity"
)

func TestCaptionSourceFromModel(t *testing.T) {
	if got := captionSourceFromModel(nil); got != entity.SrcImage {
		t.Fatalf("expected SrcImage for nil model, got %s", got)
	}

	openAIModel := &vision.Model{
		Service: vision.Service{RequestFormat: vision.ApiFormatOpenAI},
	}

	if got := captionSourceFromModel(openAIModel); got != entity.SrcOpenAI {
		t.Fatalf("expected SrcOpenAI for openai model, got %s", got)
	}

	providerModel := &vision.Model{Provider: "ollama"}
	if got := captionSourceFromModel(providerModel); got != entity.SrcOllama {
		t.Fatalf("expected SrcOllama from provider, got %s", got)
	}

	fallbackModel := &vision.Model{}
	if got := captionSourceFromModel(fallbackModel); got != entity.SrcImage {
		t.Fatalf("expected SrcImage fallback, got %s", got)
	}
}

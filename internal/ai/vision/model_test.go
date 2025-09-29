package vision

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/ai/vision/ollama"
	"github.com/photoprism/photoprism/internal/ai/vision/schema"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/service/http/scheme"
)

func TestModel(t *testing.T) {
	t.Run("Nasnet", func(t *testing.T) {
		ServiceUri = "https://app.localssl.dev/api/v1/vision"
		uri, method := NasnetModel.Endpoint()
		ServiceUri = ""
		assert.Equal(t, "https://app.localssl.dev/api/v1/vision/labels", uri)
		assert.Equal(t, http.MethodPost, method)

		uri, method = NasnetModel.Endpoint()
		assert.Equal(t, "", uri)
		assert.Equal(t, "", method)
	})
	t.Run("Caption", func(t *testing.T) {
		uri, method := CaptionModel.Endpoint()
		assert.Equal(t, "http://ollama:11434/api/generate", uri)
		assert.Equal(t, "POST", method)
		assert.Equal(t, ollama.ProviderName, CaptionModel.Provider)
		assert.Equal(t, ApiFormatOllama, CaptionModel.Service.RequestFormat)
		assert.Equal(t, ApiFormatOllama, CaptionModel.Service.ResponseFormat)
		assert.Equal(t, string(scheme.Base64), CaptionModel.Service.FileScheme)
		assert.Equal(t, ollama.ProviderName, CaptionModel.ProviderName())

		model, name, version := CaptionModel.Model()
		assert.Equal(t, "gemma3:latest", model)
		assert.Equal(t, "gemma3", name)
		assert.Equal(t, "latest", version)
	})
	t.Run("ParseName", func(t *testing.T) {
		m := &Model{
			Type:       ModelTypeCaption,
			Name:       "deepseek-r1:1.5b",
			Version:    "",
			Resolution: 720,
			Prompt:     ollama.CaptionPrompt,
			Service: Service{
				Uri:            "http://foo:bar@photoprism-vision:5000/api/v1/vision/caption",
				FileScheme:     scheme.Data,
				RequestFormat:  ApiFormatVision,
				ResponseFormat: ApiFormatVision,
			},
		}

		uri, method := m.Endpoint()
		assert.Equal(t, "http://foo:bar@photoprism-vision:5000/api/v1/vision/caption", uri)
		assert.Equal(t, "POST", method)

		model, name, version := m.Model()
		assert.Equal(t, "deepseek-r1:1.5b", model)
		assert.Equal(t, "deepseek-r1", name)
		assert.Equal(t, "1.5b", version)
	})
}

func TestParseTypes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		result := ParseTypes("nsfw, labels, Caption, generate")
		assert.Equal(t, ModelTypes{"nsfw", "labels", "caption", "generate"}, result)
	})
	t.Run("None", func(t *testing.T) {
		result := ParseTypes("")
		assert.Equal(t, ModelTypes{}, result)
	})
	t.Run("Invalid", func(t *testing.T) {
		result := ParseTypes("foo, captions")
		assert.Equal(t, ModelTypes{}, result)
	})
}

func TestModelFormatAndSchema(t *testing.T) {
	t.Run("DefaultOllamaFormat", func(t *testing.T) {
		m := &Model{
			Type: ModelTypeLabels,
			Service: Service{
				RequestFormat:  ApiFormatOllama,
				ResponseFormat: ApiFormatOllama,
			},
		}

		assert.Equal(t, FormatJSON, m.GetFormat())
	})
	t.Run("InlineSchema", func(t *testing.T) {
		s := "{\n  \"labels\": []\n}"
		m := &Model{Schema: s}

		assert.Equal(t, s, m.SchemaTemplate())
		assert.Contains(t, m.SchemaInstructions(), "Return JSON")
	})
	t.Run("SchemaFileAndEnv", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "schema.json")
		content := "{\n  \"labels\": [{\"name\": \"test\"}]\n}"
		assert.NoError(t, os.WriteFile(filePath, []byte(content), fs.ModeConfigFile))

		m := &Model{
			Type:       ModelTypeLabels,
			SchemaFile: filePath,
		}

		// First read should use file content.
		assert.Equal(t, content, m.SchemaTemplate())

		// Reset and use env override with a different file.
		otherFile := filepath.Join(tempDir, "schema-override.json")
		otherContent := "{\n  \"labels\": []\n,  \"markers\": []\n}"
		assert.NoError(t, os.WriteFile(otherFile, []byte(otherContent), fs.ModeConfigFile))

		t.Setenv(labelSchemaEnvVar, otherFile)

		m2 := &Model{Type: ModelTypeLabels}
		assert.Equal(t, otherContent, m2.SchemaTemplate())
	})
	t.Run("DefaultLabelSchema", func(t *testing.T) {
		m := &Model{Type: ModelTypeLabels}
		assert.Equal(t, strings.TrimSpace(schema.LabelDefaultV1), m.SchemaTemplate())
		assert.Contains(t, m.SchemaInstructions(), "Return JSON")
	})
	t.Run("FormatOverride", func(t *testing.T) {
		m := &Model{Format: "JSON"}
		assert.Equal(t, FormatJSON, m.GetFormat())
	})
	t.Run("DefaultLabelPrompts", func(t *testing.T) {
		m := &Model{Type: ModelTypeLabels,
			Service: Service{RequestFormat: ApiFormatOllama, ResponseFormat: ApiFormatOllama},
		}
		assert.Equal(t, ollama.LabelPrompt, m.GetPrompt())
		assert.Equal(t, ollama.LabelSystem, m.GetSystemPrompt())
	})
	t.Run("ProviderDefaults", func(t *testing.T) {
		m := &Model{
			Type:     ModelTypeCaption,
			Provider: ollama.ProviderName,
			Service: Service{
				Uri: "http://localhost:11434/api/generate",
			},
		}

		m.ApplyProviderDefaults()

		assert.Equal(t, ApiFormatOllama, m.Service.RequestFormat)
		assert.Equal(t, ApiFormatOllama, m.Service.ResponseFormat)
		assert.Equal(t, string(scheme.Base64), m.Service.FileScheme)
		assert.Equal(t, ollama.ProviderName, m.ProviderName())
	})
	t.Run("ProviderInference", func(t *testing.T) {
		m := &Model{
			Type: ModelTypeGenerate,
			Service: Service{
				Uri:           "https://example.com/api/vision",
				RequestFormat: ApiFormatVision,
			},
		}

		assert.Equal(t, ProviderVision, m.ProviderName())

		m2 := &Model{
			Type: ModelTypeLabels,
			Service: Service{
				Uri:           "https://example.com/api/ollama",
				RequestFormat: ApiFormatOllama,
			},
		}

		assert.Equal(t, ollama.ProviderName, m2.ProviderName())

		assert.Equal(t, ProviderTensorFlow, NasnetModel.ProviderName())

		m3 := &Model{}
		assert.Equal(t, ProviderLocal, m3.ProviderName())
	})
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/ai/vision"
)

func TestConfig_DisableFrontend(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.False(t, c.DisableFrontend())
}

func TestConfig_DisableSettings(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.False(t, c.DisableSettings())
}

func TestConfig_DisableWebDAV(t *testing.T) {
	c := NewConfig(CliTestContext())

	c.options.Public = false
	c.options.ReadOnly = false
	c.options.Demo = false

	assert.False(t, c.DisableWebDAV())

	c.options.Public = true
	c.options.ReadOnly = false
	c.options.Demo = false

	assert.True(t, c.DisableWebDAV())

	c.options.Public = false
	c.options.ReadOnly = true
	c.options.Demo = false

	assert.False(t, c.DisableWebDAV())

	c.options.Public = false
	c.options.ReadOnly = false
	c.options.Demo = true

	assert.True(t, c.DisableWebDAV())

	c.options.Public = true
	c.options.ReadOnly = true
	c.options.Demo = true

	assert.True(t, c.DisableWebDAV())

	c.options.Public = false
	c.options.ReadOnly = false
	c.options.Demo = false

	assert.False(t, c.DisableWebDAV())
}

func TestConfig_DisableExifTool(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.False(t, c.DisableExifTool())

	c.options.ExifToolBin = "XXX"
	assert.True(t, c.DisableExifTool())
}

func TestConfig_ExifToolEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.ExifToolEnabled())

	c.options.ExifToolBin = "XXX"
	assert.False(t, c.ExifToolEnabled())
}

func TestConfig_DisableFaces(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.False(t, c.DisableFaces())
	c.options.DisableFaces = true
	assert.True(t, c.DisableFaces())
	c.options.DisableFaces = false
	c.options.DisableTensorFlow = true
	assert.True(t, c.DisableFaces())
	c.options.DisableTensorFlow = false
	assert.False(t, c.DisableFaces())
}

func TestConfig_DisableClassification(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.False(t, c.DisableClassification())
	c.options.DisableClassification = true
	assert.True(t, c.DisableClassification())
	c.options.DisableClassification = false
	c.options.DisableTensorFlow = true
	assert.True(t, c.DisableClassification())
	c.options.DisableTensorFlow = false
	assert.False(t, c.DisableClassification())
}

func TestConfig_GenerateLabelsWhileIndexing(t *testing.T) {
	t.Run("ClassificationDisabled", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		c.options.DisableClassification = true
		withVisionConfig(t, vision.NewConfig())
		if c.GenerateLabelsWhileIndexing() {
			t.Fatalf("expected labels to be skipped when classification disabled")
		}
	})
	t.Run("NilVisionConfig", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		withVisionConfig(t, nil)
		if c.GenerateLabelsWhileIndexing() {
			t.Fatalf("expected labels to be skipped without vision config")
		}
	})
	t.Run("DefaultModelUsesInlineGeneration", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		defaultModel := cloneVisionModel(vision.NasnetModel)
		withVisionConfig(t, &vision.ConfigValues{Models: vision.Models{defaultModel}})
		if !c.GenerateLabelsWhileIndexing() {
			t.Fatalf("expected default model to run during indexing")
		}
	})
	t.Run("CustomModelDefersGeneration", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		defaultModel := cloneVisionModel(vision.NasnetModel)
		custom := &vision.Model{Type: vision.ModelTypeLabels, Name: "custom", Provider: "ollama"}
		withVisionConfig(t, &vision.ConfigValues{Models: vision.Models{defaultModel, custom}})
		if c.GenerateLabelsWhileIndexing() {
			t.Fatalf("expected custom labels model to defer indexing generation")
		}
	})
}

func TestConfig_GenerateLabelsAfterIndexing(t *testing.T) {
	t.Run("DefaultModel", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		defaultModel := cloneVisionModel(vision.NasnetModel)
		withVisionConfig(t, &vision.ConfigValues{Models: vision.Models{defaultModel}})
		if c.GenerateLabelsAfterIndexing() {
			t.Fatalf("expected default model to skip post-index generation")
		}
	})
	t.Run("CustomModel", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		defaultModel := cloneVisionModel(vision.NasnetModel)
		custom := &vision.Model{Type: vision.ModelTypeLabels, Name: "custom", Provider: "ollama"}
		withVisionConfig(t, &vision.ConfigValues{Models: vision.Models{defaultModel, custom}})
		if !c.GenerateLabelsAfterIndexing() {
			t.Fatalf("expected custom model to run after indexing")
		}
	})
	t.Run("NilVisionConfig", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		withVisionConfig(t, nil)
		if c.GenerateLabelsAfterIndexing() {
			t.Fatalf("expected nil vision config to skip label generation")
		}
	})
}

func TestConfig_GenerateCaptionsAfterIndexing(t *testing.T) {
	t.Run("ClassificationDisabled", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		c.options.DisableClassification = true
		withVisionConfig(t, vision.NewConfig())
		if c.GenerateCaptionsAfterIndexing() {
			t.Fatalf("expected captions to be skipped when classification disabled")
		}
	})
	t.Run("DefaultCaptionModel", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		withVisionConfig(t, vision.NewConfig())
		if !c.GenerateCaptionsAfterIndexing() {
			t.Fatalf("expected default caption model to run after indexing")
		}
	})
	t.Run("ExplicitDefaultCaption", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		caption := cloneVisionModel(vision.CaptionModel)
		caption.Default = true
		withVisionConfig(t, &vision.ConfigValues{Models: vision.Models{caption}})
		if c.GenerateCaptionsAfterIndexing() {
			t.Fatalf("expected explicit default caption model to skip post-index generation")
		}
	})
	t.Run("NilVisionConfig", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		withVisionConfig(t, nil)
		if c.GenerateCaptionsAfterIndexing() {
			t.Fatalf("expected nil vision config to skip caption generation")
		}
	})
}

func TestConfig_DisableDarktable(t *testing.T) {
	c := NewConfig(CliTestContext())
	missing := c.DarktableBin() == ""

	assert.Equal(t, missing, c.DisableDarktable())
	c.options.DisableRaw = true
	assert.True(t, c.DisableDarktable())
	c.options.DisableRaw = false
	assert.Equal(t, missing, c.DisableDarktable())
	c.options.DisableDarktable = true
	assert.True(t, c.DisableDarktable())
	c.options.DisableDarktable = false
	assert.Equal(t, missing, c.DisableDarktable())
}

func TestConfig_DisableRawTherapee(t *testing.T) {
	c := NewConfig(CliTestContext())
	missing := c.RawTherapeeBin() == ""

	assert.Equal(t, missing, c.DisableRawTherapee())
	c.options.DisableRaw = true
	assert.True(t, c.DisableRawTherapee())
	c.options.DisableRaw = false
	assert.Equal(t, missing, c.DisableRawTherapee())
	c.options.DisableRawTherapee = true
	assert.True(t, c.DisableRawTherapee())
	c.options.DisableRawTherapee = false
	assert.Equal(t, missing, c.DisableRawTherapee())
}

func TestConfig_DisableImageMagick(t *testing.T) {
	c := NewConfig(CliTestContext())
	missing := c.ImageMagickBin() == ""

	assert.Equal(t, missing, c.DisableImageMagick())
	c.options.DisableRaw = true
	assert.Equal(t, missing, c.DisableImageMagick())
	c.options.DisableRaw = false
	assert.Equal(t, missing, c.DisableImageMagick())
	c.options.DisableImageMagick = true
	assert.True(t, c.DisableImageMagick())
	c.options.DisableImageMagick = false
	assert.Equal(t, missing, c.DisableImageMagick())
}

func TestConfig_DisableVips(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, false, c.DisableVips())
	c.options.DisableVips = true
	assert.True(t, c.DisableVips())
	c.options.DisableVips = false
	assert.Equal(t, false, c.DisableVips())
}

func TestConfig_DisableSips(t *testing.T) {
	c := NewConfig(CliTestContext())
	missing := c.SipsBin() == ""

	assert.Equal(t, missing, c.DisableSips())
	c.options.DisableSips = true
	assert.True(t, c.DisableSips())
	c.options.DisableSips = false
	assert.Equal(t, missing, c.DisableSips())
}

func TestConfig_DisableVector(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, c.Sponsor(), !c.DisableVectors())
	c.options.DisableVectors = true
	assert.True(t, c.DisableVectors())
	c.options.DisableVectors = false
	assert.Equal(t, c.Sponsor(), !c.DisableVectors())
}

func TestConfig_DisableRsvgConvert(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, c.Sponsor(), !c.DisableRsvgConvert())
	c.options.DisableVectors = true
	assert.True(t, c.DisableRsvgConvert())
	c.options.DisableVectors = false
	assert.Equal(t, c.Sponsor(), !c.DisableVectors())
}

func TestConfig_DisableRaw(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.False(t, c.DisableRaw())
	c.options.DisableRaw = true
	assert.True(t, c.DisableRaw())
	assert.True(t, c.DisableDarktable())
	assert.True(t, c.DisableRawTherapee())
	c.options.DisableRaw = false
	assert.False(t, c.DisableRaw())
	c.options.DisableDarktable = true
	c.options.DisableRawTherapee = true
	assert.False(t, c.DisableRaw())
	c.options.DisableDarktable = false
	c.options.DisableRawTherapee = false
	assert.False(t, c.DisableRaw())
	assert.False(t, c.DisableDarktable())
	assert.False(t, c.DisableRawTherapee())
}

func withVisionConfig(t *testing.T, cfg *vision.ConfigValues) {
	t.Helper()
	prev := vision.Config
	vision.Config = cfg
	t.Cleanup(func() {
		vision.Config = prev
	})
}

func cloneVisionModel(m *vision.Model) *vision.Model {
	copy := *m
	return &copy
}

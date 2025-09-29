package vision

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
)

func TestOptions(t *testing.T) {
	var configPath = fs.Abs("testdata")
	var configFile = filepath.Join(configPath, "vision.yml")

	t.Run("Save", func(t *testing.T) {
		_ = os.Remove(configFile)
		options := NewConfig()
		err := options.Save(configFile)
		assert.NoError(t, err)
		err = options.Load(configFile)
		assert.NoError(t, err)
	})
	t.Run("LoadMissingFile", func(t *testing.T) {
		options := NewConfig()
		err := options.Load(filepath.Join(configPath, "invalid.yml"))
		assert.Error(t, err)
	})
}

func TestConfigModelPrefersLastEnabled(t *testing.T) {
	defaultModel := *NasnetModel
	defaultModel.Disabled = false
	defaultModel.Name = "nasnet-default"

	customModel := &Model{
		Type:     ModelTypeLabels,
		Name:     "ollama-labels",
		Provider: "ollama",
		Disabled: false,
	}

	cfg := &ConfigValues{
		Models: Models{
			&defaultModel,
			customModel,
		},
	}

	got := cfg.Model(ModelTypeLabels)
	if got != customModel {
		t.Fatalf("expected last enabled model, got %v", got)
	}

	customModel.Disabled = true
	got = cfg.Model(ModelTypeLabels)
	if got == nil || got.Name != defaultModel.Name {
		t.Fatalf("expected fallback to default model, got %v", got)
	}
}

func TestConfigValues_IsDefaultAndIsCustom(t *testing.T) {
	nasnetModel := *NasnetModel
	defaultModel := &nasnetModel
	defaultModel.Default = false

	t.Run("DefaultModel", func(t *testing.T) {
		cfg := &ConfigValues{Models: Models{defaultModel}}
		if !cfg.IsDefault(ModelTypeLabels) {
			t.Fatalf("expected default model to be reported as default")
		}
		if cfg.IsCustom(ModelTypeLabels) {
			t.Fatalf("expected default model not to be reported as custom")
		}
	})
	t.Run("CustomOverridesDefault", func(t *testing.T) {
		custom := &Model{Type: ModelTypeLabels, Name: "custom", Provider: "ollama"}
		cfg := &ConfigValues{Models: Models{defaultModel, custom}}
		if cfg.IsDefault(ModelTypeLabels) {
			t.Fatalf("expected custom model to disable default detection")
		}
		if !cfg.IsCustom(ModelTypeLabels) {
			t.Fatalf("expected custom model to be detected")
		}
	})
	t.Run("DisabledCustomFallsBackToDefault", func(t *testing.T) {
		custom := &Model{Type: ModelTypeLabels, Name: "custom", Provider: "ollama", Disabled: true}
		cfg := &ConfigValues{Models: Models{defaultModel, custom}}
		if !cfg.IsDefault(ModelTypeLabels) {
			t.Fatalf("expected disabled custom model to fall back to default")
		}
		if cfg.IsCustom(ModelTypeLabels) {
			t.Fatalf("expected disabled custom model not to force custom mode")
		}
	})
	t.Run("MissingModel", func(t *testing.T) {
		cfg := &ConfigValues{}
		if cfg.IsDefault(ModelTypeLabels) {
			t.Fatalf("expected missing model to return false for default detection")
		}
		if cfg.IsCustom(ModelTypeLabels) {
			t.Fatalf("expected missing model to return false for custom detection")
		}
	})
}

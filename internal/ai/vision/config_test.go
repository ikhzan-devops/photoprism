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

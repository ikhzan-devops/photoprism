package photoprism

import (
	"testing"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/config"
)

func TestIndex_Start(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cfg := config.TestConfig()
	initErr := cfg.InitializeTestData()
	assert.NoError(t, initErr)

	convert := NewConvert(cfg)
	ind := NewIndex(cfg, convert, NewFiles(), NewPhotos())
	imp := NewImport(cfg, ind, convert)
	opt := ImportOptionsMove(cfg.ImportPath(), "")

	imp.Start(opt)

	indexOpt := IndexOptionsAll()
	indexOpt.Rescan = false

	found, updated := ind.Start(indexOpt)
	assert.GreaterOrEqual(t, len(found), 0)
	assert.GreaterOrEqual(t, updated, 0)

	t.Logf("index run 1: found %s", english.Plural(updated, "file", "files"))
	t.Logf("index run 1: updated %s", english.Plural(updated, "file", "files"))

	time.Sleep(time.Second)

	found, updated = ind.Start(indexOpt)
	assert.GreaterOrEqual(t, len(found), 0)
	assert.GreaterOrEqual(t, updated, 0)

	t.Logf("index run 2: found %s", english.Plural(updated, "file", "files"))
	t.Logf("index run 2: updated %s", english.Plural(updated, "file", "files"))

	time.Sleep(time.Second)

	found, updated = ind.Start(indexOpt)
	assert.GreaterOrEqual(t, len(found), 0)
	assert.GreaterOrEqual(t, updated, 0)

	t.Logf("index run 3: found %s", english.Plural(updated, "file", "files"))
	t.Logf("index run 3: updated %s", english.Plural(updated, "file", "files"))
}

func TestNewIndexFindLabelsUsesVisionModelShouldRun(t *testing.T) {
	prevVision := vision.Config
	vision.Config = vision.NewConfig()
	t.Cleanup(func() {
		vision.Config = prevVision
	})

	cfg := config.NewConfig(config.CliTestContext())
	ind := NewIndex(cfg, NewConvert(cfg), NewFiles(), NewPhotos())

	if ind == nil {
		t.Fatalf("expected index instance")
	}

	if !ind.findLabels {
		t.Fatalf("expected labels to be generated for default configuration")
	}

	vision.Config = &vision.ConfigValues{Models: vision.Models{
		&vision.Model{Type: vision.ModelTypeLabels, Run: vision.RunManual},
	}}

	cfgManual := config.NewConfig(config.CliTestContext())
	indManual := NewIndex(cfgManual, NewConvert(cfgManual), NewFiles(), NewPhotos())

	if indManual.findLabels {
		t.Fatalf("expected labels to be skipped when vision config disallows run")
	}
}

func TestIndex_File(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cfg := config.TestConfig()
	initErr := cfg.InitializeTestData()
	assert.NoError(t, initErr)

	convert := NewConvert(cfg)
	ind := NewIndex(cfg, convert, NewFiles(), NewPhotos())

	err := ind.FileName("xxx", IndexOptionsAll())

	assert.Equal(t, IndexFailed, err.Status)
}

// TestIndexConfigureFaceDetectionFacesOnlyManual ensures faces-only runs override manual scheduling.
func TestIndexConfigureFaceDetectionFacesOnlyManual(t *testing.T) {
	cfg := config.NewConfig(config.CliTestContext())
	cfg.Options().FaceEngineRun = string(vision.RunManual)

	ind := NewIndex(cfg, nil, nil, nil)
	require.NotNil(t, ind)
	require.False(t, ind.findFaces)

	opt := NewIndexOptions("", true, false, true, true, true)
	ind.configureFaceDetection(opt)

	require.True(t, ind.findFaces)
}

// TestIndexConfigureFaceDetectionFacesOnlyNever confirms the scheduler honors the "never" run mode.
func TestIndexConfigureFaceDetectionFacesOnlyNever(t *testing.T) {
	cfg := config.NewConfig(config.CliTestContext())
	cfg.Options().FaceEngineRun = string(vision.RunNever)

	ind := NewIndex(cfg, nil, nil, nil)
	require.NotNil(t, ind)

	opt := NewIndexOptions("", true, false, true, true, true)
	ind.configureFaceDetection(opt)

	require.False(t, ind.findFaces)
}

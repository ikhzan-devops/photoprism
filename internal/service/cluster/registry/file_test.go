package registry

import (
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/config"
)

// TestFindByNameDeterministic verifies that FindByName returns the most
// recently updated node when multiple registry entries share the same Name.
func TestFindByNameDeterministic(t *testing.T) {
	// Isolate storage/config to avoid interference from other tests.
	tmp := t.TempDir()
	t.Setenv("PHOTOPRISM_STORAGE_PATH", tmp)

	conf := config.NewTestConfig("cluster-registry-findbyname")

	r, err := NewFileRegistry(conf)
	assert.NoError(t, err)

	// Two nodes with the same name but different UpdatedAt timestamps.
	old := Node{
		ID:        "id-old",
		Name:      "pp-node-01",
		Role:      "instance",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}
	newer := Node{
		ID:        "id-new",
		Name:      "pp-node-01",
		Role:      "instance",
		CreatedAt: "2024-02-01T00:00:00Z",
		UpdatedAt: "2024-02-01T00:00:00Z",
	}

	// Write YAML files directly to avoid timing flakiness.
	b1, err := yaml.Marshal(old)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(r.fileName(old.ID), b1, 0o600))

	b2, err := yaml.Marshal(newer)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(r.fileName(newer.ID), b2, 0o600))

	// Expect the most recently updated node (id-new).
	got, err := r.FindByName("pp-node-01")
	assert.NoError(t, err)
	if assert.NotNil(t, got) {
		assert.Equal(t, "id-new", got.ID)
		assert.Equal(t, "pp-node-01", got.Name)
	}

	// Non-existent name should return os.ErrNotExist.
	_, err = r.FindByName("does-not-exist")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

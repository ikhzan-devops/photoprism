package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/fs"
)

func TestConfig_Cluster(t *testing.T) {
	t.Run("Flags", func(t *testing.T) {
		c := NewConfig(CliTestContext())

		// Defaults
		assert.False(t, c.ClusterPortal())
		assert.False(t, c.ClusterNode())

		// Toggle values
		c.options.ClusterPortal = true
		c.options.ClusterNode = true
		assert.True(t, c.ClusterPortal())
		assert.True(t, c.ClusterNode())
	})

	t.Run("Paths", func(t *testing.T) {
		c := NewConfig(CliTestContext())

		// Use an isolated config path so we don't affect repo storage fixtures.
		tempCfg := t.TempDir()
		c.options.ConfigPath = tempCfg

		// ClusterConfigPath always points to a "cluster" subfolder under ConfigPath.
		expectedCluster := filepath.Join(c.ConfigPath(), fs.ClusterDir)
		assert.Equal(t, expectedCluster, c.ClusterConfigPath())

		// ClusterThemePath falls back to ThemePath if cluster dir does not exist.
		expectedTheme := filepath.Join(c.ConfigPath(), fs.ThemeDir)
		assert.Equal(t, expectedTheme, c.ClusterThemePath())

		// When only the cluster directory exists (without a theme subfolder), it still falls back to ThemePath.
		assert.NoError(t, os.MkdirAll(expectedCluster, 0o755))
		assert.Equal(t, expectedTheme, c.ClusterThemePath())

		// When the cluster theme directory exists, ClusterThemePath returns it.
		expectedClusterTheme := filepath.Join(expectedCluster, fs.ThemeDir)
		assert.NoError(t, os.MkdirAll(expectedClusterTheme, 0o755))
		assert.Equal(t, expectedClusterTheme, c.ClusterThemePath())
	})

	t.Run("PortalAndSecrets", func(t *testing.T) {
		c := NewConfig(CliTestContext())

		// Defaults
		assert.Equal(t, "", c.PortalUrl())
		assert.Equal(t, "", c.PortalClient())
		assert.Equal(t, "", c.PortalSecret())
		assert.Equal(t, "", c.NodeSecret())

		// Set and read back values
		c.options.PortalUrl = "https://portal.example.test"
		c.options.PortalClient = "client-id"
		c.options.PortalSecret = "client-secret"
		c.options.NodeSecret = "node-secret"

		assert.Equal(t, "https://portal.example.test", c.PortalUrl())
		assert.Equal(t, "client-id", c.PortalClient())
		assert.Equal(t, "client-secret", c.PortalSecret())
		assert.Equal(t, "node-secret", c.NodeSecret())
	})

	t.Run("AbsolutePaths", func(t *testing.T) {
		c := NewConfig(CliTestContext())
		tempCfg := t.TempDir()
		c.options.ConfigPath = tempCfg

		// ThemePath should be absolute.
		assert.True(t, filepath.IsAbs(c.ThemePath()))

		// ClusterThemePath should be absolute (fallback case).
		assert.True(t, filepath.IsAbs(c.ClusterThemePath()))

		// Create cluster theme directory and verify again.
		clusterTheme := filepath.Join(c.ClusterConfigPath(), fs.ThemeDir)
		assert.NoError(t, os.MkdirAll(clusterTheme, 0o755))
		assert.True(t, filepath.IsAbs(c.ClusterThemePath()))
	})
}

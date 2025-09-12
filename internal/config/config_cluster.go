package config

import (
	"path/filepath"

	"github.com/photoprism/photoprism/pkg/fs"
)

// NodeName returns the human-readable instance name, if specified.
func (c *Config) NodeName() string {
	return c.options.NodeName
}

// NodeSecret returns the node's authentication secret, if specified.
func (c *Config) NodeSecret() string {
	return c.options.NodeSecret
}

// PortalUrl returns the URL of the cluster portal server, if configured.
func (c *Config) PortalUrl() string {
	return c.options.PortalUrl
}

// PortalClient returns the portal client ID, if configured.
func (c *Config) PortalClient() string {
	return c.options.PortalClient
}

// PortalSecret returns the portal client secret, if configured.
func (c *Config) PortalSecret() string {
	return c.options.PortalSecret
}

// ClusterPortal returns true if this instance should act as a cluster portal.
func (c *Config) ClusterPortal() bool {
	return c.options.ClusterPortal
}

// ClusterNode returns true if this instance should be configured as a cluster node.
func (c *Config) ClusterNode() bool {
	return c.options.ClusterNode
}

// ClusterConfigPath returns the path to portal config files.
func (c *Config) ClusterConfigPath() string {
	return filepath.Join(c.ConfigPath(), fs.ClusterDir)
}

// ClusterThemePath returns the path to the shared theme files.
func (c *Config) ClusterThemePath() string {
	// Prefer the cluster-specific theme directory if it exists.
	if dir := filepath.Join(c.ClusterConfigPath(), fs.ThemeDir); fs.PathExists(dir) {
		return dir
	}

	// Fallback to the default theme directory in the main config path.
	return c.ThemePath()
}

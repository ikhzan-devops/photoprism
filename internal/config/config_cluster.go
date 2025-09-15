package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/photoprism/photoprism/pkg/service/cluster"
)

// NodeName returns the unique name of this node within the cluster (lowercase letters and numbers only).
func (c *Config) NodeName() string {
	return clean.TypeLowerDash(c.options.NodeName)
}

// NodeType returns the type of this node for cluster operation (portal, instance, service).
func (c *Config) NodeType() string {
	switch c.options.NodeType {
	case cluster.Portal, cluster.Instance, cluster.Service:
		return c.options.NodeType
	default:
		return cluster.Instance
	}
}

// NodeSecret returns the private node key for intra-cluster communication.
func (c *Config) NodeSecret() string {
	if c.options.NodeSecret != "" {
		return c.options.NodeSecret
	} else if fileName := FlagFilePath("NODE_SECRET"); fileName == "" {
		return ""
	} else if b, err := os.ReadFile(fileName); err != nil || len(b) == 0 {
		log.Warnf("config: failed to read node secret from %s (%s)", fileName, err)
		return ""
	} else {
		return string(b)
	}
}

// PortalUrl returns the URL of the cluster portal server, if configured.
func (c *Config) PortalUrl() string {
	return c.options.PortalUrl
}

// PortalToken returns the token required to access the portal API endpoints.
func (c *Config) PortalToken() string {
	if c.options.PortalToken != "" {
		return c.options.PortalToken
	} else if fileName := FlagFilePath("PORTAL_TOKEN"); fileName == "" {
		return ""
	} else if b, err := os.ReadFile(fileName); err != nil || len(b) == 0 {
		log.Warnf("config: failed to read portal token from %s (%s)", fileName, err)
		return ""
	} else {
		return string(b)
	}
}

// ClusterPortal returns true if this instance should act as a cluster portal.
func (c *Config) ClusterPortal() bool {
	return c.IsPortal()
}

// IsPortal returns true if the configured node type is "portal".
func (c *Config) IsPortal() bool {
	return c.NodeType() == cluster.Portal
}

// PortalConfigPath returns the path to the default configuration for cluster nodes.
func (c *Config) PortalConfigPath() string {
	return filepath.Join(c.ConfigPath(), fs.ClusterDir)
}

// PortalThemePath returns the path to the theme files for cluster nodes to use.
func (c *Config) PortalThemePath() string {
	// Prefer the cluster-specific theme directory if it exists.
	if dir := filepath.Join(c.PortalConfigPath(), fs.ThemeDir); fs.PathExists(dir) {
		return dir
	}

	// Fallback to the default theme directory in the main config path.
	return c.ThemePath()
}

// PortalUUID returns a stable UUIDv4 that uniquely identifies the Portal.
// Precedence: env PHOTOPRISM_PORTAL_UUID -> options.yml (PortalUUID) -> auto-generate and persist.
func (c *Config) PortalUUID() string {
	// Use value loaded into options only if it is persisted in the current options.yml.
	// This avoids tests (or defaults) loading a UUID from an unrelated file path.
	if c.options.PortalUUID != "" {
		// Respect explicit CLI value if provided.
		if c.cliCtx != nil && c.cliCtx.IsSet("portal-uuid") {
			return c.options.PortalUUID
		}
		// Otherwise, only trust a persisted value from the current options.yml.
		if fs.FileExists(c.OptionsYaml()) {
			return c.options.PortalUUID
		}
	}

	// Generate, persist, and cache in memory if still empty.
	id := rnd.UUID()
	c.options.PortalUUID = id

	if err := c.savePortalUUID(id); err != nil {
		log.Warnf("config: failed to persist PortalUUID to %s (%s)", c.OptionsYaml(), err)
	}

	return id
}

// savePortalUUID writes or updates the PortalUUID key in options.yml without
// touching unrelated keys. Creates the file and directories if needed.
func (c *Config) savePortalUUID(id string) error {
	// Always resolve against the current ConfigPath and remember it explicitly
	// so subsequent calls don't accidentally point to a previous default.
	cfgDir := c.ConfigPath()
	if err := fs.MkdirAll(cfgDir); err != nil {
		return err
	}
	fileName := filepath.Join(cfgDir, "options.yml")

	var m map[string]interface{}

	if fs.FileExists(fileName) {
		if b, err := os.ReadFile(fileName); err == nil && len(b) > 0 {
			_ = yaml.Unmarshal(b, &m)
		}
	}

	if m == nil {
		m = map[string]interface{}{}
	}

	m["PortalUUID"] = id

	if b, err := yaml.Marshal(m); err != nil {
		return err
	} else if err = os.WriteFile(fileName, b, 0o644); err != nil {
		return err
	}

	// Remember options.yml path for subsequent loads and ensure in-memory options see the value.
	if c.options != nil {
		c.options.OptionsYaml = fileName
		_ = c.options.Load(fileName)
	}

	return nil
}

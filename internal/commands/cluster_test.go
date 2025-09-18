package commands

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/photoprism/get"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/fs"
)

func TestClusterSummaryCommand(t *testing.T) {
	t.Run("NotPortal", func(t *testing.T) {
		out, err := RunWithTestContext(ClusterSummaryCommand, []string{"summary"})
		assert.Error(t, err)
		_ = out
	})
}

func TestClusterNodesListCommand(t *testing.T) {
	t.Run("NotPortal", func(t *testing.T) {
		out, err := RunWithTestContext(ClusterNodesListCommand, []string{"ls"})
		assert.Error(t, err)
		_ = out
	})
}

func TestClusterNodesShowCommand(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		_ = os.Setenv("PHOTOPRISM_NODE_TYPE", "portal")
		defer os.Unsetenv("PHOTOPRISM_NODE_TYPE")
		out, err := RunWithTestContext(ClusterNodesShowCommand, []string{"show", "does-not-exist"})
		assert.Error(t, err)
		_ = out
	})
}

func TestClusterThemePullCommand(t *testing.T) {
	t.Run("NotPortal", func(t *testing.T) {
		out, err := RunWithTestContext(ClusterThemePullCommand.Subcommands[0], []string{"pull"})
		assert.Error(t, err)
		_ = out
	})
}

func TestClusterRegisterCommand(t *testing.T) {
	t.Run("ValidationMissingURL", func(t *testing.T) {
		out, err := RunWithTestContext(ClusterRegisterCommand, []string{"register", "--name", "pp-node-01", "--type", "instance", "--portal-token", "token"})
		assert.Error(t, err)
		_ = out
	})
}

func TestClusterSuccessPaths_PortalLocal(t *testing.T) {
	// Enable portal mode for local admin commands.
	c := get.Config()
	c.Options().NodeType = "portal"

	// Ensure registry and theme paths exist.
	portCfg := c.PortalConfigPath()
	nodesDir := filepath.Join(portCfg, "nodes")
	themeDir := filepath.Join(portCfg, "theme")
	assert.NoError(t, fs.MkdirAll(nodesDir))
	assert.NoError(t, fs.MkdirAll(themeDir))

	// Create a theme file to zip.
	themeFile := filepath.Join(themeDir, "test.txt")
	assert.NoError(t, os.WriteFile(themeFile, []byte("ok"), 0o600))

	// Create a registry node via FileRegistry.
	r, err := reg.NewFileRegistry(c)
	assert.NoError(t, err)
	n := &reg.Node{Name: "pp-node-01", Type: "instance", Labels: map[string]string{"env": "test"}}
	assert.NoError(t, r.Put(n))

	// nodes ls (JSON)
	out, err := RunWithTestContext(ClusterNodesListCommand, []string{"ls", "--json"})
	assert.NoError(t, err)
	assert.Contains(t, out, "pp-node-01")

	// nodes show by name
	out, err = RunWithTestContext(ClusterNodesShowCommand, []string{"show", "pp-node-01"})
	assert.NoError(t, err)
	assert.Contains(t, out, "pp-node-01")

	// nodes mod: add another label (non-interactive)
	out, err = RunWithTestContext(ClusterNodesModCommand, []string{"mod", "pp-node-01", "--label", "region=us-east-1", "-y"})
	assert.NoError(t, err)
	_ = out

	// theme pull via HTTP: fake portal endpoint returns a zip with test.txt
	// Prepare temp destination
	destDir := t.TempDir()

	// Create a fake portal theme zip server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/cluster/theme" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		// Build a small zip in-memory
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		f, _ := zw.Create("test.txt")
		_, _ = f.Write([]byte("ok"))
		_ = zw.Close()
		_, _ = w.Write(buf.Bytes())
	}))
	defer ts.Close()

	_ = os.Setenv("PHOTOPRISM_PORTAL_URL", ts.URL)
	_ = os.Setenv("PHOTOPRISM_PORTAL_TOKEN", "test-token")
	defer os.Unsetenv("PHOTOPRISM_PORTAL_URL")
	defer os.Unsetenv("PHOTOPRISM_PORTAL_TOKEN")

	out, err = RunWithTestContext(ClusterThemePullCommand.Subcommands[0], []string{"pull", "--dest", destDir, "-f", "--portal-url=" + ts.URL, "--portal-token=test-token"})
	assert.NoError(t, err)
	// Expect extracted file
	assert.FileExists(t, filepath.Join(destDir, "test.txt"))
}

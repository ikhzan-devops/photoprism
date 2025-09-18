package instance

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/service/cluster"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

var log = event.Log

func init() {
	// Register early so this can adjust DB settings before connectDb().
	config.RegisterEarly("cluster-instance", InitConfig, nil)
}

// InitConfig performs instance bootstrap: optional registration with the Portal
// and theme installation. Runs early during config.Init().
func InitConfig(c *config.Config) error {
	if !cluster.BootstrapAutoJoinEnabled && !cluster.BootstrapAutoThemeEnabled {
		return nil
	}

	// Skip on portal nodes and unknown node types.
	if c.IsPortal() || c.NodeType() != cluster.Instance {
		return nil
	}

	portalURL := strings.TrimSpace(c.PortalUrl())
	portalToken := strings.TrimSpace(c.PortalToken())
	if portalURL == "" || portalToken == "" {
		return nil
	}

	u, err := url.Parse(portalURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		log.Warnf("cluster: invalid portal url %s", clean.Log(portalURL))
		return nil
	}

	// Enforce TLS for non-local URLs.
	if u.Scheme != "https" && !isLocalHost(u.Hostname()) {
		log.Warnf("cluster: refusing non-TLS portal url %s on non-local host", clean.Log(portalURL))
		return nil
	}

	// Register with retry policy.
	if cluster.BootstrapAutoJoinEnabled {
		if err := registerWithPortal(c, u, portalToken); err != nil {
			// Registration errors are expected when the Portal is temporarily unavailable
			// or not configured with cluster endpoints (404). Keep as warn to signal
			// exhaustion/terminal errors; per-attempt details are logged at debug level.
			log.Warnf("cluster: register failed (%s)", clean.Error(err))
		}
	}

	// Pull theme if missing.
	if cluster.BootstrapAutoThemeEnabled {
		if err := installThemeIfMissing(c, u, portalToken); err != nil {
			// Theme install failures are non-critical; log at debug to avoid noise.
			log.Debugf("cluster: theme install skipped/failed (%s)", clean.Error(err))
		}
	}

	return nil
}

func isLocalHost(h string) bool {
	switch strings.ToLower(h) {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		// TODO: Consider treating RFC1918/link-local hosts as local for TLS enforcement
		// if the operator explicitly opts in (e.g., via a policy var). Keep simple for now.
		return false
	}
}

func newHTTPClient(timeout time.Duration) *http.Client {
	// TODO: Consider reusing a shared *http.Transport with sane defaults and enabling
	// proxy support explicitly if required. For now, rely on net/http defaults and
	// the HTTPS_PROXY set in config.Init().
	return &http.Client{Timeout: timeout}
}

func registerWithPortal(c *config.Config, portal *url.URL, token string) error {
	maxAttempts := cluster.BootstrapRegisterMaxAttempts
	delay := cluster.BootstrapRegisterRetryDelay
	timeout := cluster.BootstrapRegisterTimeout

	endpoint := *portal
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/api/v1/cluster/nodes/register"

	// Decide if DB rotation is desired as per spec: only if driver is MySQL/MariaDB
	// and no DSN/fields are set (raw options) and no password is provided via file.
	opts := c.Options()
	driver := c.DatabaseDriver()
	wantRotateDB := (driver == config.MySQL || driver == config.MariaDB) &&
		opts.DatabaseDsn == "" && opts.DatabaseName == "" && opts.DatabaseUser == "" && opts.DatabasePassword == "" &&
		c.DatabasePassword() == ""

	payload := map[string]interface{}{
		"nodeName":    c.NodeName(),
		"nodeType":    string(cluster.Instance), // JSON wire format is string
		"internalUrl": c.InternalUrl(),
	}
	if wantRotateDB {
		payload["rotate"] = true
	}

	bodyBytes, _ := json.Marshal(payload)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, _ := http.NewRequest(http.MethodPost, endpoint.String(), strings.NewReader(string(bodyBytes)))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := newHTTPClient(timeout).Do(req)
		if err != nil {
			if attempt < maxAttempts {
				log.Debugf("cluster: register attempt %d/%d error: %s", attempt, maxAttempts, clean.Error(err))
				time.Sleep(delay)
				continue
			}
			return err
		}

		// Ensure body is closed after handling the response.
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			var r cluster.RegisterResponse
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&r); err != nil {
				return err
			}
			if err := persistRegistration(c, &r, wantRotateDB); err != nil {
				return err
			}
			if resp.StatusCode == http.StatusCreated {
				log.Infof("cluster: registered as %s (%d)", clean.LogQuote(r.Node.Name), resp.StatusCode)
			} else {
				log.Infof("cluster: registration ok (%d)", resp.StatusCode)
			}
			return nil
		case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
			// Terminal errors (no retry). 404 likely indicates a Portal without cluster endpoints.
			return errors.New(resp.Status)
		case http.StatusTooManyRequests:
			if attempt < maxAttempts {
				log.Debugf("cluster: register attempt %d/%d rate limited", attempt, maxAttempts)
				time.Sleep(delay)
				continue
			}
			return errors.New(resp.Status)
		case http.StatusConflict, http.StatusBadRequest:
			// Do not retry on 400/409 per spec intent.
			return errors.New(resp.Status)
		default:
			if attempt < maxAttempts {
				log.Debugf("cluster: register attempt %d/%d server responded %s", attempt, maxAttempts, resp.Status)
				// TODO: Consider exponential backoff with jitter instead of constant delay.
				time.Sleep(delay)
				continue
			}
			return errors.New(resp.Status)
		}
	}
	return nil
}

func isTemporary(err error) bool {
	var nerr net.Error
	return errors.As(err, &nerr) && nerr.Timeout()
}

func persistRegistration(c *config.Config, r *cluster.RegisterResponse, wantRotateDB bool) error {
	updates := map[string]interface{}{}

	// Persist node secret only if missing locally and provided by server.
	if r.Secrets != nil && r.Secrets.NodeSecret != "" && c.NodeSecret() == "" {
		updates["NodeSecret"] = r.Secrets.NodeSecret
	}

	// Persist DB settings only if rotation was requested and driver is MySQL/MariaDB
	// and local DB not configured (as checked before calling).
	if wantRotateDB {
		if r.DB.DSN != "" {
			updates["DatabaseDriver"] = config.MySQL
			updates["DatabaseDsn"] = r.DB.DSN
		} else if r.DB.Name != "" && r.DB.User != "" && r.DB.Password != "" {
			server := r.DB.Host
			if r.DB.Port > 0 {
				server = net.JoinHostPort(r.DB.Host, strconv.Itoa(r.DB.Port))
			}
			updates["DatabaseDriver"] = config.MySQL
			updates["DatabaseServer"] = server
			updates["DatabaseName"] = r.DB.Name
			updates["DatabaseUser"] = r.DB.User
			updates["DatabasePassword"] = r.DB.Password
		}
	}

	if len(updates) == 0 {
		return nil
	}

	if err := mergeOptionsYaml(c, updates); err != nil {
		return err
	}

	// Reload into memory so later code paths see updated values during this run.
	_ = c.Options().Load(c.OptionsYaml())

	if hasDBUpdate(updates) {
		log.Infof("cluster: database settings applied; restart required to take effect")
	}
	return nil
}

func hasDBUpdate(m map[string]interface{}) bool {
	if _, ok := m["DatabaseDsn"]; ok {
		return true
	}
	if _, ok := m["DatabaseName"]; ok {
		return true
	}
	if _, ok := m["DatabaseUser"]; ok {
		return true
	}
	if _, ok := m["DatabasePassword"]; ok {
		return true
	}
	if _, ok := m["DatabaseServer"]; ok {
		return true
	}
	return false
}

func mergeOptionsYaml(c *config.Config, updates map[string]interface{}) error {
	if err := fs.MkdirAll(c.ConfigPath()); err != nil {
		return err
	}
	fileName := c.OptionsYaml()

	var m map[string]interface{}
	if fs.FileExists(fileName) {
		if b, err := os.ReadFile(fileName); err == nil && len(b) > 0 {
			_ = yaml.Unmarshal(b, &m)
		}
	}
	if m == nil {
		m = map[string]interface{}{}
	}
	for k, v := range updates {
		m[k] = v
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, b, 0o644)
}

// installThemeIfMissing downloads and installs the Portal-provided theme if the
// local theme directory is missing or lacks an app.js file.
func installThemeIfMissing(c *config.Config, portal *url.URL, token string) error {
	themeDir := c.ThemePath()
	need := !fs.PathExists(themeDir) || (cluster.BootstrapThemeInstallOnlyIfMissingJS && !fs.FileExists(filepath.Join(themeDir, "app.js")))
	if !need && !cluster.BootstrapAllowThemeOverwrite {
		return nil
	}

	endpoint := *portal
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/api/v1/cluster/theme"

	req, _ := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/zip")

	resp, err := newHTTPClient(cluster.BootstrapRegisterTimeout).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Save to temp zip.
		if err := fs.MkdirAll(c.TempPath()); err != nil {
			return err
		}
		zipName := filepath.Join(c.TempPath(), "cluster-theme.zip")
		out, err := os.Create(zipName)
		if err != nil {
			return err
		}
		if _, err = io.Copy(out, resp.Body); err != nil {
			_ = out.Close()
			return err
		}
		_ = out.Close()

		// Extract with moderate limits.
		if err := fs.MkdirAll(themeDir); err != nil {
			return err
		}
		_, _, unzipErr := fs.Unzip(zipName, themeDir, 32*fs.MB, 512*fs.MB)
		return unzipErr
	case http.StatusNotFound:
		// No theme configured at Portal.
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New(resp.Status)
	default:
		return errors.New(resp.Status)
	}
}

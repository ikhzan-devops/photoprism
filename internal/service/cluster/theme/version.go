package theme

import (
	"os"
	"path/filepath"
	"time"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// DetectVersion returns the sanitized theme version for the given directory.
// It prefers a version.txt file when present and falls back to the app.js
// modification timestamp when the version file is missing or empty.
func DetectVersion(themePath string) (string, error) {
	versionFile := filepath.Join(themePath, fs.VersionTxtFile)

	if data, err := os.ReadFile(versionFile); err == nil {
		if v := clean.TypeUnicode(string(data)); v != "" {
			return v, nil
		}
	}

	appPath := filepath.Join(themePath, fs.AppJsFile)
	info, err := os.Stat(appPath)

	if err != nil {
		return "", err
	}

	return clean.TypeUnicode(info.ModTime().UTC().Format(time.RFC3339)), nil
}

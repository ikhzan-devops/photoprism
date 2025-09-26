package node

import (
	"os"
	"testing"

	"github.com/photoprism/photoprism/pkg/fs"
)

// TestMain ensures SQLite test DB artifacts are purged after the suite runs.
func TestMain(m *testing.M) {
	// Run unit tests.
	code := m.Run()

	// Remove temporary SQLite files after running the tests.
	fs.PurgeTestDbFiles(".", false)

	os.Exit(code)
}

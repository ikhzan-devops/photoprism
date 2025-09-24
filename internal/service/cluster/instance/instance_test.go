package instance

import (
	"os"
	"testing"

	"github.com/photoprism/photoprism/pkg/fs"
)

// TestMain ensures SQLite test DB artifacts are purged after the suite runs.
func TestMain(m *testing.M) {
	code := m.Run()
	fs.PurgeTestDbFiles(".", false)
	os.Exit(code)
}

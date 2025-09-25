package wellknown

import (
	"os"
	"testing"

	"github.com/photoprism/photoprism/pkg/fs"
)

func TestMain(m *testing.M) {
	// Run unit tests.
	code := m.Run()

	// Remove temporary SQLite files after running the tests.
	fs.PurgeTestDbFiles(".", false)

	os.Exit(code)
}

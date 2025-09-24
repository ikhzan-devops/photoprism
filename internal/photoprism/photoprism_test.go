package photoprism

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/pkg/fs"
)

func TestMain(m *testing.M) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	c := config.NewTestConfig("photoprism")
	SetConfig(c)
	defer c.CloseDb()

	code := m.Run()

	// Purge local SQLite test artifacts created during this package's tests.
	fs.PurgeTestDbFiles(".", false)

	os.Exit(code)
}

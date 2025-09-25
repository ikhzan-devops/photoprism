package avatar

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/photoprism/get"
)

func TestMain(m *testing.M) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	tempDir, err := os.MkdirTemp("", "avatar-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	c := config.NewMinimalTestConfigWithDb("avatar", tempDir)
	get.SetConfig(c)
	photoprism.SetConfig(c)
	defer c.CloseDb()

	code := m.Run()

	os.Exit(code)
}

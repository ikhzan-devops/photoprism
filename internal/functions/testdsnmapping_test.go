package functions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhotoPrismTestToDriverDsn(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
		originalDsn := os.Getenv("PHOTOPRISM_TEST_DSN_SQLITE")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "sqlite")
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITE", ":memory:?cache=shared&_foreign_keys=on")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "sqlite", driver)
		assert.Equal(t, ":memory:?cache=shared&_foreign_keys=on", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITE", originalDsn)
	})

	t.Run("sqlitefile", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
		originalDsn := os.Getenv("PHOTOPRISM_TEST_DSN_SQLITE")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "sqlitefile")
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITEFILE", ".test.db")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "sqlite", driver)
		assert.Equal(t, ".test.db", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITEFILE", originalDsn)
	})

	t.Run("mariadb", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
		originalDsn := os.Getenv("PHOTOPRISM_TEST_DSN_MARIADB")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "mariadb")
		os.Setenv("PHOTOPRISM_TEST_DSN_MARIADB", ".test.db")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "mysql", driver)
		assert.Equal(t, ".test.db", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
		os.Setenv("PHOTOPRISM_TEST_DSN_MARIADB", originalDsn)
	})

	t.Run("mysql8", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
		originalDsn := os.Getenv("PHOTOPRISM_TEST_DSN_MYSQL8")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "mysql8")
		os.Setenv("PHOTOPRISM_TEST_DSN_MYSQL8", ".test.db")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "mysql", driver)
		assert.Equal(t, ".test.db", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
		os.Setenv("PHOTOPRISM_TEST_DSN_MYSQL8", originalDsn)
	})

	t.Run("postgres", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
		originalDsn := os.Getenv("PHOTOPRISM_TEST_DSN_POSTGRES")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "postgres")
		os.Setenv("PHOTOPRISM_TEST_DSN_POSTGRES", ".test.db")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "postgres", driver)
		assert.Equal(t, ".test.db", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
		os.Setenv("PHOTOPRISM_TEST_DSN_POSTGRES", originalDsn)
	})

	t.Run("default", func(t *testing.T) {
		originalDsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", "unknown")

		driver, dsn := PhotoPrismTestToDriverDsn()

		assert.Equal(t, "sqlite", driver)
		assert.Equal(t, ":memory:?cache=shared&_foreign_keys=on", dsn)

		os.Setenv("PHOTOPRISM_TEST_DSN_NAME", originalDsnName)
	})
}

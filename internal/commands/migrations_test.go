package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/event"
)

func TestMigrationCommand(t *testing.T) {
	t.Run("NoMigrateSettings", func(t *testing.T) {
		// Run command with test context.
		output, err := RunWithTestContext(MigrationsCommands, []string{"migrations", "transfer"})

		// Check command output for plausibility.
		// t.Logf(output)
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), "config: transfer config must be provided")
		}
		assert.Equal(t, "", output)
	})

	t.Run("InvalidCommand", func(t *testing.T) {
		// Run command with test context.
		output, err := RunWithTestContext(MigrationsCommands, []string{"migrations", "--magles"})

		// Check command output for plausibility.
		// t.Logf(output)
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), "flag provided but not defined: -magles")
		}
		assert.Contains(t, output, "flag provided but not defined: -magles")
	})

	t.Run("TargetPopulated", func(t *testing.T) {
		// Setup target database
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/targetpopulated.test.db")
		if err := copyFile("/go/src/github.com/photoprism/photoprism/storage/acceptance/backup.db", "/go/src/github.com/photoprism/photoprism/storage/targetpopulated.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "acceptance:acceptance@tcp(mariadb:4001)/acceptance?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
			"--transfer-driver", "sqlite",
			"--transfer-dsn", "/go/src/github.com/photoprism/photoprism/storage/targetpopulated.test.db?_busy_timeout=5000&_foreign_keys=on"}
		cmdArgs := []string{"migrations", "transfer"}

		ctx := NewTestContextWithParse(appArgs, cmdArgs)

		output, err := RunWithProvidedTestContext(ctx, MigrationsCommands, cmdArgs)

		// Check command output for plausibility.
		// t.Logf(output)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "migrate: transfer target database is not empty")
		assert.NotContains(t, output, "Usage")

		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/targetpopulated.test.db")
		}
	})

	t.Run("MySQLtoSQLite", func(t *testing.T) {
		// Remove target database file
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlite.test.db")

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "acceptance:acceptance@tcp(mariadb:4001)/acceptance?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
			"--transfer-driver", "sqlite",
			"--transfer-dsn", "/go/src/github.com/photoprism/photoprism/storage/mysqltosqlite.test.db?_busy_timeout=5000&_foreign_keys=on"}
		cmdArgs := []string{"migrations", "transfer"}

		ctx := NewTestContextWithParse(appArgs, cmdArgs)

		s := event.Subscribe("log.info")
		defer event.Unsubscribe(s)

		var l string

		assert.IsType(t, hub.Subscription{}, s)

		go func() {
			for msg := range s.Receiver {
				l += msg.Fields["message"].(string) + "\n"
			}
		}()

		output, err := RunWithProvidedTestContext(ctx, MigrationsCommands, cmdArgs)

		// Check command output for plausibility.
		//t.Logf(output)
		assert.NoError(t, err)
		assert.NotContains(t, output, "Usage")

		time.Sleep(time.Second)

		// Check command output.
		if l == "" {
			t.Fatal("log output missing")
		}
		// t.Logf(l)

		assert.Contains(t, l, "migrate: number of albums transfered")
		assert.Contains(t, l, "migrate: number of albumusers transfered")
		assert.Contains(t, l, "migrate: number of cameras transfered")
		assert.Contains(t, l, "migrate: number of categories transfered")
		assert.Contains(t, l, "migrate: number of cells transfered")
		assert.Contains(t, l, "migrate: number of clients transfered")
		assert.Contains(t, l, "migrate: number of countries transfered")
		assert.Contains(t, l, "migrate: number of duplicates transfered")
		assert.Contains(t, l, "migrate: number of errors transfered")
		assert.Contains(t, l, "migrate: number of faces transfered")
		assert.Contains(t, l, "migrate: number of files transfered")
		assert.Contains(t, l, "migrate: number of fileshares transfered")
		assert.Contains(t, l, "migrate: number of filesyncs transfered")
		assert.Contains(t, l, "migrate: number of folders transfered")
		assert.Contains(t, l, "migrate: number of keywords transfered")
		assert.Contains(t, l, "migrate: number of labels transfered")
		assert.Contains(t, l, "migrate: number of lenses transfered")
		assert.Contains(t, l, "migrate: number of links transfered")
		assert.Contains(t, l, "migrate: number of markers transfered")
		assert.Contains(t, l, "migrate: number of passcodes transfered")
		assert.Contains(t, l, "migrate: number of passwords transfered")
		assert.Contains(t, l, "migrate: number of photos transfered")
		assert.Contains(t, l, "migrate: number of photousers transfered")
		assert.Contains(t, l, "migrate: number of places transfered")
		assert.Contains(t, l, "migrate: number of reactions transfered")
		assert.Contains(t, l, "migrate: number of sessions transfered")
		assert.Contains(t, l, "migrate: number of services transfered")
		assert.Contains(t, l, "migrate: number of subjects transfered")
		assert.Contains(t, l, "migrate: number of users transfered")
		assert.Contains(t, l, "migrate: number of userdetails transfered")
		assert.Contains(t, l, "migrate: number of usersettings transfered")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlite.test.db")
		}
	})

	t.Run("MySQLtoSQLitePopulated", func(t *testing.T) {
		// Remove target database file
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db")
		if err := copyFile("/go/src/github.com/photoprism/photoprism/storage/acceptance/backup.db", "/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "acceptance:acceptance@tcp(mariadb:4001)/acceptance?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
			"--transfer-driver", "sqlite",
			"--transfer-dsn", "/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db?_busy_timeout=5000&_foreign_keys=on"}
		cmdArgs := []string{"migrations", "transfer", "-force"}

		ctx := NewTestContextWithParse(appArgs, cmdArgs)

		s := event.Subscribe("log.info")
		defer event.Unsubscribe(s)

		var l string

		assert.IsType(t, hub.Subscription{}, s)

		go func() {
			for msg := range s.Receiver {
				l += msg.Fields["message"].(string) + "\n"
			}
		}()

		output, err := RunWithProvidedTestContext(ctx, MigrationsCommands, cmdArgs)

		// Check command output for plausibility.
		//t.Logf(output)
		assert.NoError(t, err)
		assert.NotContains(t, output, "Usage")

		time.Sleep(time.Second)

		// Check command output.
		if l == "" {
			t.Fatal("log output missing")
		}
		// t.Logf(l)

		assert.Contains(t, l, "migrate: number of albums transfered")
		assert.Contains(t, l, "migrate: number of albumusers transfered")
		assert.Contains(t, l, "migrate: number of cameras transfered")
		assert.Contains(t, l, "migrate: number of categories transfered")
		assert.Contains(t, l, "migrate: number of cells transfered")
		assert.Contains(t, l, "migrate: number of clients transfered")
		assert.Contains(t, l, "migrate: number of countries transfered")
		assert.Contains(t, l, "migrate: number of duplicates transfered")
		assert.Contains(t, l, "migrate: number of errors transfered")
		assert.Contains(t, l, "migrate: number of faces transfered")
		assert.Contains(t, l, "migrate: number of files transfered")
		assert.Contains(t, l, "migrate: number of fileshares transfered")
		assert.Contains(t, l, "migrate: number of filesyncs transfered")
		assert.Contains(t, l, "migrate: number of folders transfered")
		assert.Contains(t, l, "migrate: number of keywords transfered")
		assert.Contains(t, l, "migrate: number of labels transfered")
		assert.Contains(t, l, "migrate: number of lenses transfered")
		assert.Contains(t, l, "migrate: number of links transfered")
		assert.Contains(t, l, "migrate: number of markers transfered")
		assert.Contains(t, l, "migrate: number of passcodes transfered")
		assert.Contains(t, l, "migrate: number of passwords transfered")
		assert.Contains(t, l, "migrate: number of photos transfered")
		assert.Contains(t, l, "migrate: number of photousers transfered")
		assert.Contains(t, l, "migrate: number of places transfered")
		assert.Contains(t, l, "migrate: number of reactions transfered")
		assert.Contains(t, l, "migrate: number of sessions transfered")
		assert.Contains(t, l, "migrate: number of services transfered")
		assert.Contains(t, l, "migrate: number of subjects transfered")
		assert.Contains(t, l, "migrate: number of users transfered")
		assert.Contains(t, l, "migrate: number of userdetails transfered")
		assert.Contains(t, l, "migrate: number of usersettings transfered")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db")
		}
	})
}

func copyFile(source, target string) error {
	if _, err := os.Stat(source); err != nil {
		return fmt.Errorf("copyFile: source file %s is required", source)
	}

	if _, err := os.Stat(target); err != nil {
		if err = os.Remove(target); err != nil {
			if !strings.Contains(err.Error(), "no such file or directory") {
				return fmt.Errorf("copyFile: target file %s can not be removed with error %s", target, err.Error())
			}
		}
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("copyFile: source file %s can not be opened with error %s", source, err.Error())
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("copyFile: target file %s can not be opened with error %s", target, err.Error())
	}

	defer func() {
		closeErr := targetFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("copyFile: copy failed with error %s", err.Error())
	}

	if err = targetFile.Sync(); err != nil {
		return fmt.Errorf("copyFile: target sync failed with error %s", err.Error())
	}

	return nil
}

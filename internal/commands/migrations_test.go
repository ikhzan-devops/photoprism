package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
		if err := copyFile("/go/src/github.com/photoprism/photoprism/internal/commands/testdata/transfer_sqlite3", "/go/src/github.com/photoprism/photoprism/storage/targetpopulated.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
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

		// Load migrate database as source
		if dumpName, err := filepath.Abs("./testdata/transfer_mysql"); err != nil {
			t.Fatal(err)
		} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
			"-e", "source "+dumpName).Run(); err != nil {
			t.Fatal(err)
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
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

		assert.Contains(t, l, "migrate: number of albums transfered 31")
		assert.Contains(t, l, "migrate: number of albumusers transfered 0")
		assert.Contains(t, l, "migrate: number of cameras transfered 6")
		assert.Contains(t, l, "migrate: number of categories transfered 1")
		assert.Contains(t, l, "migrate: number of cells transfered 9")
		assert.Contains(t, l, "migrate: number of clients transfered 7")
		assert.Contains(t, l, "migrate: number of countries transfered 1")
		assert.Contains(t, l, "migrate: number of duplicates transfered 0")
		assert.Contains(t, l, "migrate: number of errors transfered 0")
		assert.Contains(t, l, "migrate: number of faces transfered 7")
		assert.Contains(t, l, "migrate: number of files transfered 71")
		assert.Contains(t, l, "migrate: number of fileshares transfered 2")
		assert.Contains(t, l, "migrate: number of filesyncs transfered 3")
		assert.Contains(t, l, "migrate: number of folders transfered 3")
		assert.Contains(t, l, "migrate: number of keywords transfered 26")
		assert.Contains(t, l, "migrate: number of labels transfered 32")
		assert.Contains(t, l, "migrate: number of lenses transfered 2")
		assert.Contains(t, l, "migrate: number of links transfered 5")
		assert.Contains(t, l, "migrate: number of markers transfered 18")
		assert.Contains(t, l, "migrate: number of passcodes transfered 3")
		assert.Contains(t, l, "migrate: number of passwords transfered 11")
		assert.Contains(t, l, "migrate: number of photos transfered 58")
		assert.Contains(t, l, "migrate: number of photousers transfered 0")
		assert.Contains(t, l, "migrate: number of places transfered 10")
		assert.Contains(t, l, "migrate: number of reactions transfered 3")
		assert.Contains(t, l, "migrate: number of sessions transfered 21")
		assert.Contains(t, l, "migrate: number of services transfered 2")
		assert.Contains(t, l, "migrate: number of subjects transfered 6")
		assert.Contains(t, l, "migrate: number of users transfered 11")
		assert.Contains(t, l, "migrate: number of userdetails transfered 9")
		assert.Contains(t, l, "migrate: number of usersettings transfered 13")
		assert.Contains(t, l, "migrate: number of usershares transfered 1")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlite.test.db")
		}
	})

	t.Run("MySQLtoSQLitePopulated", func(t *testing.T) {
		// Remove target database file
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db")
		if err := copyFile("/go/src/github.com/photoprism/photoprism/internal/commands/testdata/transfer_sqlite3", "/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Load migrate database as source
		if dumpName, err := filepath.Abs("./testdata/transfer_mysql"); err != nil {
			t.Fatal(err)
		} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
			"-e", "source "+dumpName).Run(); err != nil {
			t.Fatal(err)
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "mysql",
			"--database-dsn", "migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s",
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

		assert.Contains(t, l, "migrate: number of albums transfered 31")
		assert.Contains(t, l, "migrate: number of albumusers transfered 0")
		assert.Contains(t, l, "migrate: number of cameras transfered 6")
		assert.Contains(t, l, "migrate: number of categories transfered 1")
		assert.Contains(t, l, "migrate: number of cells transfered 9")
		assert.Contains(t, l, "migrate: number of clients transfered 7")
		assert.Contains(t, l, "migrate: number of countries transfered 1")
		assert.Contains(t, l, "migrate: number of duplicates transfered 0")
		assert.Contains(t, l, "migrate: number of errors transfered 0")
		assert.Contains(t, l, "migrate: number of faces transfered 7")
		assert.Contains(t, l, "migrate: number of files transfered 71")
		assert.Contains(t, l, "migrate: number of fileshares transfered 2")
		assert.Contains(t, l, "migrate: number of filesyncs transfered 3")
		assert.Contains(t, l, "migrate: number of folders transfered 3")
		assert.Contains(t, l, "migrate: number of keywords transfered 26")
		assert.Contains(t, l, "migrate: number of labels transfered 32")
		assert.Contains(t, l, "migrate: number of lenses transfered 2")
		assert.Contains(t, l, "migrate: number of links transfered 5")
		assert.Contains(t, l, "migrate: number of markers transfered 18")
		assert.Contains(t, l, "migrate: number of passcodes transfered 3")
		assert.Contains(t, l, "migrate: number of passwords transfered 11")
		assert.Contains(t, l, "migrate: number of photos transfered 58")
		assert.Contains(t, l, "migrate: number of photousers transfered 0")
		assert.Contains(t, l, "migrate: number of places transfered 10")
		assert.Contains(t, l, "migrate: number of reactions transfered 3")
		assert.Contains(t, l, "migrate: number of sessions transfered 21")
		assert.Contains(t, l, "migrate: number of services transfered 2")
		assert.Contains(t, l, "migrate: number of subjects transfered 6")
		assert.Contains(t, l, "migrate: number of users transfered 11")
		assert.Contains(t, l, "migrate: number of userdetails transfered 9")
		assert.Contains(t, l, "migrate: number of usersettings transfered 13")
		assert.Contains(t, l, "migrate: number of usershares transfered 1")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/mysqltosqlitepopulated.test.db")
		}
	})

	t.Run("SQLiteToMySQL", func(t *testing.T) {
		// Remove target database file
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db")

		// Load migrate database as source
		if err := copyFile("/go/src/github.com/photoprism/photoprism/internal/commands/testdata/transfer_sqlite3", "/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Clear MySQL target (migrate)
		if dumpName, err := filepath.Abs("./testdata/reset-migrate.sql"); err != nil {
			t.Fatal(err)
		} else {
			resetFile, err := os.Open(dumpName)
			if err != nil {
				t.Log("unable to open reset file")
				t.Fatal(err)
			}
			defer resetFile.Close()

			cmd := exec.Command("mysql")
			cmd.Stdin = resetFile

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(output)
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "sqlite",
			"--database-dsn", "/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db?_busy_timeout=5000&_foreign_keys=on",
			"--transfer-driver", "mysql",
			"--transfer-dsn", "migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s"}
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

		assert.Contains(t, l, "migrate: number of albums transfered 31")
		assert.Contains(t, l, "migrate: number of albumusers transfered 0")
		assert.Contains(t, l, "migrate: number of cameras transfered 6")
		assert.Contains(t, l, "migrate: number of categories transfered 1")
		assert.Contains(t, l, "migrate: number of cells transfered 9")
		assert.Contains(t, l, "migrate: number of clients transfered 7")
		assert.Contains(t, l, "migrate: number of countries transfered 1")
		assert.Contains(t, l, "migrate: number of duplicates transfered 0")
		assert.Contains(t, l, "migrate: number of errors transfered 0")
		assert.Contains(t, l, "migrate: number of faces transfered 7")
		assert.Contains(t, l, "migrate: number of files transfered 71")
		assert.Contains(t, l, "migrate: number of fileshares transfered 2")
		assert.Contains(t, l, "migrate: number of filesyncs transfered 3")
		assert.Contains(t, l, "migrate: number of folders transfered 3")
		assert.Contains(t, l, "migrate: number of keywords transfered 26")
		assert.Contains(t, l, "migrate: number of labels transfered 32")
		assert.Contains(t, l, "migrate: number of lenses transfered 2")
		assert.Contains(t, l, "migrate: number of links transfered 5")
		assert.Contains(t, l, "migrate: number of markers transfered 18")
		assert.Contains(t, l, "migrate: number of passcodes transfered 3")
		assert.Contains(t, l, "migrate: number of passwords transfered 11")
		assert.Contains(t, l, "migrate: number of photos transfered 58")
		assert.Contains(t, l, "migrate: number of photousers transfered 0")
		assert.Contains(t, l, "migrate: number of places transfered 10")
		assert.Contains(t, l, "migrate: number of reactions transfered 3")
		assert.Contains(t, l, "migrate: number of sessions transfered 21")
		assert.Contains(t, l, "migrate: number of services transfered 2")
		assert.Contains(t, l, "migrate: number of subjects transfered 6")
		assert.Contains(t, l, "migrate: number of users transfered 11")
		assert.Contains(t, l, "migrate: number of userdetails transfered 9")
		assert.Contains(t, l, "migrate: number of usersettings transfered 13")
		assert.Contains(t, l, "migrate: number of usershares transfered 1")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db")
		}
	})

	t.Run("SQLiteToMySQLPopulated", func(t *testing.T) {
		// Remove target database file
		os.Remove("/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db")

		// Load migrate database as source
		if err := copyFile("/go/src/github.com/photoprism/photoprism/internal/commands/testdata/transfer_sqlite3", "/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db"); err != nil {
			t.Fatal(err.Error())
		}

		// Load migrate database as target
		if dumpName, err := filepath.Abs("./testdata/transfer_mysql"); err != nil {
			t.Fatal(err)
		} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
			"-e", "source "+dumpName).Run(); err != nil {
			t.Fatal(err)
		}

		// Run command with test context.
		log = event.Log

		appArgs := []string{"photoprism",
			"--database-driver", "sqlite",
			"--database-dsn", "/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db?_busy_timeout=5000&_foreign_keys=on",
			"--transfer-driver", "mysql",
			"--transfer-dsn", "migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=15s"}
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

		assert.Contains(t, l, "migrate: number of albums transfered 31")
		assert.Contains(t, l, "migrate: number of albumusers transfered 0")
		assert.Contains(t, l, "migrate: number of cameras transfered 6")
		assert.Contains(t, l, "migrate: number of categories transfered 1")
		assert.Contains(t, l, "migrate: number of cells transfered 9")
		assert.Contains(t, l, "migrate: number of clients transfered 7")
		assert.Contains(t, l, "migrate: number of countries transfered 1")
		assert.Contains(t, l, "migrate: number of duplicates transfered 0")
		assert.Contains(t, l, "migrate: number of errors transfered 0")
		assert.Contains(t, l, "migrate: number of faces transfered 7")
		assert.Contains(t, l, "migrate: number of files transfered 71")
		assert.Contains(t, l, "migrate: number of fileshares transfered 2")
		assert.Contains(t, l, "migrate: number of filesyncs transfered 3")
		assert.Contains(t, l, "migrate: number of folders transfered 3")
		assert.Contains(t, l, "migrate: number of keywords transfered 26")
		assert.Contains(t, l, "migrate: number of labels transfered 32")
		assert.Contains(t, l, "migrate: number of lenses transfered 2")
		assert.Contains(t, l, "migrate: number of links transfered 5")
		assert.Contains(t, l, "migrate: number of markers transfered 18")
		assert.Contains(t, l, "migrate: number of passcodes transfered 3")
		assert.Contains(t, l, "migrate: number of passwords transfered 11")
		assert.Contains(t, l, "migrate: number of photos transfered 58")
		assert.Contains(t, l, "migrate: number of photousers transfered 0")
		assert.Contains(t, l, "migrate: number of places transfered 10")
		assert.Contains(t, l, "migrate: number of reactions transfered 3")
		assert.Contains(t, l, "migrate: number of sessions transfered 21")
		assert.Contains(t, l, "migrate: number of services transfered 2")
		assert.Contains(t, l, "migrate: number of subjects transfered 6")
		assert.Contains(t, l, "migrate: number of users transfered 11")
		assert.Contains(t, l, "migrate: number of userdetails transfered 9")
		assert.Contains(t, l, "migrate: number of usersettings transfered 13")
		assert.Contains(t, l, "migrate: number of usershares transfered 1")
		// Remove target database file
		if !t.Failed() {
			os.Remove("/go/src/github.com/photoprism/photoprism/storage/sqlitetomysql.test.db")
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

package performancetest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/migrate"
	"github.com/photoprism/photoprism/pkg/fs"
)

func sqliteMigration(original string, temp string, numberOfRecords int, skipSpeedup bool, testname string, expectedDuration time.Duration, b *testing.B) {

	b.StopTimer()
	// Prepare temporary sqlite db.
	testDbOriginal := original
	testDbTemp := temp
	dumpName, err := filepath.Abs(testDbTemp)
	_ = os.Remove(dumpName)
	if err != nil {
		b.Fatal(err)
	} else if err = fs.Copy(testDbOriginal, dumpName); err != nil {
		b.Fatal(err)
	}
	defer os.Remove(dumpName)
	b.StartTimer()

	log = logrus.StandardLogger()
	log.SetLevel(logrus.ErrorLevel)

	start := time.Now()
	dsn := fmt.Sprintf("%v?_busy_timeout=5000", dumpName)

	db, err := gorm.Open(
		"sqlite3",
		dsn,
	)

	if err != nil || db == nil {
		if err != nil {
			b.Fatal(err)
		}

		return
	}

	defer db.Close()

	db.LogMode(false)
	db.SetLogger(log)

	if err != nil || db == nil {
		if err != nil {
			b.Fatal(err)
		}

		return
	}

	opt := migrate.Opt(true, true, nil)

	// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
	if err = db.AutoMigrate(&migrate.Migration{}).Error; err != nil {
		b.Fatal(err)
	}
	if err = db.AutoMigrate(&migrate.Version{}).Error; err != nil {
		b.Fatal(err)
	}

	// Setup and capture SQL Logging output
	buffer := bytes.Buffer{}
	log.SetOutput(&buffer)

	entity.Entities.Migrate(db, opt)
	// The bad thing is that the above panics, but doesn't return an error.

	// Reset logger
	log.SetOutput(os.Stdout)

	// Expect 0 errors (no such table accounts, and missing account_id in files_sync and files_share)
	// And a blank record.
	assert.Equal(b, 1, len(strings.Split(buffer.String(), "\n")))
	if len(strings.Split(buffer.String(), "\n")) == 1 {
		assert.Equal(b, 0, len(strings.Split(buffer.String(), "\n")[0]))
	} else {
		log.Error("Migration result not as expected.  Results follow:")
		for i := 0; i < len(strings.Split(buffer.String(), "\n")); i++ {
			log.Error(strings.Split(buffer.String(), "\n")[i])
		}
	}

	elapsed := time.Since(start)

	stmt := db.Table("photos").Where("photo_uid IS NOT NULL")

	count := int64(0)

	// Fetch count from database.
	if err = stmt.Count(&count).Error; err != nil {
		b.Error(err)
	} else {
		assert.Equal(b, int64(numberOfRecords), count)
	}

	log.Info(testname, " sqlite took ", elapsed)
	assert.LessOrEqual(b, elapsed, expectedDuration)
}

func mysqlMigration(testDbOriginal string, numberOfRecords int, testname string, expectedDuration time.Duration, b *testing.B) {
	b.StopTimer()
	// Prepare migrate mariadb db.
	if dumpName, err := filepath.Abs(testDbOriginal); err != nil {
		b.Fatal(err)
	} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
		"-e", "source "+dumpName).Run(); err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	start := time.Now()

	log = logrus.StandardLogger()
	log.SetLevel(logrus.ErrorLevel)

	db, err := gorm.Open(
		"mysql",
		"migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
	)

	if err != nil || db == nil {
		if err != nil {
			b.Fatal(err)
		}

		return
	}

	defer db.Close()

	db.LogMode(false)
	db.SetLogger(log)

	opt := migrate.Opt(true, true, nil)

	// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
	if err = db.AutoMigrate(&migrate.Migration{}).Error; err != nil {
		b.Fatal(err)
	}
	if err = db.AutoMigrate(&migrate.Version{}).Error; err != nil {
		b.Fatal(err)
	}

	// Setup and capture SQL Logging output
	buffer := bytes.Buffer{}
	log.SetOutput(&buffer)

	entity.Entities.Migrate(db, opt)
	// The bad thing is that the above panics, but doesn't return an error.

	// Reset logger
	log.SetOutput(os.Stdout)

	// Expect 3 errors (no such table accounts, and missing account_id in files_sync and files_share)
	// And a blank record.
	assert.Equal(b, 1, len(strings.Split(buffer.String(), "\n")))
	if len(strings.Split(buffer.String(), "\n")) == 1 {
		assert.Equal(b, 0, len(strings.Split(buffer.String(), "\n")[0]))
	} else {
		log.Error("Migration result not as expected.  Results follow:")
		for i := 0; i < len(strings.Split(buffer.String(), "\n")); i++ {
			log.Error(strings.Split(buffer.String(), "\n")[i])
		}
	}

	elapsed := time.Since(start)

	stmt := db.Table("photos").Where("photo_uid IS NOT NULL")

	count := int64(0)

	// Fetch count from database.
	if err = stmt.Count(&count).Error; err != nil {
		b.Error(err)
	} else {
		assert.Equal(b, int64(numberOfRecords), count)
	}

	log.Info(testname, " mysql took ", elapsed)
	assert.LessOrEqual(b, elapsed, expectedDuration)
}

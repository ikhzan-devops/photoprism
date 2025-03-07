package migrate

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDialectPostgreSQL(t *testing.T) {
	if dumpName, err := filepath.Abs("./testdata/migrate_postgres.sql"); err != nil {
		t.Fatal(err)
	} else {
		if err = exec.Command("psql", "postgresql://photoprism:photoprism@postgres:5432/postgres", "--file="+dumpName).Run(); err != nil {
			t.Fatal(err)
		}
	}

	log.Info("Expect potential table does not exist or no such table Error or SQLSTATE from migration.go")
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	db, err := gorm.Open(postgres.Open(
		"postgresql://migrate:migrate@postgres:5432/migrate?TimeZone=UTC&connect_timeout=15&lock_timeout=5000&sslmode=disable"),
		&gorm.Config{
			Logger: logger.New(
				log,
				logger.Config{
					SlowThreshold:             time.Second,  // Slow SQL threshold
					LogLevel:                  logger.Error, // Log level
					IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
					ParameterizedQueries:      true,         // Don't include params in the SQL log
					Colorful:                  false,        // Disable color
				},
			),
		},
	)

	if err != nil || db == nil {
		if err != nil {
			t.Fatal(err)
		}

		return
	}

	sqldb, _ := db.DB()
	defer sqldb.Close()

	opt := Opt(true, true, nil)

	// Run pre-migrations.
	if err = Run(db, opt.Pre()); err != nil {
		t.Error(err)
	}

	// Run migrations.
	if err = Run(db, opt); err != nil {
		t.Error(err)
	}

	stmt := db.Table("photos").Where("photo_caption = '' OR photo_caption IS NULL")

	count := int64(0)

	// Fetch count from database.
	if err = stmt.Count(&count).Error; err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, int64(77), count)
	}
	log.Info("End Expect potential table does not exist or no such table Error or SQLSTATE from migration.go")
}

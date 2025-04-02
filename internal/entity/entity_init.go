package entity

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/photoprism/photoprism/internal/entity/migrate"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// onReady contains init functions to be called when the
// initialization of the database is complete.
var onReady []func()

// ready executes init callbacks when the initialization of the
// database is complete.
func ready() {
	for _, init := range onReady {
		init()
	}
}

// InitDb creates database tables and inserts default fixtures as needed.
func InitDb(opt migrate.Options) {
	if !HasDbProvider() {
		log.Error("migrate: no database provider")
		return
	}

	start := time.Now()

	if opt.DropDeprecated {
		DeprecatedTables.Drop(Db())
	}

	Entities.Migrate(Db(), opt)

	if err := Entities.WaitForMigration(Db()); err != nil {
		log.Errorf("migrate: %s", err)
	}

	CreateDefaultFixtures()

	ready()

	log.Debugf("migrate: completed in %s", time.Since(start))
}

// InitTestDb connects to and completely initializes the test database incl fixtures.
func InitTestDb(driver, dsn string) *DbConn {
	if HasDbProvider() {
		return nil
	}

	start := time.Now()

	// Set default test database driver.
	if driver == "test" || driver == "sqlite" || driver == "" || dsn == "" {
		driver = SQLite3
	}

	// Set default database DSN.
	if driver == SQLite3 {
		if dsn == "" || dsn == SQLiteTestDB {
			dsn = fmt.Sprintf("%s?_foreign_keys=on", SQLiteTestDB)
			if !fs.FileExists(SQLiteTestDB) {
				log.Debugf("sqlite: test database %s does not already exist", clean.Log(SQLiteTestDB))
			} else if err := os.Remove(SQLiteTestDB); err != nil {
				log.Errorf("sqlite: failed to remove existing test database %s (%s)", clean.Log(SQLiteTestDB), err)
			}
		}
	}

	log.Infof("initializing %s test db in %s", driver, dsn)

	allowDelete := os.Getenv("PHOTOPRISM_TEST_DBDROP")
	if driver == MySQL && allowDelete == "true" {
		basedsn := dsn[0 : strings.Index(dsn, "/")+1]
		basedbname := dsn[strings.Index(dsn, "/")+1 : strings.Index(dsn, "?")]
		log.Infof("Connecting to %v", basedsn)
		database, err := gorm.Open(mysql.Open(basedsn), &gorm.Config{})
		if err != nil {
			log.Errorf("Unable to connect to MariaDB %v", err)
		}
		log.Infof("Dropping database %v if it exists", basedbname)
		if res := database.Exec("DROP DATABASE IF EXISTS " + basedbname + ";"); res.Error != nil {
			log.Errorf("Unable to drop database %v", res.Error)
			return nil
		}

		log.Infof("Creating database %v if it doesnt exist", basedbname)
		if res := database.Exec("CREATE DATABASE IF NOT EXISTS " + basedbname + ";"); res.Error != nil {
			log.Errorf("Unable to create database %v", res.Error)
			return nil
		}
	}

	// Create gorm.DB connection provider.
	db := &DbConn{
		Driver: driver,
		Dsn:    dsn,
	}

	// Insert test fixtures into the database.
	SetDbProvider(db)
	ResetTestFixtures()
	File{}.RegenerateIndex()

	ready()

	log.Debugf("migrate: completed in %s", time.Since(start))

	return db
}

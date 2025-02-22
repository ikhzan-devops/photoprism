package performancetest

import (
	"fmt"
	"sync"
	"time"

<<<<<<< HEAD
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
=======
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
>>>>>>> origin/Benchmarking
)

// Supported test databases.
const (
	MySQL           = "mysql"
<<<<<<< HEAD
	Postgres        = "postgres"
	SQLite3         = "sqlite"
	SQLiteTestDB    = ".test.db"
	SQLiteMemoryDSN = ":memory:?cache=shared&_foreign_keys=on"
)

var drivers = map[string]func(string) gorm.Dialector{
	MySQL:    mysql.Open,
	Postgres: postgres.Open,
	SQLite3:  sqlite.Open,
}

=======
	SQLite3         = "sqlite3"
	SQLiteTestDB    = ".test.db"
	SQLiteMemoryDSN = ":memory:?cache=shared"
)

>>>>>>> origin/Benchmarking
// dbConn is the global gorm.DB connection provider.
var dbConn Gorm

// Gorm is a gorm.DB connection provider interface.
type Gorm interface {
	Db() *gorm.DB
}

// DbConn is a gorm.DB connection provider.
type DbConn struct {
	Driver string
	Dsn    string

	once sync.Once
	db   *gorm.DB
}

// Db returns the gorm db connection.
func (g *DbConn) Db() *gorm.DB {
	g.once.Do(g.Open)

	if g.db == nil {
		log.Fatal("migrate: database not connected")
	}

	return g.db
}

// Open creates a new gorm db connection.
func (g *DbConn) Open() {
<<<<<<< HEAD
	log.Infof("Opening DB connection with driver %s", g.Driver)
	db, err := gorm.Open(drivers[g.Driver](g.Dsn), gormConfig())
=======
	db, err := gorm.Open(g.Driver, g.Dsn)
>>>>>>> origin/Benchmarking

	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			fmt.Printf("gorm.Open(%s, %s) %d\n", g.Driver, g.Dsn, i)
<<<<<<< HEAD
			db, err = gorm.Open(drivers[g.Driver](g.Dsn), gormConfig())
=======
			db, err = gorm.Open(g.Driver, g.Dsn)
>>>>>>> origin/Benchmarking

			if db != nil && err == nil {
				break
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		if err != nil || db == nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
<<<<<<< HEAD
	log.Info("DB connection established successfully")

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(4)   // in config_db it uses c.DatabaseConnsIdle(), but we don't have the c here.
	sqlDB.SetMaxOpenConns(256) // in config_db it uses c.DatabaseConns(), but we don't have the c here.
=======

	db.LogMode(true)
	db.SetLogger(log)
	db.DB().SetMaxIdleConns(4)
	db.DB().SetMaxOpenConns(256)
>>>>>>> origin/Benchmarking

	g.db = db
}

// Close closes the gorm db connection.
func (g *DbConn) Close() {
	if g.db != nil {
<<<<<<< HEAD
		sqlDB, _ := g.db.DB()
		if err := sqlDB.Close(); err != nil {
=======
		if err := g.db.Close(); err != nil {
>>>>>>> origin/Benchmarking
			log.Fatal(err)
		}

		g.db = nil
	}
}

// IsDialect returns true if the given sql dialect is used.
func IsDialect(name string) bool {
<<<<<<< HEAD
	return name == Db().Dialector.Name()
=======
	return name == dbConn.Db().Dialect().GetName()
>>>>>>> origin/Benchmarking
}

// DbDialect returns the sql dialect name.
func DbDialect() string {
<<<<<<< HEAD
	return Db().Dialector.Name()
=======
	return dbConn.Db().Dialect().GetName()
>>>>>>> origin/Benchmarking
}

// SetDbProvider sets the Gorm database connection provider.
func SetDbProvider(conn Gorm) {
	dbConn = conn
}

// HasDbProvider returns true if a db provider exists.
func HasDbProvider() bool {
	return dbConn != nil
}

<<<<<<< HEAD
func gormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log,
			logger.Config{
				SlowThreshold:             time.Second,  // Slow SQL threshold
				LogLevel:                  logger.Error, // Log level
				IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,        // Don't include params in the SQL log
				Colorful:                  false,        // Disable color
			},
		),
		// Set UTC as the default for created and updated timestamps.
		NowFunc: func() time.Time {
			return UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
=======
// Db returns the default *gorm.DB connection.
func Db() *gorm.DB {
	if dbConn == nil {
		return nil
	}

	return dbConn.Db()
}

// UnscopedDb returns an unscoped *gorm.DB connection
// that returns all records including deleted records.
func UnscopedDb() *gorm.DB {
	return Db().Unscoped()
>>>>>>> origin/Benchmarking
}

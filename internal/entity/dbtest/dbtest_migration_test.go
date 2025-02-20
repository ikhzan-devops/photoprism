package entity

import (
	"bytes"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/migrate"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
)

func TestDialectSQLite3(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	dbtestMutex.Lock()
	defer dbtestMutex.Unlock()

	t.Run("ValidMigration", func(t *testing.T) {
		// Prepare temporary sqlite db.
		testDbOriginal := "../migrate/testdata/migrate_sqlite3"
		testDbTemp := "../migrate/testdata/migrate_sqlite3.db"
		if !fs.FileExists(testDbOriginal) {
			t.Fatal(testDbOriginal + " not found")
		}
		dumpName, err := filepath.Abs(testDbTemp)
		_ = os.Remove(dumpName)
		if err != nil {
			t.Fatal(err)
		} else if err = fs.Copy(testDbOriginal, dumpName); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(dumpName)

		log = logrus.StandardLogger()
		log.SetLevel(logrus.TraceLevel)

		dsn := fmt.Sprintf("%v?_foreign_keys=on&_busy_timeout=5000", dumpName)

		db, err := gorm.Open(sqlite.Open(dsn),
			&gorm.Config{
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

		opt := migrate.Opt(true, true, nil)

		// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
		if err = db.AutoMigrate(&migrate.Migration{}); err != nil {
			t.Error(err)
		}
		if err = db.AutoMigrate(&migrate.Version{}); err != nil {
			t.Error(err)
		}

		// Setup and capture SQL Logging output
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)

		entity.Entities.Migrate(db, opt)
		// The bad thing is that the above panics, but doesn't return an error.

		// Reset logger
		log.SetOutput(os.Stdout)

		for i := 0; i < len(strings.Split(buffer.String(), "\n")); i++ {
			log.Info(strings.Split(buffer.String(), "\n")[i])
		}

		// Expect 4 errors (auth_? tables missing)
		// And a blank record.
		assert.Equal(t, 5, len(strings.Split(buffer.String(), "\n")))
		assert.Equal(t, 0, len(strings.Split(buffer.String(), "\n")[4]))

		stmt := db.Table("photos").Where("photo_caption = '' OR photo_caption IS NULL")

		count := int64(0)

		// Fetch count from database.
		if err = stmt.Count(&count).Error; err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int64(0), count)
		}

		// Test that the minimum values can be added to the database
		populatePhotoPrismStructsWithMin(t, db)

		// Test that the maximum values can be added to the database
		populatePhotoPrismStructsWithMax(t, db)

	})

	t.Run("InvalidDataUpgrade", func(t *testing.T) {
		// Prepare temporary sqlite db.
		testDbOriginal := "../migrate/testdata/migrate_sqlite3"
		testDbTemp := "../migrate/testdata/migrate_sqlite3.db"
		if !fs.FileExists(testDbOriginal) {
			t.Fatal(testDbOriginal + " not found")
		}
		dumpName, err := filepath.Abs(testDbTemp)
		_ = os.Remove(dumpName)
		if err != nil {
			t.Fatal(err)
		} else if err = fs.Copy(testDbOriginal, dumpName); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(dumpName)

		log = logrus.StandardLogger()
		log.SetLevel(logrus.TraceLevel)

		dsn := fmt.Sprintf("%v?_foreign_keys=on&_busy_timeout=5000", dumpName)

		db, err := gorm.Open(sqlite.Open(dsn),
			&gorm.Config{
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

		// Load some invalid data into the database
		AlbumUID := byte('a')
		PhotoUID := byte('p')

		_, err = sqldb.Exec("INSERT INTO photos_albums (photo_uid, album_uid, `order`, hidden, missing, created_at, updated_at) VALUES (?, ?, '0', '0', '0', ?, ?)", rnd.GenerateUID(PhotoUID), rnd.GenerateUID(AlbumUID), time.Now().UTC(), time.Now().UTC())
		assert.Nil(t, err)

		opt := migrate.Opt(true, true, nil)

		// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
		err = db.AutoMigrate(&migrate.Migration{})
		err = db.AutoMigrate(&migrate.Version{})

		// Setup and capture SQL Logging output
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)

		entity.Entities.Migrate(db, opt)
		// The bad thing is that the above panics, but doesn't return an error.

		// Reset logger
		log.SetOutput(os.Stdout)

		// Expect 4 errors (auth_? tables missing)
		// And a blank record.
		assert.Equal(t, 5, len(strings.Split(buffer.String(), "\n")))
		assert.Equal(t, 0, len(strings.Split(buffer.String(), "\n")[4]))

		stmt := db.Table("photos").Where("photo_caption = '' OR photo_caption IS NULL")

		count := int64(0)

		// Fetch count from database.
		if err = stmt.Count(&count).Error; err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int64(0), count)
		}
	})

}

func TestDialectMysql(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	dbtestMutex.Lock()
	defer dbtestMutex.Unlock()

	t.Run("ValidMigration", func(t *testing.T) {
		// Prepare migrate mariadb db.
		if dumpName, err := filepath.Abs("../migrate/testdata/migrate_mysql.sql"); err != nil {
			t.Fatal(err)
		} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
			"-e", "source "+dumpName).Run(); err != nil {
			t.Fatal(err)
		}

		log = logrus.StandardLogger()
		log.SetLevel(logrus.TraceLevel)

		db, err := gorm.Open(mysql.Open(
			"migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true"),
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

		opt := migrate.Opt(true, true, nil)

		// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
		if err = db.AutoMigrate(&migrate.Migration{}); err != nil {
			t.Error(err)
		}
		if err = db.AutoMigrate(&migrate.Version{}); err != nil {
			t.Error(err)
		}

		// Setup and capture SQL Logging output
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)

		entity.Entities.Migrate(db, opt)
		// The bad thing is that the above panics, but doesn't return an error.

		// Reset logger
		log.SetOutput(os.Stdout)

		assert.Contains(t, buffer.String(), "Table 'migrate.auth_sessions' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_details' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_settings' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_shares' doesn't exist")
		// There is a blank record.
		assert.Equal(t, 5, len(strings.Split(buffer.String(), "\n")))
		if len(strings.Split(buffer.String(), "\n")) != 5 {
			for i := 0; i < len(strings.Split(buffer.String(), "\n")); i++ {
				assert.Empty(t, strings.Split(buffer.String(), "\n")[i])
			}
		}
		// Detect a foreign key issue
		assert.NotContains(t, buffer.String(), "errno: 150")

		stmt := db.Table("photos").Where("photo_caption = '' OR photo_caption IS NULL")

		count := int64(0)

		// Fetch count from database.
		if err = stmt.Count(&count).Error; err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int64(0), count)
		}

		// Test that the minimum values can be added to the database
		populatePhotoPrismStructsWithMin(t, db)

		// Test that the maximum values can be added to the database
		populatePhotoPrismStructsWithMax(t, db)

	})

	t.Run("InvalidDataUpgrade", func(t *testing.T) {
		// Prepare migrate mariadb db.
		if dumpName, err := filepath.Abs("../migrate/testdata/migrate_mysql.sql"); err != nil {
			t.Fatal(err)
		} else if err = exec.Command("mariadb", "-u", "migrate", "-pmigrate", "migrate",
			"-e", "source "+dumpName).Run(); err != nil {
			t.Fatal(err)
		}

		log = logrus.StandardLogger()
		log.SetLevel(logrus.TraceLevel)

		db, err := gorm.Open(mysql.Open(
			"migrate:migrate@tcp(mariadb:4001)/migrate?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true"),
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

		// Load some invalid data into the database
		AlbumUID := byte('a')
		PhotoUID := byte('p')

		_, err = sqldb.Exec("INSERT INTO photos_albums (photo_uid, album_uid, `order`, hidden, missing, created_at, updated_at) VALUES (?, ?, '0', '0', '0', ?, ?)", rnd.GenerateUID(PhotoUID), rnd.GenerateUID(AlbumUID), time.Now().UTC(), time.Now().UTC())
		assert.Nil(t, err)

		opt := migrate.Opt(true, true, nil)

		// Make sure that migrate and version is done, as the Once doesn't work as it has already been set before we opened the new database..
		err = db.AutoMigrate(&migrate.Migration{})
		err = db.AutoMigrate(&migrate.Version{})

		// Setup and capture SQL Logging output
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)

		entity.Entities.Migrate(db, opt)
		// The bad thing is that the above panics, but doesn't return an error.

		// Reset logger
		log.SetOutput(os.Stdout)

		assert.Contains(t, buffer.String(), "Table 'migrate.auth_sessions' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_details' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_settings' doesn't exist")
		assert.Contains(t, buffer.String(), "Table 'migrate.auth_users_shares' doesn't exist")
		// There is a blank record.
		assert.Equal(t, 5, len(strings.Split(buffer.String(), "\n")))
		if len(strings.Split(buffer.String(), "\n")) != 5 {
			for i := 0; i < len(strings.Split(buffer.String(), "\n")); i++ {
				assert.Empty(t, strings.Split(buffer.String(), "\n")[i])
			}
		}

		stmt := db.Table("photos").Where("photo_caption = '' OR photo_caption IS NULL")

		count := int64(0)

		// Fetch count from database.
		if err = stmt.Count(&count).Error; err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int64(0), count)
		}
	})

}

func populatePhotoPrismStructsWithMin(t *testing.T, db *gorm.DB) {
	savedID := uint(123456789)

	albumUser := entity.AlbumUser{}
	if err := populateStructWithMin(&albumUser); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&albumUser); result.Error != nil {
			t.Error(result.Error)
		}
	}

	album := entity.Album{}
	if err := populateStructWithMin(&album); err != nil {
		t.Error(err)
	} else {
		album.ID = savedID
		if result := db.Create(&album); result.Error != nil {
			t.Error(result.Error)
		}
	}

	client := entity.Client{}
	if err := populateStructWithMin(&client); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&client); result.Error != nil {
			t.Error(result.Error)
		}
	}

	session := entity.Session{}
	if err := populateStructWithMin(&session); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&session); result.Error != nil {
			t.Error(result.Error)
		}
	}

	photo := entity.Photo{}
	if err := populateStructWithMin(&photo); err != nil {
		t.Error(err)
	} else {
		photo.ID = savedID
		if result := db.Create(&photo); result.Error != nil {
			t.Error(result.Error)
		}
	}
}

func populatePhotoPrismStructsWithMax(t *testing.T, db *gorm.DB) {
	uintMaxInt64 := db.Dialector.Name() == entity.SQLite3

	savedID := uint(math.MaxUint)
	if uintMaxInt64 {
		savedID = uint(math.MaxInt64)
	}

	albumUser := entity.AlbumUser{}
	if err := populateStructWithMax(&albumUser, uintMaxInt64); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&albumUser); result.Error != nil {
			t.Error(result.Error)
		}
	}

	album := entity.Album{}
	if err := populateStructWithMax(&album, uintMaxInt64); err != nil {
		t.Error(err)
	} else {
		album.ID = savedID
		if result := db.Create(&album); result.Error != nil {
			t.Error(result.Error)
		}
	}

	client := entity.Client{}
	if err := populateStructWithMax(&client, uintMaxInt64); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&client); result.Error != nil {
			t.Error(result.Error)
		}
	}

	session := entity.Session{}
	if err := populateStructWithMax(&session, uintMaxInt64); err != nil {
		t.Error(err)
	} else {
		if result := db.Create(&session); result.Error != nil {
			t.Error(result.Error)
		}
	}

	photo := entity.Photo{}
	if err := populateStructWithMax(&photo, uintMaxInt64); err != nil {
		t.Error(err)
	} else {
		photo.CameraID = 0 // use the records created above.
		photo.CellID = "zz"
		photo.LensID = 0
		photo.PlaceID = "zz"

		photo.ID = savedID
		//log.Debug(photo)
		if result := db.Create(&photo); result.Error != nil {
			t.Error(result.Error)
		}
	}

}

func populateStructWithMin(m interface{}) (err error) {

	r := reflect.ValueOf(m)

	if r.Kind() != reflect.Pointer {
		return fmt.Errorf("model interface expected")
	}

	values := r.Elem()

	if kind := values.Kind(); kind != reflect.Struct {
		return fmt.Errorf("model expected")
	}

	reType := regexp.MustCompile(`(?:type:)(?P<type>[a-zA-Z]*)`)
	reSize := regexp.MustCompile(`(?:size:)(?P<size>[0-9]*)`)

	vT := values.Type()
	num := vT.NumField()
	for i := 0; i < num; i++ {
		field := vT.Field(i)

		// Skip non-exported fields.
		if !field.IsExported() {
			continue
		}

		// fieldName := field.Name

		v := values.Field(i)
		tag, _ := readTag(field, "gorm")
		// log.Debugf("field %v is %v with name %v or string %v and is exported %v with tag %v", fieldName, v.Kind(), v.Type().Name(), v.Type().String(), field.IsExported(), tag)

		// Skip read-only fields.
		if !v.CanSet() {
			continue
		}

		// Skip ignored fields by tag
		if tag == "-" {
			continue
		}

		if strings.Contains(tag, "type") {
			typeString := reFindStringGroup(reType, tag, "type")
			switch typeString {
			// case "bytes":
			// 	v.Set(reflect.ValueOf(""))
			case "float":
				strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
				if strLen == 32 {
					v.SetFloat(math.SmallestNonzeroFloat32)
				} else {
					v.SetFloat(math.SmallestNonzeroFloat64)
				}
			case "int":
				strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
				switch strLen {
				case 8:
					v.SetInt(math.MinInt8)
				case 16:
					v.SetInt(math.MinInt16)
				case 32:
					v.SetInt(math.MinInt32)
				case 64:
					v.SetInt(math.MinInt64)
				}
			}
		} else {
			switch v.Kind() {
			case reflect.Struct:
				switch v.Type().String() {
				case "time.Time":
					v.Set(reflect.ValueOf(time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC))) // Mariadb limitation
				default:
					// case "sql.NullTime", "time.Duration":
					log.Debugf("reflect.Struct with Unhandled type %s for field %s", v.Type().String(), field.Name)

				}
			case reflect.Pointer:
				switch v.Type().String() {
				case "*time.Time":
					timeTime := time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC)
					v.Set(reflect.ValueOf(&timeTime)) // Mariadb limitation
				default:
					//case "*time.Time", "*time.Duration", "*bool", "*uint", "*uint64", "*uint32", "*int", "*int64", "*int32", "*string", "*float32", "*float64", "*otp.Key", "*sql.NullTime", "*json.RawMessage":
					log.Debugf("reflect.Pointer with Unhandled type %s for field %s", v.Type().String(), field.Name)
				}
			case reflect.Uint:
				v.Set(reflect.ValueOf(uint(0)))
			case reflect.Int8:
				v.SetInt(math.MinInt8)
			case reflect.Int16:
				v.SetInt(math.MinInt16)
			case reflect.Int32:
				v.SetInt(math.MinInt32)
			case reflect.Int64:
				v.SetInt(math.MinInt64)
			case reflect.Int:
				if strings.Contains(tag, "size") {
					tagSize, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
					switch tagSize {
					case 8:
						v.SetInt(math.MinInt8)
					case 16:
						v.SetInt(math.MinInt16)
					case 32:
						v.SetInt(math.MinInt32)
					case 64:
						v.SetInt(math.MinInt64)
					}
				} else {
					v.SetInt(math.MinInt)
				}
			case reflect.Float32:
				v.SetFloat(math.SmallestNonzeroFloat32)
			case reflect.Float64:
				v.SetFloat(math.SmallestNonzeroFloat64)
			case reflect.String:
				if strings.Contains(tag, "type") {
					typeString := reFindStringGroup(reType, tag, "type")
					switch typeString {
					case "bytes":
						v.Set(reflect.ValueOf(""))
					case "float":
						strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
						if strLen == 32 {
							v.SetFloat(math.SmallestNonzeroFloat32)
						} else {
							v.SetFloat(math.SmallestNonzeroFloat64)
						}
					case "int":
						strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
						switch strLen {
						case 8:
							v.SetInt(math.MinInt8)
						case 16:
							v.SetInt(math.MinInt16)
						case 32:
							v.SetInt(math.MinInt32)
						case 64:
							v.SetInt(math.MinInt64)
						}
					}
				} else if strings.Contains(tag, "size") {
					v.Set(reflect.ValueOf(""))
				}
			case reflect.Bool:
				v.SetBool(false)
			default:
				log.Debugf("vKind %s with Unhandled type %s for field %s", v.Kind(), v.Type().String(), field.Name)
			}
		}
	}
	return nil
}

func populateStructWithMax(m interface{}, uintMaxInt64 bool) (err error) {

	r := reflect.ValueOf(m)

	if r.Kind() != reflect.Pointer {
		return fmt.Errorf("model interface expected")
	}

	values := r.Elem()

	if kind := values.Kind(); kind != reflect.Struct {
		return fmt.Errorf("model expected")
	}

	reType := regexp.MustCompile(`(?:type:)(?P<type>[a-zA-Z]*)`)
	reSize := regexp.MustCompile(`(?:size:)(?P<size>[0-9]*)`)

	vT := values.Type()
	num := vT.NumField()
	for i := 0; i < num; i++ {
		field := vT.Field(i)

		// Skip non-exported fields.
		if !field.IsExported() {
			continue
		}

		// fieldName := field.Name

		v := values.Field(i)
		tag, _ := readTag(field, "gorm")
		// log.Debugf("field %v is %v with name %v or string %v and is exported %v with tag %v", fieldName, v.Kind(), v.Type().Name(), v.Type().String(), field.IsExported(), tag)

		// Skip read-only fields.
		if !v.CanSet() {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			switch v.Type().String() {
			case "time.Time":
				v.Set(reflect.ValueOf(time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC))) // Mariadb limitation
				//case "sql.NullTime", "time.Duration":

			}
		case reflect.Pointer:
			switch v.Type().String() {
			case "*time.Time":
				timeTime := time.Date(9999, 12, 31, 23, 59, 59, 999999, time.UTC)
				v.Set(reflect.ValueOf(&timeTime)) // Mariadb limitation

				// case "*time.Time", "*time.Duration", "*bool", "*uint", "*uint64", "*uint32", "*int", "*int64", "*int32", "*string", "*float32", "*float64", "*otp.Key", "*sql.NullTime", "*json.RawMessage":
			}
		case reflect.Uint:
			if strings.Contains(tag, "size") {
				tagSize, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
				switch tagSize {
				case 8:
					v.Set(reflect.ValueOf(uint(math.MaxUint8)))
				case 16:
					v.Set(reflect.ValueOf(uint(math.MaxUint16)))
				case 32:
					v.Set(reflect.ValueOf(uint(math.MaxUint32)))
				case 64:
					if uintMaxInt64 {
						v.Set(reflect.ValueOf(uint(math.MaxUint32)))
					} else {
						v.Set(reflect.ValueOf(uint(math.MaxUint64)))
					}
				}
			} else {
				if uintMaxInt64 {
					v.Set(reflect.ValueOf(uint(math.MaxUint32)))
				} else {
					v.Set(reflect.ValueOf(uint(math.MaxUint64)))
				}
			}
		case reflect.Int8:
			v.SetInt(math.MaxInt8)
		case reflect.Int16:
			v.SetInt(math.MaxInt16)
		case reflect.Int32:
			v.SetInt(math.MaxInt32)
		case reflect.Int64:
			v.SetInt(math.MaxInt64)
		case reflect.Int:
			if strings.Contains(tag, "size") {
				tagSize, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
				switch tagSize {
				case 8:
					v.SetInt(math.MaxInt8)
				case 16:
					v.SetInt(math.MaxInt16)
				case 32:
					v.SetInt(math.MaxInt32)
				case 64:
					v.SetInt(math.MaxInt64)
				}
			} else {
				v.SetInt(math.MaxInt)
			}
		case reflect.Float32:
			v.SetFloat(math.MaxFloat32)
		case reflect.Float64:
			v.SetFloat(math.MaxFloat64)
		case reflect.String:
			// tag
			if tag == "-" {
				continue
			}
			if strings.Contains(tag, "type") {
				typeString := reFindStringGroup(reType, tag, "type")
				switch typeString {
				case "bytes":
					strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
					if strLen <= 0 || strLen > 8192 {
						v.Set(reflect.ValueOf(randomString(8192)))
					} else {
						v.Set(reflect.ValueOf(randomString(strLen)))
					}
				case "float":
					strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
					if strLen == 32 {
						v.SetFloat(math.SmallestNonzeroFloat32)
					} else {
						v.SetFloat(math.SmallestNonzeroFloat64)
					}
				case "int":
					strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
					switch strLen {
					case 8:
						v.SetInt(math.MaxInt8)
					case 16:
						v.SetInt(math.MaxInt16)
					case 32:
						v.SetInt(math.MaxInt32)
					case 64:
						v.SetInt(math.MaxInt64)
					}
				}
			} else if strings.Contains(tag, "size") {
				strLen, _ := strconv.Atoi(reFindStringGroup(reSize, tag, "size"))
				if strLen <= 0 || strLen > 8192 {
					v.Set(reflect.ValueOf(randomString(8192)))
				} else {
					v.Set(reflect.ValueOf(randomString(strLen)))
				}
			}
		case reflect.Bool:
			v.SetBool(true)

		}
	}
	return nil
}

func reFindStringGroup(expression *regexp.Regexp, textToSearch string, groupName string) string {
	match := expression.FindStringSubmatch(textToSearch)
	groupIndex := expression.SubexpIndex(groupName)

	return match[groupIndex]
}

func readTag(f reflect.StructField, id string) (result string, valid bool) {
	val, ok := f.Tag.Lookup(id)
	if !ok {
		return id, false
	} else {
		return val, true
	}
}

const characterRunes = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(len int) string {
	sb := strings.Builder{}
	sb.Grow(len)
	for i := 0; i < len; {
		sb.WriteByte(characterRunes[rand.IntN(53)])
		i++
	}

	return sb.String()
}

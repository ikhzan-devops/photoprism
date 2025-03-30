package functions

import (
	"os"
)

// function to transform PHOTOPRISM_TEST_DSN environment variables to driver and dsn strings
func PhotoPrismTestToDriverDsn() (driver string, dsn string) {
	dsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
	switch dsnName {
	case "sqlite":
		driver = "sqlite"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_SQLITE")
	case "sqlitefile":
		driver = "sqlite"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_SQLITEFILE")
	case "mariadb":
		driver = "mysql"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_MARIADB")
	case "mysql8":
		driver = "mysql"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_MYSQL8")
	case "postgres":
		driver = "postgres"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_POSTGRES")
	default:
		driver = "sqlite"
		dsn = ""
	}
	return driver, dsn
}

// Gets the folder name to use to enforce folder separation for DBMS tests
func PhotoPrismTestToFolderName() (folderName string) {
	folderName = os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
	if folderName == "" {
		folderName, _ = PhotoPrismTestToDriverDsn()
	}
	return folderName
}

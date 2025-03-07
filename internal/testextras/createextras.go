package testextras

import (
	"time"

	"gorm.io/gorm"
)

func MigrateTestExtras(db *gorm.DB) {
	var err error
	for migrateRetry := 0; migrateRetry < 10; migrateRetry++ {
		if err = db.AutoMigrate(&TestLog{}); err != nil {
			log.Warnf("migratetestextras: automigrate testlog encountered %+v on loop %s", err, migrateRetry)
			if migrateRetry < 10 {
				time.Sleep(time.Second * 5)
			}
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}
	for migrateRetry := 0; migrateRetry < 10; migrateRetry++ {
		if err = db.AutoMigrate(&TestDBMutex{}); err != nil {
			log.Warnf("migratetestextras: automigrate testdbmutex encountered %+v on loop %s", err, migrateRetry)
			if migrateRetry < 10 {
				time.Sleep(time.Second * 5)
			}
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}
}

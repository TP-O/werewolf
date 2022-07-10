package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"uwwolf/database/migration"
	"uwwolf/database/seeder"
)

var dbInstance *gorm.DB

func init() {
	if db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{}); err != nil {
		panic("Failed to connect database")
	} else {
		dbInstance = db

		migration.Migrate(dbInstance)
		seeder.Seed(dbInstance)
	}
}

func DB() *gorm.DB {
	return dbInstance
}

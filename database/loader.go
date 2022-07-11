package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"uwwolf/database/migration"
	"uwwolf/database/seeder"
)

var dbInstance *gorm.DB

func LoadDatabase() error {
	if db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{}); err != nil {
		return err
	} else {
		dbInstance = db

		migration.Migrate(dbInstance)
		seeder.Seed(dbInstance)
	}

	return nil
}

package migration

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Team{},
		&model.Role{},
	)
}

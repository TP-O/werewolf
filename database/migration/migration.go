package migration

import (
	"gorm.io/gorm"

	"uwwolf/app/model"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Faction{},
		&model.Phase{},
		&model.Role{},
	)
}

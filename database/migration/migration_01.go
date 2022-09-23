package migration

import (
	"uwwolf/app/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration01 = &gormigrate.Migration{
	ID: "1",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.Phase{}, &model.Faction{}, &model.Role{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("roles", "factions", "phases")
	},
}

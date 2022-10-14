package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"uwwolf/app/model"
)

var migration01 = &gormigrate.Migration{
	ID: "1",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(
			&model.Phase{},
			&model.Faction{},
			&model.Role{},
			&model.Action{},
			&model.Status{},
			&model.Player{},
			&model.Game{},
			&model.GameRecord{},
			&model.RoleAssignment{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			"role_assignments",
			"game_records",
			"games",
			"users",
			"status",
			"actions",
			"roles",
			"factions",
			"phases",
		)
	},
}

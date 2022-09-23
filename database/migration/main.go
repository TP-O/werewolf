package migration

import (
	"uwwolf/database"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrations() *gormigrate.Gormigrate {
	client := database.Client()
	migration := gormigrate.New(client, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migration01,
	})

	return migration
}

package migration

import (
	"uwwolf/db"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrations() *gormigrate.Gormigrate {
	client := db.Client()
	migration := gormigrate.New(client, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migration01,
		migration02,
	})

	return migration
}
